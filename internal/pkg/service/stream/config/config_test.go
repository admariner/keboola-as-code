package config_test

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/c2h5oh/datasize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/keboola/keboola-as-code/internal/pkg/encoding/json"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/configmap"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/configpatch"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/config"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/level"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/level/local"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/level/local/volume"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/level/local/volume/diskalloc"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/test"
	"github.com/keboola/keboola-as-code/internal/pkg/validator"
)

func TestConfig_DefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := config.New()

	bytes, err := configmap.NewDumper().Dump(&cfg).AsYAML()
	require.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(`
# Enable logging at DEBUG level.
debugLog: false
# Log HTTP client requests and responses as debug messages.
debugHTTPClient: false
# Path where CPU profile is saved.
cpuProfilePath: ""
# Unique ID of the node in the cluster. Validation rules: required
nodeID: ""
# Storage API host. Validation rules: required,hostname
storageApiHost: ""
datadog:
    # Enable DataDog integration.
    enabled: true
    # Enable DataDog debug messages.
    debug: false
etcd:
    # Etcd endpoint. Validation rules: required
    endpoint: ""
    # Etcd namespace. Validation rules: required
    namespace: ""
    # Etcd username.
    username: ""
    # Etcd password.
    password: '*****'
    # Etcd connect timeout. Validation rules: required
    connectTimeout: 30s
    # Etcd keep alive timeout. Validation rules: required
    keepAliveTimeout: 5s
    # Etcd keep alive interval. Validation rules: required
    keepAliveInterval: 10s
    # Etcd operations logging as debug messages.
    debugLog: false
metrics:
    # Prometheus scraping metrics listen address. Validation rules: required,hostname_port
    listen: 0.0.0.0:9000
api:
    # Listen address of the configuration HTTP API. Validation rules: required,hostname_port
    listen: 0.0.0.0:8000
    # Public URL of the configuration HTTP API for link generation. Validation rules: required
    publicUrl: http://localhost:8000
distribution:
    # The maximum time to wait for creating a new session. Validation rules: required,minDuration=1s,maxDuration=1m
    grantTimeout: 5s
    # Timeout for the node registration to the cluster. Validation rules: required,minDuration=1s,maxDuration=5m
    startupTimeout: 1m0s
    # Timeout for the node un-registration from the cluster. Validation rules: required,minDuration=1s,maxDuration=5m
    shutdownTimeout: 10s
    # Interval of processing changes in the topology. Use 0 to disable the grouping. Validation rules: maxDuration=30s
    eventsGroupInterval: 5s
    # Seconds after which the node is automatically un-registered if an outage occurs. Validation rules: required,min=1,max=30
    ttlSeconds: 15
source:
    http:
        # Listen address of the HTTP source. Validation rules: required,hostname_port
        listen: 0.0.0.0:7000
        # Public URL of the HTTP source for link generation.
        publicUrl: null
storage:
    statistics:
        sync:
            # Statistics synchronization interval, from memory to the etcd. Validation rules: required,minDuration=100ms,maxDuration=5s
            interval: 1s
            # Statistics synchronization timeout. Validation rules: required,minDuration=1s,maxDuration=1m
            timeout: 30s
        cache:
            L2:
                # Enable statistics L2 in-memory cache, otherwise only L1 cache is used.
                enabled: true
                # Statistics L2 in-memory cache invalidation interval. Validation rules: required,minDuration=100ms,maxDuration=5s
                interval: 1s
    cleanup:
        # Enable storage cleanup.
        enabled: true
        # Cleanup interval. Validation rules: required,minDuration=5m,maxDuration=24h
        interval: 30m0s
        # How many files are deleted in parallel. Validation rules: required,min=1,max=500
        concurrency: 100
        # Expiration interval of a file that has not yet been imported. Validation rules: required,minDuration=1h,maxDuration=720h,gtefield=ArchivedFileExpiration
        activeFileExpiration: 168h0m0s
        # Expiration interval of a file that has already been imported. Validation rules: required,minDuration=15m,maxDuration=720h
        archivedFileExpiration: 24h0m0s
    level:
        local:
            volume:
                assignment:
                    # Volumes count simultaneously utilized per sink. Validation rules: required,min=1,max=100
                    count: 1
                    # List of preferred volume types, start with the most preferred. Validation rules: required,min=1
                    preferredTypes:
                        - default
                registration:
                    # Number of seconds after the volume registration expires if the node is not available. Validation rules: required,min=1,max=60
                    ttlSeconds: 10
                sync:
                    # Sync mode: "disabled", "cache" or "disk". Validation rules: required,oneof=disabled disk cache
                    mode: disk
                    # Wait for sync to disk OS cache or to disk hardware, depending on the mode.
                    wait: true
                    # Minimal interval between syncs. Validation rules: min=0,maxDuration=2s,required_if=Mode disk,required_if=Mode cache
                    checkInterval: 5ms
                    # Written records count to trigger sync. Validation rules: min=0,max=1000000,required_if=Mode disk,required_if=Mode cache
                    countTrigger: 500
                    # Written size to trigger sync. Validation rules: maxBytes=100MB,required_if=Mode disk,required_if=Mode cache
                    bytesTrigger: 1MB
                    # Interval from the last sync to trigger sync. Validation rules: min=0,maxDuration=2s,required_if=Mode disk,required_if=Mode cache
                    intervalTrigger: 50ms
                allocation:
                    # Allocate disk space for each slice.
                    enabled: true
                    # Size of allocated disk space for a slice. Validation rules: required
                    static: 100MB
                    # Allocate disk space as % from the previous slice size. Validation rules: min=100,max=500
                    relative: 110
            compression:
                # Compression type. Validation rules: required,oneof=none gzip zstd
                type: gzip
                gzip:
                    # GZIP compression level: 1-9. Validation rules: min=1,max=9
                    level: 1
                    # GZIP implementation: standard, fast, parallel. Validation rules: required,oneof=standard fast parallel
                    implementation: parallel
                    # GZIP parallel block size. Validation rules: required,minBytes=16kB,maxBytes=100MB
                    blockSize: 256KB
                    # GZIP parallel concurrency, 0 = auto.
                    concurrency: 0
                zstd:
                    # ZSTD compression level: fastest, default, better, best. Validation rules: min=1,max=4
                    level: 1
                    # ZSTD window size. Validation rules: required,minBytes=1kB,maxBytes=512MB
                    windowSize: 1MB
                    # ZSTD concurrency, 0 = auto
                    concurrency: 0
        staging:
            # Maximum number of slices in a file, a new file is created after reaching it. Validation rules: required,min=1,max=50000
            maxSlicesPerFile: 100
            # Maximum number of the Storage API file resources created in parallel within one operation. Validation rules: required,min=1,max=500
            parallelFileCreateLimit: 50
            upload:
                # Minimal interval between uploads. Validation rules: required,minDuration=1s,maxDuration=5m
                minInterval: 5s
                trigger:
                    # Records count. Validation rules: required,min=1,max=10000000
                    count: 10000
                    # Records size. Validation rules: required,minBytes=100B,maxBytes=50MB
                    size: 1MB
                    # Duration from the last upload. Validation rules: required,minDuration=1s,maxDuration=30m
                    interval: 1m0s
        target:
            import:
                # Minimal interval between imports. Validation rules: required,minDuration=30s,maxDuration=30m
                minInterval: 1m0s
                trigger:
                    # Records count. Validation rules: required,min=1,max=10000000
                    count: 50000
                    # Records size. Validation rules: required,minBytes=100B,maxBytes=500MB
                    size: 5MB
                    # Duration from the last import. Validation rules: required,minDuration=60s,maxDuration=24h
                    interval: 5m0s
`), strings.TrimSpace(string(bytes)))

	// Add missing values, and validate it
	cfg.NodeID = "test-node"
	cfg.StorageAPIHost = "connection.keboola.local"
	cfg.Source.HTTP.PublicURL, _ = url.Parse("https://stream-in.keboola.local")
	cfg.Etcd.Endpoint = "test-etcd"
	cfg.Etcd.Namespace = "test-namespace"
	require.NoError(t, validator.New().Validate(context.Background(), cfg))
}

func TestTableSinkConfigPatch_ToKVs(t *testing.T) {
	t.Parallel()

	kvs, err := configpatch.DumpAll(
		config.New(),
		config.Patch{
			Storage: &storage.ConfigPatch{
				Level: &level.ConfigPatch{
					Local: &local.ConfigPatch{
						Volume: &volume.ConfigPatch{
							Allocation: &diskalloc.ConfigPatch{
								Static: test.Ptr(456 * datasize.MB),
							},
						},
					},
				},
			},
		},
	)
	require.NoError(t, err)

	assert.Equal(t, strings.TrimSpace(`
[
  {
    "key": "storage.level.local.compression.gzip.blockSize",
    "value": "256KB",
    "defaultValue": "256KB",
    "overwritten": false,
    "protected": false,
    "validation": "required,minBytes=16kB,maxBytes=100MB"
  },
  {
    "key": "storage.level.local.compression.gzip.concurrency",
    "value": 0,
    "defaultValue": 0,
    "overwritten": false,
    "protected": false
  },
  {
    "key": "storage.level.local.compression.gzip.implementation",
    "value": "parallel",
    "defaultValue": "parallel",
    "overwritten": false,
    "protected": false,
    "validation": "required,oneof=standard fast parallel"
  },
  {
    "key": "storage.level.local.compression.gzip.level",
    "value": 1,
    "defaultValue": 1,
    "overwritten": false,
    "protected": false,
    "validation": "min=1,max=9"
  },
  {
    "key": "storage.level.local.compression.type",
    "value": "gzip",
    "defaultValue": "gzip",
    "overwritten": false,
    "protected": false,
    "validation": "required,oneof=none gzip zstd"
  },
  {
    "key": "storage.level.local.compression.zstd.concurrency",
    "value": 0,
    "defaultValue": 0,
    "overwritten": false,
    "protected": false
  },
  {
    "key": "storage.level.local.compression.zstd.level",
    "value": 1,
    "defaultValue": 1,
    "overwritten": false,
    "protected": false,
    "validation": "min=1,max=4"
  },
  {
    "key": "storage.level.local.compression.zstd.windowSize",
    "value": "1MB",
    "defaultValue": "1MB",
    "overwritten": false,
    "protected": false,
    "validation": "required,minBytes=1kB,maxBytes=512MB"
  },
  {
    "key": "storage.level.local.volume.allocation.enabled",
    "value": true,
    "defaultValue": true,
    "overwritten": false,
    "protected": false
  },
  {
    "key": "storage.level.local.volume.allocation.relative",
    "value": 110,
    "defaultValue": 110,
    "overwritten": false,
    "protected": false,
    "validation": "min=100,max=500"
  },
  {
    "key": "storage.level.local.volume.allocation.static",
    "value": "456MB",
    "defaultValue": "100MB",
    "overwritten": true,
    "protected": false,
    "validation": "required"
  },
  {
    "key": "storage.level.local.volume.assignment.count",
    "value": 1,
    "defaultValue": 1,
    "overwritten": false,
    "protected": false,
    "validation": "required,min=1,max=100"
  },
  {
    "key": "storage.level.local.volume.assignment.preferredTypes",
    "value": [
      "default"
    ],
    "defaultValue": [
      "default"
    ],
    "overwritten": false,
    "protected": false,
    "validation": "required,min=1"
  },
  {
    "key": "storage.level.local.volume.sync.bytesTrigger",
    "value": "1MB",
    "defaultValue": "1MB",
    "overwritten": false,
    "protected": false,
    "validation": "maxBytes=100MB,required_if=Mode disk,required_if=Mode cache"
  },
  {
    "key": "storage.level.local.volume.sync.checkInterval",
    "value": "5ms",
    "defaultValue": "5ms",
    "overwritten": false,
    "protected": false,
    "validation": "min=0,maxDuration=2s,required_if=Mode disk,required_if=Mode cache"
  },
  {
    "key": "storage.level.local.volume.sync.countTrigger",
    "value": 500,
    "defaultValue": 500,
    "overwritten": false,
    "protected": false,
    "validation": "min=0,max=1000000,required_if=Mode disk,required_if=Mode cache"
  },
  {
    "key": "storage.level.local.volume.sync.intervalTrigger",
    "value": "50ms",
    "defaultValue": "50ms",
    "overwritten": false,
    "protected": false,
    "validation": "min=0,maxDuration=2s,required_if=Mode disk,required_if=Mode cache"
  },
  {
    "key": "storage.level.local.volume.sync.mode",
    "value": "disk",
    "defaultValue": "disk",
    "overwritten": false,
    "protected": false,
    "validation": "required,oneof=disabled disk cache"
  },
  {
    "key": "storage.level.local.volume.sync.wait",
    "value": true,
    "defaultValue": true,
    "overwritten": false,
    "protected": false
  },
  {
    "key": "storage.level.staging.maxSlicesPerFile",
    "value": 100,
    "defaultValue": 100,
    "overwritten": false,
    "protected": false,
    "validation": "required,min=1,max=50000"
  },
  {
    "key": "storage.level.staging.upload.trigger.count",
    "value": 10000,
    "defaultValue": 10000,
    "overwritten": false,
    "protected": false,
    "validation": "required,min=1,max=10000000"
  },
  {
    "key": "storage.level.staging.upload.trigger.interval",
    "value": "1m0s",
    "defaultValue": "1m0s",
    "overwritten": false,
    "protected": false,
    "validation": "required,minDuration=1s,maxDuration=30m"
  },
  {
    "key": "storage.level.staging.upload.trigger.size",
    "value": "1MB",
    "defaultValue": "1MB",
    "overwritten": false,
    "protected": false,
    "validation": "required,minBytes=100B,maxBytes=50MB"
  },
  {
    "key": "storage.level.target.import.trigger.count",
    "value": 50000,
    "defaultValue": 50000,
    "overwritten": false,
    "protected": false,
    "validation": "required,min=1,max=10000000"
  },
  {
    "key": "storage.level.target.import.trigger.interval",
    "value": "5m0s",
    "defaultValue": "5m0s",
    "overwritten": false,
    "protected": false,
    "validation": "required,minDuration=60s,maxDuration=24h"
  },
  {
    "key": "storage.level.target.import.trigger.size",
    "value": "5MB",
    "defaultValue": "5MB",
    "overwritten": false,
    "protected": false,
    "validation": "required,minBytes=100B,maxBytes=500MB"
  }
]
`), strings.TrimSpace(json.MustEncodeString(kvs, true)))
}

func TestConfig_BindKVs_Ok(t *testing.T) {
	t.Parallel()

	patch := config.Patch{}
	require.NoError(t, configpatch.BindKVs(&patch, []configpatch.PatchKV{
		{
			KeyPath: "storage.level.local.volume.allocation.static",
			Value:   "456MB",
		},
	}))

	assert.Equal(t, config.Patch{
		Storage: &storage.ConfigPatch{
			Level: &level.ConfigPatch{
				Local: &local.ConfigPatch{
					Volume: &volume.ConfigPatch{
						Allocation: &diskalloc.ConfigPatch{
							Static: test.Ptr(456 * datasize.MB),
						},
					},
				},
			},
		},
	}, patch)
}

func TestConfig_BindKVs_InvalidType(t *testing.T) {
	t.Parallel()

	err := configpatch.BindKVs(&config.Patch{}, []configpatch.PatchKV{
		{
			KeyPath: "storage.level.local.compression.gzip.level",
			Value:   "foo",
		},
	})

	if assert.Error(t, err) {
		assert.Equal(t, `invalid "storage.level.local.compression.gzip.level" value: found type "string", expected "int"`, err.Error())
	}
}

func TestConfig_BindKVs_InvalidValue(t *testing.T) {
	t.Parallel()

	err := configpatch.BindKVs(&config.Patch{}, []configpatch.PatchKV{
		{
			KeyPath: "storage.level.local.volume.allocation.static",
			Value:   "foo",
		},
	})

	if assert.Error(t, err) {
		assert.Equal(t, `invalid "storage.level.local.volume.allocation.static" value "foo": invalid syntax`, err.Error())
	}
}