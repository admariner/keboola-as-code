package registration

import (
	"context"
	"sync"

	etcd "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/sink/tablesink/storage/volume"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/op"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/servicectx"
)

type dependencies interface {
	Logger() log.Logger
	Process() *servicectx.Process
	EtcdClient() *etcd.Client
}

type putOpFactory func(metadata volume.Metadata, id etcd.LeaseID) op.WithResult[volume.Metadata]

// RegisterVolumes in etcd with lease, so on node failure, records are automatically removed after TTL seconds.
// On session failure, volumes are registered again by the callback.
// List of the active volumes can be read by the repository.VolumeRepository.
func RegisterVolumes[V volume.Volume](d dependencies, cfg Config, volumes *volume.Collection[V], putOpFactory putOpFactory) error {
	logger := d.Logger()
	client := d.EtcdClient()

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	d.Process().OnShutdown(func(ctx context.Context) {
		logger.InfoCtx(ctx, "received shutdown request")
		cancel()
		wg.Wait()
		logger.InfoCtx(ctx, "shutdown done")
	})

	// Register volumes
	errCh := etcdop.ResistantSession(ctx, wg, logger, client, cfg.TTLSeconds, func(session *concurrency.Session) error {
		txn := op.Txn(client)
		for _, vol := range volumes.All() {
			txn.Merge(putOpFactory(vol.Metadata(), session.Lease()))
		}
		return txn.Do(ctx).Err()
	})

	return <-errCh
}
