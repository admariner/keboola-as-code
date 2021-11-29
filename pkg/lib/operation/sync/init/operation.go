package init

import (
	"context"

	"go.uber.org/zap"

	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/manifest"
	"github.com/keboola/keboola-as-code/internal/pkg/remote"
	"github.com/keboola/keboola-as-code/internal/pkg/scheduler"
	"github.com/keboola/keboola-as-code/internal/pkg/state"
	"github.com/keboola/keboola-as-code/internal/pkg/utils"
	createEnvFiles "github.com/keboola/keboola-as-code/pkg/lib/operation/local/envfiles/create"
	createManifest "github.com/keboola/keboola-as-code/pkg/lib/operation/local/manifest/create"
	createMetaDir "github.com/keboola/keboola-as-code/pkg/lib/operation/local/metadir/create"
	genWorkflows "github.com/keboola/keboola-as-code/pkg/lib/operation/local/workflows/generate"
	loadState "github.com/keboola/keboola-as-code/pkg/lib/operation/state/load"
	"github.com/keboola/keboola-as-code/pkg/lib/operation/sync/pull"
)

type Options struct {
	Pull            bool // run pull after init?
	ManifestOptions createManifest.Options
	Workflows       genWorkflows.Options
}

type dependencies interface {
	Ctx() context.Context
	Logger() *zap.SugaredLogger
	Fs() filesystem.Fs
	AssertMetaDirNotExists() error
	StorageApi() (*remote.StorageApi, error)
	SchedulerApi() (*scheduler.Api, error)
	CreateManifest(o createManifest.Options) (*manifest.Manifest, error)
	Manifest() (*manifest.Manifest, error)
	LoadStateOnce(loadOptions loadState.Options) (*state.State, error)
}

func Run(o Options, d dependencies) (err error) {
	logger := d.Logger()

	// Is project directory already initialized?
	if err := d.AssertMetaDirNotExists(); err != nil {
		return err
	}

	// Create metadata dir
	if err := createMetaDir.Run(d); err != nil {
		return err
	}

	// Create manifest
	if _, err := d.CreateManifest(o.ManifestOptions); err != nil {
		return err
	}

	// Create ENV files
	if err := createEnvFiles.Run(d); err != nil {
		return err
	}

	// Related operations
	errors := utils.NewMultiError()

	// Generate CI workflows
	if err := genWorkflows.Run(o.Workflows, d); err != nil {
		errors.Append(utils.PrefixError(`workflows generation failed`, err))
	}

	logger.Info("Init done.")

	// First pull
	if o.Pull {
		logger.Info()
		logger.Info(`Running pull.`)
		pullOptions := pull.Options{
			DryRun:            false,
			Force:             false,
			LogUntrackedPaths: false,
		}
		if err := pull.Run(pullOptions, d); err != nil {
			errors.Append(utils.PrefixError(`pull failed`, err))
		}
	}

	return errors.ErrorOrNil()
}
