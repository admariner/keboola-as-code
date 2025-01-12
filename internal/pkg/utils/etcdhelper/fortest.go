package etcdhelper

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"         //nolint: depguard
	"go.uber.org/zap/zapcore" //nolint: depguard
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	"github.com/keboola/keboola-as-code/internal/pkg/env"
	"github.com/keboola/keboola-as-code/internal/pkg/idgenerator"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdclient"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdlogger"
)

type testOrBenchmark interface {
	Cleanup(f func())
	Skipf(format string, args ...any)
	Logf(format string, args ...any)
	Fatalf(format string, args ...any)
}

// TmpNamespace creates a temporary etcd namespace and registers cleanup after the test.
func TmpNamespace(t testOrBenchmark) etcdclient.Credentials {
	return TmpNamespaceFromEnv(t, "UNIT_ETCD_")
}

// TmpNamespaceFromEnv creates a temporary etcd namespace and registers cleanup after the test.
// Credentials are read from the provided ENV prefix.
func TmpNamespaceFromEnv(t testOrBenchmark, envPrefix string) etcdclient.Credentials {
	envs, err := env.FromOs()
	if err != nil {
		t.Fatalf("cannot get envs: %s", err)
	}

	if envs.Get(envPrefix+"ENABLED") == "false" {
		t.Skipf(fmt.Sprintf("etcd test is disabled by %s_ENABLED=false", envPrefix))
	}

	credentials := etcdclient.Credentials{
		Endpoint:  envs.Get(envPrefix + "ENDPOINT"),
		Namespace: idgenerator.EtcdNamespaceForTest(),
		Username:  envs.Get(envPrefix + "USERNAME"),
		Password:  envs.Get(envPrefix + "PASSWORD"),
	}

	if credentials.Endpoint == "" {
		t.Fatalf(`etcd endpoint is not set`)
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithCancel(context.Background())
		client := clientForTest(t, ctx, credentials)
		_, err := client.Delete(ctx, "", etcd.WithFromKey())
		cancel()
		if err != nil {
			t.Fatalf(`cannot clear etcd after test: %s`, err)
		}
	})

	return credentials
}

func ClientForTest(t testOrBenchmark, credentials etcdclient.Credentials, dialOpts ...grpc.DialOption) *etcd.Client {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
	})
	return clientForTest(t, ctx, credentials, dialOpts...)
}

func clientForTest(t testOrBenchmark, ctx context.Context, credentials etcdclient.Credentials, dialOpts ...grpc.DialOption) *etcd.Client {
	// Normalize namespace
	credentials.Namespace = strings.Trim(credentials.Namespace, " /") + "/"

	// Setup logger
	var logger *zap.Logger

	// Should be logger enabled?
	verbose := VerboseTestLogs()

	// Enable logger
	if verbose {
		encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		logger = zap.New(log.NewCallbackCore(func(entry zapcore.Entry, fields []zapcore.Field) {
			if entry.Level > log.DebugLevel {
				bytes, _ := encoder.EncodeEntry(entry, fields)
				_, _ = os.Stdout.Write(bytes.Bytes())
			}
		}))
	}

	// Dial options
	dialOpts = append(
		dialOpts,
		grpc.WithBlock(),                 // wait for the connection
		grpc.WithReturnConnectionError(), // wait for the connection error
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				Multiplier: 1.5,
				Jitter:     0.2,
				MaxDelay:   15 * time.Second,
			},
		}),
	)

	// Create etcd client
	etcdClient, err := etcd.New(etcd.Config{
		Context:              ctx,
		Endpoints:            []string{credentials.Endpoint},
		DialTimeout:          10 * time.Second,
		DialKeepAliveTimeout: 5 * time.Second,
		DialKeepAliveTime:    10 * time.Second,
		Username:             credentials.Username, // optional
		Password:             credentials.Password, // optional
		Logger:               logger,
		DialOptions:          dialOpts,
	})
	if err != nil {
		t.Fatalf("cannot create etcd client: %s, %s", err, debug.Stack())
	}

	// Use namespace
	etcdclient.UseNamespace(etcdClient, credentials.Namespace)

	// Add operations logger
	if verbose {
		etcdClient.KV = etcdlogger.KVLogWrapper(etcdClient.KV, os.Stdout)
	}

	return etcdClient
}

func VerboseTestLogs() bool {
	str, found := os.LookupEnv("ETCD_VERBOSE")
	return found && strings.ToLower(str) == "true"
}
