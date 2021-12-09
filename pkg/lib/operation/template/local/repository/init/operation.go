package init

import (
	"context"

	"go.uber.org/zap"

	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/template/repository/manifest"
	createMetaDir "github.com/keboola/keboola-as-code/pkg/lib/operation/local/metadir/create"
)

type dependencies interface {
	Ctx() context.Context
	Logger() *zap.SugaredLogger
	EmptyDir() (filesystem.Fs, error)
	CreateRepositoryManifest() (*manifest.Manifest, error)
}

func Run(d dependencies) (err error) {
	logger := d.Logger()

	// Create metadata dir
	if err := createMetaDir.Run(d); err != nil {
		return err
	}

	// Create manifest
	if _, err := d.CreateRepositoryManifest(); err != nil {
		return err
	}

	logger.Info("Repository init done.")
	return nil
}