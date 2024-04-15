// Package dependencies provides dependencies for Apps Proxy.
//
// # Dependency Containers
//
// This package extends common dependencies from [pkg/github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies].
//
// Following dependencies containers are implemented:
//   - [ServiceScope] long-lived dependencies that exist during the entire run of the proxy server.
//
// Dependency containers creation:
//   - [ServiceScope] is created at startup in main.go.
//
// The package also provides mocked dependency implementations for tests:
//   - [NewMockedServiceScope]
package dependencies

import (
	"context"
	"io"

	"github.com/benbjohnson/clock"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/appsproxy/config"
	"github.com/keboola/keboola-as-code/internal/pkg/service/appsproxy/dataapps/api"
	"github.com/keboola/keboola-as-code/internal/pkg/service/appsproxy/dataapps/appconfig"
	"github.com/keboola/keboola-as-code/internal/pkg/service/appsproxy/dataapps/notify"
	"github.com/keboola/keboola-as-code/internal/pkg/service/appsproxy/dataapps/wakeup"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/httpclient"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/servicectx"
	"github.com/keboola/keboola-as-code/internal/pkg/telemetry"
)

const (
	userAgent = "keboola-apps-proxy"
)

// ServiceScope interface provides dependencies for Apps Proxy server.
// The container exists during the entire run of the Apps Proxy server.
type ServiceScope interface {
	dependencies.BaseScope
	Config() config.Config
	AppsAPI() *api.API
	AppConfigLoader() *appconfig.Loader
	NotifyManager() *notify.Manager
	WakeupManager() *wakeup.Manager
}

type Mocked interface {
	dependencies.Mocked
	TestConfig() config.Config
}

// serviceScope implements APIScope interface.
type serviceScope struct {
	parentScopes
	config          config.Config
	appsAPI         *api.API
	appConfigLoader *appconfig.Loader
	notifyManager   *notify.Manager
	wakeupManager   *wakeup.Manager
}

type parentScopes interface {
	dependencies.BaseScope
}

type parentScopesImpl struct {
	dependencies.BaseScope
}

func NewServiceScope(
	ctx context.Context,
	cfg config.Config,
	proc *servicectx.Process,
	logger log.Logger,
	tel telemetry.Telemetry,
	stdout io.Writer,
	stderr io.Writer,
) (v ServiceScope, err error) {
	parentScp := newParentScopes(ctx, cfg, proc, logger, tel, stdout, stderr)
	return newServiceScope(ctx, parentScp, cfg)
}

func newParentScopes(
	ctx context.Context,
	cfg config.Config,
	proc *servicectx.Process,
	logger log.Logger,
	tel telemetry.Telemetry,
	stdout io.Writer,
	stderr io.Writer,
) (v parentScopes) {
	ctx, span := tel.Tracer().Start(ctx, "keboola.go.appsproxy.dependencies.newParentScopes")
	defer span.End(nil)

	httpClient := httpclient.New(
		httpclient.WithoutForcedHTTP2(), // We're currently unable to connect to Sandboxes Service using HTTP2.
		httpclient.WithTelemetry(tel),
		httpclient.WithUserAgent(userAgent),
		func(c *httpclient.Config) {
			if cfg.DebugLog {
				httpclient.WithDebugOutput(stdout)(c)
			}
			if cfg.DebugHTTPClient {
				httpclient.WithDumpOutput(stdout)(c)
			}
		},
	)

	d := &parentScopesImpl{}
	d.BaseScope = dependencies.NewBaseScope(ctx, logger, tel, stdout, stderr, clock.New(), proc, httpClient)
	return d
}

func newServiceScope(ctx context.Context, parentScp parentScopes, cfg config.Config) (v *serviceScope, err error) {
	ctx, span := parentScp.Telemetry().Tracer().Start(ctx, "keboola.go.appsproxy.dependencies.newServiceScope")
	defer span.End(&err)

	d := &serviceScope{}
	d.parentScopes = parentScp
	d.config = cfg
	d.appsAPI = api.New(d.HTTPClient(), cfg.SandboxesAPI.URL, cfg.SandboxesAPI.Token)
	d.appConfigLoader = appconfig.NewLoader(d)
	d.notifyManager = notify.NewManager(d)
	d.wakeupManager = wakeup.NewManager(d)

	return d, nil
}

func (v *serviceScope) Config() config.Config {
	return v.config
}

func (v *serviceScope) AppsAPI() *api.API {
	return v.appsAPI
}

func (v *serviceScope) AppConfigLoader() *appconfig.Loader {
	return v.appConfigLoader
}

func (v *serviceScope) NotifyManager() *notify.Manager {
	return v.notifyManager
}

func (v *serviceScope) WakeupManager() *wakeup.Manager {
	return v.wakeupManager
}
