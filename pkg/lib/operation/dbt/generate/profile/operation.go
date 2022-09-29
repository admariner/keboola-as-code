package profile

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"gopkg.in/yaml.v3"

	"github.com/keboola/keboola-as-code/internal/pkg/dbt"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/telemetry"
)

type dependencies interface {
	Logger() log.Logger
	Tracer() trace.Tracer
	LocalDbtProject(ctx context.Context) (*dbt.Project, bool, error)
}

const profilePath = "profiles.yml"

func Run(ctx context.Context, targetName string, d dependencies) (err error) {
	ctx, span := d.Tracer().Start(ctx, "kac.lib.operation.dbt.generate.profile")
	defer telemetry.EndSpan(span, &err)

	// Get dbt project
	project, _, err := d.LocalDbtProject(ctx)
	if err != nil {
		return err
	}
	fs := project.Fs()

	targetUpper := strings.ToUpper(targetName)
	profileDetails := map[string]interface{}{
		"target": targetName,
		"outputs": map[string]interface{}{
			targetName: map[string]interface{}{
				"account":   fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_ACCOUNT\") }}", targetUpper),
				"database":  fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_DATABASE\") }}", targetUpper),
				"password":  fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_PASSWORD\") }}", targetUpper),
				"schema":    fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_SCHEMA\") }}", targetUpper),
				"type":      fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_TYPE\") }}", targetUpper),
				"user":      fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_USER\") }}", targetUpper),
				"warehouse": fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_WAREHOUSE\") }}", targetUpper),
			},
		},
	}
	profilesFile := make(map[string]interface{})
	profilesFile["send_anonymous_usage_stats"] = false

	if fs.Exists(profilePath) {
		file, err := fs.ReadFile(filesystem.NewFileDef(profilePath))
		if err != nil {
			return err
		}
		err = yaml.Unmarshal([]byte(file.Content), &profilesFile)
		if err != nil {
			return fmt.Errorf(`profiles file "%s" is not valid yaml: %w`, profilePath, err)
		}
	}
	profilesFile[project.Profile()] = profileDetails

	yamlEnc, err := yaml.Marshal(&profilesFile)
	if err != nil {
		return err
	}
	err = fs.WriteFile(filesystem.NewRawFile(profilePath, string(yamlEnc)))
	if err != nil {
		return err
	}

	d.Logger().Infof(`Profile stored in "%s".`, profilePath)
	return nil
}
