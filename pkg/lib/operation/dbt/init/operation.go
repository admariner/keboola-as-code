package init

import (
	"context"
	"fmt"
	"time"

	"github.com/keboola/go-client/pkg/client"
	"github.com/keboola/go-client/pkg/sandboxesapi"
	"github.com/keboola/go-client/pkg/storageapi"
	"go.opentelemetry.io/otel/trace"

	"github.com/keboola/keboola-as-code/internal/pkg/dbt"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/telemetry"
	"github.com/keboola/keboola-as-code/pkg/lib/operation/dbt/generate/env"
	"github.com/keboola/keboola-as-code/pkg/lib/operation/dbt/generate/profile"
	"github.com/keboola/keboola-as-code/pkg/lib/operation/dbt/generate/sources"
)

type DbtInitOptions struct {
	TargetName    string
	WorkspaceName string
}

type dependencies interface {
	JobsQueueApiClient() client.Sender
	Logger() log.Logger
	Tracer() trace.Tracer
	LocalDbtProject(ctx context.Context) (*dbt.Project, bool, error)
	SandboxesApiClient() client.Sender
	StorageApiClient() client.Sender
}

func Run(ctx context.Context, opts DbtInitOptions, d dependencies) (err error) {
	ctx, span := d.Tracer().Start(ctx, "kac.lib.operation.dbt.init")
	defer telemetry.EndSpan(span, &err)

	// Check that we are in dbt directory
	if _, _, err := d.LocalDbtProject(ctx); err != nil {
		return err
	}

	branch, err := storageapi.GetDefaultBranchRequest().Send(ctx, d.StorageApiClient())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	d.Logger().Info(`Creating new workspace, please wait.`)
	// Create workspace
	s, err := sandboxesapi.Create(
		ctx,
		d.StorageApiClient(),
		d.JobsQueueApiClient(),
		d.SandboxesApiClient(),
		branch.ID,
		opts.WorkspaceName,
		sandboxesapi.TypeSnowflake,
	)
	if err != nil {
		return fmt.Errorf("cannot create workspace: %w", err)
	}
	d.Logger().Infof(`Created new workspace "%s".`, opts.WorkspaceName)

	workspace := s.Sandbox

	// Generate profile
	err = profile.Run(ctx, opts.TargetName, d)
	if err != nil {
		return fmt.Errorf("could not generate profile: %w", err)
	}

	// Generate sources
	err = sources.Run(ctx, opts.TargetName, d)
	if err != nil {
		return fmt.Errorf("could not generate sources: %w", err)
	}

	// Generate env
	err = env.Run(ctx, env.GenerateEnvOptions{
		TargetName: opts.TargetName,
		Workspace:  workspace,
	}, d)
	if err != nil {
		return fmt.Errorf("could not generate env: %w", err)
	}

	return nil
}