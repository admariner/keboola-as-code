package dependencies

import (
	"context"
	"fmt"
	stdLog "log"
	"net/url"
	"strings"

	"github.com/keboola/keboola-as-code/internal/pkg/api/client/storageapi"
	"github.com/keboola/keboola-as-code/internal/pkg/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/env"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem/aferofs"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/template/repository"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/strhelper"
	loadRepositoryManifest "github.com/keboola/keboola-as-code/pkg/lib/operation/template/local/repository/manifest/load"
)

type ctxKey string

const CtxKey = ctxKey("dependencies")

// Container provides dependencies used only in the API + common dependencies.
type Container interface {
	dependencies.Common
	Ctx() context.Context
	CtxCancelFn() context.CancelFunc
	WithCtx(ctx context.Context, cancelFn context.CancelFunc) Container
	PrefixLogger() log.PrefixLogger
	RepositoryManager() (*repository.Manager, error)
	DefaultRepository() model.TemplateRepository
	TemplateRepository(definition model.TemplateRepository, forTemplate model.TemplateRef) (*repository.Repository, error)
	WithLoggerPrefix(prefix string) *container
	WithStorageApi(api *storageapi.Api) (*container, error)
}

// NewContainer returns dependencies for API and add them to the context.
func NewContainer(ctx context.Context, repoPath string, debug bool, logger *stdLog.Logger, envs *env.Map) (Container, error) {
	ctx, cancel := context.WithCancel(ctx)
	c := &container{ctx: ctx, ctxCancelFn: cancel, repositoryPath: repoPath, debug: debug, envs: envs, logger: log.NewApiLogger(logger, "", debug)}
	c.commonDeps = dependencies.NewCommonContainer(c)
	repo, err := defaultRepository(repoPath)
	if err != nil {
		return c, err
	}
	c.defaultRepository = repo
	return c, nil
}

type commonDeps = dependencies.CommonContainer

type container struct {
	*commonDeps
	ctx               context.Context
	ctxCancelFn       context.CancelFunc
	debug             bool
	logger            log.PrefixLogger
	envs              *env.Map
	repositoryManager *repository.Manager
	storageApi        *storageapi.Api
	repositoryPath    string
	defaultRepository model.TemplateRepository
}

func (v *container) Ctx() context.Context {
	return v.ctx
}

func (v *container) CtxCancelFn() context.CancelFunc {
	return v.ctxCancelFn
}

func (v *container) DefaultRepository() model.TemplateRepository {
	return v.defaultRepository
}

func defaultRepository(repositoryPath string) (model.TemplateRepository, error) {
	_, err := url.ParseRequestURI(repositoryPath)
	if err != nil {
		return model.TemplateRepository{}, fmt.Errorf("cannot parse repository-path parameter: " + err.Error())
	}

	if strings.HasPrefix(repositoryPath, "file://") {
		return model.TemplateRepository{
			Type: model.RepositoryTypeDir,
			Name: repository.DefaultTemplateRepositoryName,
			Url:  strings.TrimPrefix(repositoryPath, "file://"),
		}, nil
	}

	if strings.HasPrefix(repositoryPath, "https://") {
		uri := strings.TrimPrefix(repositoryPath, "https://")
		refDividerPos := strings.Index(uri, ":")
		if refDividerPos < 0 {
			return model.TemplateRepository{}, fmt.Errorf("invalid repository path url, no ref provided")
		}
		return model.TemplateRepository{
			Type: model.RepositoryTypeGit,
			Name: repository.DefaultTemplateRepositoryName,
			Url:  "https://" + uri[:refDividerPos],
			Ref:  uri[refDividerPos+1:],
		}, nil
	}
	return model.TemplateRepository{}, fmt.Errorf("invalid repository path url provided")
}

func (v *container) WithCtx(ctx context.Context, cancelFn context.CancelFunc) Container {
	clone := v.Clone()
	clone.ctx = ctx
	if cancelFn != nil {
		clone.ctxCancelFn = cancelFn
	}
	return clone
}

// WithLoggerPrefix returns dependencies clone with modified logger.
func (v *container) WithLoggerPrefix(prefix string) *container {
	clone := v.Clone()
	clone.logger = v.logger.WithAdditionalPrefix(prefix)
	return clone
}

// WithStorageApi returns dependencies clone with modified Storage API.
func (v *container) WithStorageApi(api *storageapi.Api) (*container, error) {
	clone := v.Clone()
	clone.storageApi = api
	clone.commonDeps = clone.commonDeps.WithStorageApi(api)
	return clone, nil
}

func (v *container) Logger() log.Logger {
	return v.logger
}

func (v *container) PrefixLogger() log.PrefixLogger {
	return v.logger
}

func (v *container) RepositoryManager() (*repository.Manager, error) {
	if v.repositoryManager == nil {
		if manager, err := repository.NewManager(v.Ctx(), v.Logger(), v.defaultRepository); err != nil {
			return nil, err
		} else {
			v.repositoryManager = manager
		}
	}
	return v.repositoryManager, nil
}

func (v *container) TemplateRepository(definition model.TemplateRepository, _ model.TemplateRef) (*repository.Repository, error) {
	var fs filesystem.Fs
	var err error
	if definition.Type == model.RepositoryTypeGit {
		// Get manager
		manager, err := v.RepositoryManager()
		if err != nil {
			return nil, err
		}

		// Get git repository
		gitRepository, err := manager.Repository(definition)
		if err != nil {
			return nil, err
		}

		// Acquire read lock and release it after request,
		// so pull cannot occur in the middle of the request.
		gitRepository.RLock()
		go func() {
			<-v.ctx.Done()
			gitRepository.RUnlock()
		}()
		fs = gitRepository.Fs()
	} else {
		fs, err = aferofs.NewLocalFs(v.logger, definition.Url, ".")
		if err != nil {
			return nil, err
		}
	}

	// Load manifest from FS
	manifest, err := loadRepositoryManifest.Run(fs, v)
	if err != nil {
		return nil, err
	}

	// Return repository
	return repository.New(definition, fs, manifest)
}

func (v *container) Envs() *env.Map {
	return v.envs
}

func (v *container) ApiVerboseLogs() bool {
	return v.debug
}

func (v *container) StorageApi() (*storageapi.Api, error) {
	// Store API instance, so it can be cloned, see WithStorageApi
	if v.storageApi == nil {
		api, err := v.commonDeps.StorageApi()
		if err != nil {
			return nil, err
		}
		v.storageApi = api
	}

	return v.storageApi, nil
}

func (v *container) StorageApiHost() (string, error) {
	return strhelper.NormalizeHost(v.envs.MustGet("KBC_STORAGE_API_HOST")), nil
}

func (v *container) StorageApiToken() (string, error) {
	// The API is authorized separately in each request
	return "", nil
}

func (v *container) Clone() *container {
	clone := *v
	clone.commonDeps = clone.commonDeps.Clone()
	clone.commonDeps.Abstract = &clone
	return &clone
}
