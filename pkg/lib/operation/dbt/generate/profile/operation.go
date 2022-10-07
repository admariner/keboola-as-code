package profile

import (
	"context"
	"fmt"
	"strings"

	"github.com/keboola/go-utils/pkg/orderedmap"
	"go.opentelemetry.io/otel/trace"

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

	// Load profiles file if exists
	profilesFile := orderedmap.New()
	profilesFileDef := filesystem.NewFileDef(dbt.ProfilesPath).SetDescription("dbt profiles")
	if fs.Exists(dbt.ProfilesPath) {
		if _, err := fs.FileLoader().ReadYamlFileTo(profilesFileDef, profilesFile); err != nil {
			return err
		}
	}

	// Set send_anonymous_usage_stats if missing
	if _, found := profilesFile.Get("send_anonymous_usage_stats"); !found {
		profilesFile.Set("send_anonymous_usage_stats", false)
	}

	// Set profile
	targetUpper := strings.ToUpper(targetName)
	profilesFile.Set(project.Profile(), orderedmap.FromPairs([]orderedmap.Pair{
		{
			Key:   "target",
			Value: targetName,
		},
		{
			Key: "outputs",
			Value: orderedmap.FromPairs([]orderedmap.Pair{
				{
					Key: targetName,
					Value: orderedmap.FromPairs([]orderedmap.Pair{
						{
							Key:   "account",
							Value: fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_ACCOUNT\") }}", targetUpper),
						},
						{
							Key:   "database",
							Value: fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_DATABASE\") }}", targetUpper),
						},
						{
							Key:   "password",
							Value: fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_PASSWORD\") }}", targetUpper),
						},
						{
							Key:   "schema",
							Value: fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_SCHEMA\") }}", targetUpper),
						},
						{
							Key:   "type",
							Value: fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_TYPE\") }}", targetUpper),
						},
						{
							Key:   "user",
							Value: fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_USER\") }}", targetUpper),
						},
						{
							Key:   "warehouse",
							Value: fmt.Sprintf("{{ env_var(\"DBT_KBC_%s_WAREHOUSE\") }}", targetUpper),
						},
					}),
				},
			}),
		},
	}))

	// Save file
	if err := fs.WriteFile(filesystem.NewYamlFile(dbt.ProfilesPath, profilesFile).SetDescription("dbt profiles")); err != nil {
		return err
	}

	d.Logger().Infof(`Profile stored in "%s".`, profilePath)
	return nil
}