// Package cache provides local cache for files and slices statistics.
// It is used for fast resolving of the upload/import conditions.
// The cache is synced via the etcd Watch API.
package cache

import (
	"context"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	etcd "go.etcd.io/etcd/client/v3"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/key"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/model"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/schema"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/slicestate"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/prefixtree"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/servicectx"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

const (
	// prefixBuffered contains cached stats about records buffered in the etcd.
	prefixBuffered = "buffered/"
	// prefixUploading contains cached stats about records in the process of uploading from the etcd to the file storage.
	prefixUploading = "uploading/"
	// prefixUploaded contains cached stats about uploaded records.
	prefixUploaded = "uploaded/"
)

type Node struct {
	logger log.Logger
	clock  clock.Clock
	client *etcd.Client
	schema *schema.Schema
	cache  *prefixtree.AtomicTree[model.Stats]
}

type Dependencies interface {
	Logger() log.Logger
	Clock() clock.Clock
	Process() *servicectx.Process
	EtcdClient() *etcd.Client
	Schema() *schema.Schema
}

func NewNode(d Dependencies) (*Node, error) {
	// Create
	n := &Node{
		logger: d.Logger().AddPrefix("[stats][cache]"),
		clock:  d.Clock(),
		client: d.EtcdClient(),
		schema: d.Schema(),
		cache:  prefixtree.New[model.Stats](),
	}

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	d.Process().OnShutdown(func() {
		n.logger.Info("received shutdown request")
		cancel()
		wg.Wait()
		n.logger.Info("shutdown done")
	})

	// Stop on initialization error
	startTime := time.Now()
	if err := <-watchOpenedSlices(ctx, wg, n); err != nil {
		return nil, err
	}
	if err := <-watchClosedSlices(ctx, wg, n); err != nil {
		return nil, err
	}

	n.logger.Infof(`initialized | %s`, time.Since(startTime))
	return n, nil
}

func (n *Node) SliceStats(k key.SliceKey) model.StatsByType {
	return n.statsFor(k.String())
}

func (n *Node) FileStats(k key.FileKey) model.StatsByType {
	return n.statsFor(k.String())
}

func (n *Node) ExportStats(k key.ExportKey) model.StatsByType {
	return n.statsFor(k.String())
}

func (n *Node) statsFor(prefix string) (out model.StatsByType) {
	n.cache.Atomic(func(t *prefixtree.Tree[model.Stats]) {
		t.WalkPrefix(prefixBuffered+prefix, func(_ string, v model.Stats) bool {
			out.Total = out.Total.Add(v)
			out.Buffered = out.Buffered.Add(v)
			return false
		})
		t.WalkPrefix(prefixUploading+prefix, func(_ string, v model.Stats) bool {
			out.Total = out.Total.Add(v)
			out.Uploading = out.Uploading.Add(v)
			return false
		})
		t.WalkPrefix(prefixUploaded+prefix, func(_ string, v model.Stats) bool {
			out.Total = out.Total.Add(v)
			out.Uploaded = out.Uploaded.Add(v)
			return false
		})
	})
	return out
}

// watchOpenedSlices operation watches for statistics of slices in writing/closing state.
// These temporary statistics are stored separately.
// The key has format "<sliceKey>/<apiNode>".
func watchOpenedSlices(ctx context.Context, wg *sync.WaitGroup, n *Node) <-chan error {
	// The WithFilterDelete option is used, so only PUT events are watched and statistics are only inserted.
	// Delete operation is part of the watchClosedSlices, to make transitions between states atomic and prevent duplicates.
	return n.schema.ReceivedStats().
		GetAllAndWatch(ctx, n.client, etcd.WithFilterDelete()).
		SetupConsumer(n.logger).
		WithForEach(func(events []etcdop.WatchEventT[model.SliceStats], header *etcdop.Header, restart bool) {
			n.cache.Atomic(func(t *prefixtree.Tree[model.Stats]) {
				// Process all PUt events, the keys will be cleared by the watchClosedSlices.
				for _, event := range events {
					statsPerAPINode := event.Value
					switch event.Type {
					case etcdop.CreateEvent, etcdop.UpdateEvent:
						// This event may arrive later than the event in watchClosedSlices.
						// Therefore, we have to check whether the state of the slice has not changed.
						_, found1 := t.Get(prefixUploading + statsPerAPINode.SliceKey.String())
						_, found2 := t.Get(prefixUploaded + statsPerAPINode.SliceKey.String())
						if !found1 && !found2 {
							// Slice is still open, insert statistics per API node.
							t.Insert(keyForActiveStats(statsPerAPINode.SliceNodeKey), statsPerAPINode.GetStats())
						}
					default:
						panic(errors.Errorf(`unexpected event type "%v"`, event.Type))
					}
				}
			})
		}).
		StartConsumer(wg)
}

// watchActiveSlices operation watches for statistics of slices in uploading/failed/uploaded state.
// These statistics are stored inside the model.Slice structure.
func watchClosedSlices(ctx context.Context, wg *sync.WaitGroup, n *Node) <-chan error {
	return n.schema.Slices().AllClosed().
		GetAllAndWatch(ctx, n.client, etcd.WithPrevKV()).
		SetupConsumer(n.logger).
		WithForEach(func(events []etcdop.WatchEventT[model.Slice], header *etcdop.Header, restart bool) {
			n.cache.Atomic(func(t *prefixtree.Tree[model.Stats]) {
				// Reset the tree after receiving the first batch after the restart.
				if restart {
					t.Reset()
				}

				// Atomically process all events
				for _, event := range events {
					slice := event.Value
					cacheKey := keyForClosedSlice(slice)

					// Update slice stats
					switch event.Type {
					case etcdop.CreateEvent:
						// Clear temporary stats to prevent duplicates.
						// The slice was closed and the statistics were moved directly to the slice.
						t.DeletePrefix(prefixForActiveStats(slice.SliceKey))
						fallthrough
					case etcdop.UpdateEvent:
						stats := event.Value.GetStats()
						t.Insert(cacheKey, stats)
					case etcdop.DeleteEvent:
						t.Delete(cacheKey)
					default:
						panic(errors.Errorf(`unexpected event type "%v"`, event.Type))
					}
				}
			})
		}).
		StartConsumer(wg)
}

func keyForActiveStats(v key.SliceNodeKey) string {
	return prefixBuffered + v.String() + "/"
}

func prefixForActiveStats(v key.SliceKey) string {
	return prefixBuffered + v.String() + "/"
}

func keyForClosedSlice(v model.Slice) string {
	switch v.State {
	case slicestate.Uploading, slicestate.Failed:
		// uploading/<sliceKey>
		return prefixUploading + v.SliceKey.String()
	case slicestate.Uploaded:
		// uploaded/<sliceKey>
		return prefixUploaded + v.SliceKey.String()
	default:
		panic(errors.Errorf(`unexpected state "%s"`, v.State.String()))
	}
}