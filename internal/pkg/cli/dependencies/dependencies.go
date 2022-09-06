// Package dependencies provides dependencies for command line interface.
//
// # Dependency Containers
//
// This package extends common dependencies from [pkg/github.com/keboola/keboola-as-code/internal/pkg/dependencies].
//
// These dependencies containers are implemented:
//   - [Base] interface provides basic CLI dependencies.
//   - [ForLocalCommand] interface provides dependencies for commands that do not modify the remote project
//   - [ForRemoteCommand] interface provides dependencies for commands that modify the remote project.
//
// These containers can be obtained from the [Provider], it can be created by [NewProvider].
package dependencies

import (
	"context"
	"fmt"

	"github.com/keboola/keboola-as-code/internal/pkg/cli/dialog"
	"github.com/keboola/keboola-as-code/internal/pkg/cli/options"
	"github.com/keboola/keboola-as-code/internal/pkg/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	projectPkg "github.com/keboola/keboola-as-code/internal/pkg/project"
	"github.com/keboola/keboola-as-code/internal/pkg/template"
	"github.com/keboola/keboola-as-code/internal/pkg/template/repository"
)

var (
	ErrMissingStorageApiHost          = fmt.Errorf(`missing Storage API host`)
	ErrMissingStorageApiToken         = fmt.Errorf(`missing Storage API token`)
	ErrInvalidStorageApiToken         = fmt.Errorf(`invalid Storage API token`)
	ErrProjectDirFound                = fmt.Errorf("remote directory not expected, but found")
	ErrProjectManifestNotFound        = fmt.Errorf("remote manifest not found")
	ErrTemplateDirFound               = fmt.Errorf("template directory not expected, but found")
	ErrTemplateManifestNotFound       = fmt.Errorf("template manifest not found")
	ErrRepositoryDirFound             = fmt.Errorf("repository directory not expected, but found")
	ErrRepositoryManifestNotFound     = fmt.Errorf("repository manifest not found")
	ErrExpectedRepositoryFoundProject = fmt.Errorf("repository manifest not found, found remote manifest")
	ErrExpectedProjectFoundRepository = fmt.Errorf("remote manifest not found, found repository manifest")
	ErrExpectedProjectFoundTemplate   = fmt.Errorf("remote manifest not found, found template manifest")
)

// Base interface provides basic CLI dependencies.
type Base interface {
	dependencies.Base
	CommandCtx() context.Context
	Fs() filesystem.Fs
	FsInfo() FsInfo
	Dialogs() *dialog.Dialogs
	Options() *options.Options
	EmptyDir() (filesystem.Fs, error)
}

// ForLocalCommand interface provides dependencies for commands that do not modify the remote project.
// It contains CLI dependencies that are available from the Storage API and other sources without authentication / Storage API token.
type ForLocalCommand interface {
	Base
	dependencies.Public
	LocalProject(ignoreErrors bool) (*projectPkg.Project, bool, error)
	LocalTemplate(ctx context.Context) (*template.Template, bool, error)
	LocalTemplateRepository(ctx context.Context) (*repository.Repository, bool, error)
}

// ForRemoteCommand interface provides dependencies for commands that modify remote project.
// It contains CLI dependencies that require authentication / Storage API token.
type ForRemoteCommand interface {
	ForLocalCommand
	dependencies.Project
}

// Provider of CLI dependencies.
type Provider interface {
	BaseDependencies() Base
	DependenciesForLocalCommand() (ForLocalCommand, error)
	DependenciesForRemoteCommand() (ForRemoteCommand, error)
	// LocalProject method can be used by a CLI command that must be run in the local project directory.
	// First, the local project is loaded, and then the authentication is performed,
	// so the error that we are not in a project directory takes precedence over an invalid/missing token.
	LocalProject(ignoreErrors bool) (*projectPkg.Project, ForRemoteCommand, error)
	// LocalRepository method can be used by a CLI command that must be run in the local repository directory.
	LocalRepository() (*repository.Repository, ForLocalCommand, error)
}