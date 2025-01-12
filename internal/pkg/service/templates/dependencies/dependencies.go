// Package dependencies provides dependencies for Templates API.
//
// # Dependency Containers
//
// This package extends common dependencies from [pkg/github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies].
//
// Following dependencies containers are implemented:
//   - [APIScope] long-lived dependencies that exist during the entire run of the API server.
//   - [PublicRequestScope] short-lived dependencies for a public request without authentication.
//   - [ProjectRequestScope] short-lived dependencies for a request with authentication.
//
// Dependency containers creation:
//   - [APIScope] is created at startup in main.go.
//   - [PublicRequestScope] is created for each HTTP request by Muxer.Use callback in main.go.
//   - [ProjectRequestScope] is created for each authenticated HTTP request in the service.APIKeyAuth method.
//
// The package also provides mocked dependency implementations for tests:
//   - [NewMockedAPIScope]
//   - [NewMockedPublicRequestScope]
//   - [NewMockedProjectRequestScope]
//
// Dependencies injection to service endpoints:
//   - Each service endpoint method gets [PublicRequestScope] as a parameter.
//   - Authorized endpoints gets [ProjectRequestScope] instead.
//   - The injection is generated by "internal/pkg/service/common/goaextension/dependencies" package.
//   - See service implementation for details [src/github.com/keboola/keboola-as-code/internal/pkg/service/biffer/api/service/service.go].
package dependencies

import (
	"context"

	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/templates/api/config"
	"github.com/keboola/keboola-as-code/internal/pkg/service/templates/store"
	"github.com/keboola/keboola-as-code/internal/pkg/service/templates/store/schema"
	"github.com/keboola/keboola-as-code/internal/pkg/template"
	"github.com/keboola/keboola-as-code/internal/pkg/template/repository"
	repositoryManager "github.com/keboola/keboola-as-code/internal/pkg/template/repository/manager"
)

type ctxKey string

const (
	PublicRequestScopeCtxKey  = ctxKey("PublicRequestScope")
	ProjectRequestScopeCtxKey = ctxKey("ProjectRequestScope")
)

// APIScope interface provides dependencies for Templates API server.
// The container exists during the entire run of the API server.
type APIScope interface {
	dependencies.BaseScope
	dependencies.PublicScope
	dependencies.EtcdClientScope
	dependencies.TaskScope
	dependencies.DistributionScope
	APIConfig() config.Config
	Schema() *schema.Schema
	Store() *store.Store
	RepositoryManager() *repositoryManager.Manager
	ProjectLocker() *Locker
}

// PublicRequestScope interface provides dependencies for a public request that doesn't need the Storage API token.
// The container exists only during request processing.
type PublicRequestScope interface {
	APIScope
	dependencies.RequestInfo
}

// ProjectRequestScope interface provides dependencies for an request authenticated by the Storage API token.
// The container exists only during request processing.
type ProjectRequestScope interface {
	PublicRequestScope
	dependencies.ProjectScope
	Template(ctx context.Context, reference model.TemplateRef) (*template.Template, error)
	TemplateRepository(ctx context.Context, reference model.TemplateRepository) (*repository.Repository, error)
	ProjectRepositories() *model.TemplateRepositories
}
