package dependencies

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/keboola/keboola-sdk-go/v2/pkg/keboola"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/distlock"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/httpclient"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/servicectx"
	aggregationRepo "github.com/keboola/keboola-as-code/internal/pkg/service/stream/aggregation/repository"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/config"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/definition"
	definitionRepo "github.com/keboola/keboola-as-code/internal/pkg/service/stream/definition/repository"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/plugin"
	keboolaSinkBridge "github.com/keboola/keboola-as-code/internal/pkg/service/stream/sink/type/tablesink/keboola/bridge"
	keboolaBridgeRepo "github.com/keboola/keboola-as-code/internal/pkg/service/stream/sink/type/tablesink/keboola/bridge/model/repository"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/model"
	storageRepo "github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/model/repository"
	statsRepo "github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/statistics/repository"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/test/dummy"
	"github.com/keboola/keboola-as-code/internal/pkg/telemetry"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/ulid"
)

const (
	userAgent         = "keboola-stream"
	exceptionIDPrefix = "keboola-stream-task-"
)

// serviceScope implements ServiceScope interface.
type serviceScope struct {
	dependencies.BaseScope
	dependencies.PublicScope
	dependencies.EtcdClientScope
	dependencies.EncryptionScope
	dependencies.DistributedLockScope
	logger                      log.Logger
	plugins                     *plugin.Plugins
	definitionRepository        *definitionRepo.Repository
	storageRepository           *storageRepo.Repository
	storageStatisticsRepository *statsRepo.Repository
	aggregationRepository       *aggregationRepo.Repository
	keboolaBridge               *keboolaSinkBridge.Bridge
	keboolaBridgeRepository     *keboolaBridgeRepo.Repository
	watchTelemetryInterval      time.Duration
}

type parentScopes struct {
	dependencies.BaseScope
	dependencies.PublicScope
	dependencies.EtcdClientScope
	dependencies.EncryptionScope
	dependencies.DistributionScope
	dependencies.DistributedLockScope
	WatchTelemetryInterval time.Duration
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
	ctx, span := tel.Tracer().Start(ctx, "keboola.go.stream.dependencies.NewServiceScope")
	defer span.End(&err)

	p, err := newParentScopes(ctx, cfg, proc, logger, tel, stdout, stderr)
	if err != nil {
		return nil, err
	}

	return newServiceScope(p.BaseScope, p.PublicScope, p.EtcdClientScope, p.EncryptionScope, p.DistributedLockScope, cfg, model.DefaultBackoff())
}

func newParentScopes(
	ctx context.Context,
	cfg config.Config,
	proc *servicectx.Process,
	logger log.Logger,
	tel telemetry.Telemetry,
	stdout io.Writer,
	stderr io.Writer,
) (v *parentScopes, err error) {
	ctx, span := tel.Tracer().Start(ctx, "keboola.go.stream.dependencies.newParentScopes")
	defer span.End(&err)

	// Create base HTTP client for all API requests to other APIs
	httpClient := httpclient.New(
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

	d := &parentScopes{}

	d.BaseScope = dependencies.NewBaseScope(
		ctx,
		logger,
		tel,
		stdout,
		stderr,
		clockwork.NewRealClock(),
		proc,
		httpClient,
		ulid.NewDefaultGenerator(),
	)

	d.PublicScope, err = dependencies.NewPublicScope(ctx, d, cfg.StorageAPIHost, dependencies.WithLogIndexLoading(true))
	if err != nil {
		return nil, err
	}

	d.EtcdClientScope, err = dependencies.NewEtcdClientScope(ctx, d, cfg.Etcd)
	if err != nil {
		return nil, err
	}

	d.DistributedLockScope, err = dependencies.NewDistributedLockScope(ctx, distlock.NewConfig(), d)
	if err != nil {
		return nil, err
	}

	d.EncryptionScope, err = dependencies.NewEncryptionScope(ctx, cfg.Encryption, d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func NewMockedServiceScope(tb testing.TB, ctx context.Context, opts ...dependencies.MockedOption) (ServiceScope, Mocked) {
	tb.Helper()
	return NewMockedServiceScopeWithConfig(tb, ctx, nil, opts...)
}

func NewMockedServiceScopeWithConfig(tb testing.TB, ctx context.Context, modifyConfig func(*config.Config), opts ...dependencies.MockedOption) (ServiceScope, Mocked) {
	tb.Helper()

	// Create common mocked dependencies
	commonMock := dependencies.NewMocked(tb, ctx, append(
		[]dependencies.MockedOption{
			dependencies.WithEnabledEtcdClient(),
			dependencies.WithMockedStorageAPIHost("connection.keboola.local"),
		},
		opts...,
	)...)

	// Get and modify test config
	cfg := testConfig(tb, commonMock)
	if modifyConfig != nil {
		modifyConfig(&cfg)
	}

	// Create service mocked dependencies
	mock := &mocked{Mocked: commonMock, config: cfg, dummySinkController: dummy.NewController()}

	distLockScope, err := dependencies.NewDistributedLockScope(ctx, distlock.NewConfig(), mock)
	require.NoError(tb, err)

	encryptionScope, err := dependencies.NewEncryptionScope(ctx, cfg.Encryption, mock)
	require.NoError(tb, err)

	backoff := model.NoRandomizationBackoff()
	serviceScp, err := newServiceScope(mock, mock, mock, encryptionScope, distLockScope, cfg, backoff)
	require.NoError(tb, err)

	mock.DebugLogger().Truncate()
	if !commonMock.UseRealAPIs() {
		mock.MockedHTTPTransport().Reset()
	}

	mock.dummySinkController.RegisterDummySinkTypes(serviceScp.Plugins(), mock.TestDummySinkController())

	return serviceScp, mock
}

func newServiceScope(
	baseScp dependencies.BaseScope,
	publicScp dependencies.PublicScope,
	etcdClientScp dependencies.EtcdClientScope,
	encryptionScp dependencies.EncryptionScope,
	distLockScp dependencies.DistributedLockScope,
	cfg config.Config,
	storageBackoff model.RetryBackoff,
) (ServiceScope, error) {
	var err error

	d := &serviceScope{}

	d.BaseScope = baseScp

	d.PublicScope = publicScp

	d.EtcdClientScope = etcdClientScp

	d.DistributedLockScope = distLockScp

	d.EncryptionScope = encryptionScp

	d.watchTelemetryInterval = cfg.WatchTelemetryInterval

	d.logger = baseScp.Logger().With(attribute.String("nodeId", cfg.NodeID))

	d.plugins = plugin.New(d.Logger())

	d.definitionRepository = definitionRepo.New(d)

	d.storageRepository, err = storageRepo.New(cfg.Storage.Level, d, storageBackoff)
	if err != nil {
		return nil, err
	}

	d.plugins.RegisterSinkWithLocalStorage(func(sinkType definition.SinkType) bool {
		return sinkType == definition.SinkTypeTable
	})

	apiCtxProvider := func(ctx context.Context) *keboola.AuthorizedAPI {
		api, _ := ctx.Value(KeboolaProjectAPICtxKey).(*keboola.AuthorizedAPI)
		return api
	}

	d.keboolaBridgeRepository, err = keboolaBridgeRepo.New(cfg.Storage.Level, d)
	if err != nil {
		return nil, err
	}

	d.keboolaBridge, err = keboolaSinkBridge.New(d, apiCtxProvider, cfg.Sink.Table.Keboola)
	if err != nil {
		return nil, err
	}

	d.storageStatisticsRepository = statsRepo.New(d)

	d.aggregationRepository = aggregationRepo.New(d)

	return d, nil
}

func (v *serviceScope) Logger() log.Logger {
	return v.logger
}

func (v *serviceScope) Plugins() *plugin.Plugins {
	return v.plugins
}

func (v *serviceScope) DefinitionRepository() *definitionRepo.Repository {
	return v.definitionRepository
}

func (v *serviceScope) KeboolaSinkBridge() *keboolaSinkBridge.Bridge {
	return v.keboolaBridge
}

func (v *serviceScope) KeboolaBridgeRepository() *keboolaBridgeRepo.Repository {
	return v.keboolaBridgeRepository
}

func (v *serviceScope) StorageRepository() *storageRepo.Repository {
	return v.storageRepository
}

func (v *serviceScope) StatisticsRepository() *statsRepo.Repository {
	return v.storageStatisticsRepository
}

func (v *serviceScope) AggregationRepository() *aggregationRepo.Repository {
	return v.aggregationRepository
}

func (v *serviceScope) WatchTelemetryInterval() time.Duration {
	return v.watchTelemetryInterval
}
