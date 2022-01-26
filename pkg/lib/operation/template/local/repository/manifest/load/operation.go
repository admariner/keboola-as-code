package load

import (
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	repositoryManifest "github.com/keboola/keboola-as-code/internal/pkg/template/repository/manifest"
)

type dependencies interface {
	Logger() log.Logger
	TemplateRepositoryDir() (filesystem.Fs, error)
}

func Run(d dependencies) (*repositoryManifest.Manifest, error) {
	logger := d.Logger()

	fs, err := d.TemplateRepositoryDir()
	if err != nil {
		return nil, err
	}

	m, err := repositoryManifest.Load(fs)
	if err != nil {
		return nil, err
	}

	logger.Debugf(`ProjectManifest loaded.`)
	return m, nil
}