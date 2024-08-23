package testnode

import (
	"context"
	"testing"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	commonDeps "github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdclient"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/config"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/dependencies"
)

func StartSourceNode(tb testing.TB, ctx context.Context, logger log.DebugLogger, etcdCfg etcdclient.Config, modifyConfig func(cfg *config.Config), opts ...commonDeps.MockedOption) (dependencies.SourceScope, dependencies.Mocked) {
	tb.Helper()
	opts = append(opts, commonDeps.WithDebugLogger(logger), commonDeps.WithEtcdConfig(etcdCfg))
	return dependencies.NewMockedSourceScopeWithConfig(
		tb,
		ctx,
		func(cfg *config.Config) {
			if modifyConfig != nil {
				modifyConfig(cfg)
			}
			cfg.NodeID = "source"
		},
		opts...,
	)
}
