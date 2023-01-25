package apinode

import (
	"context"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/spf13/cast"
	etcd "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

// RevisionSyncer syncs current revision of the cached state to the etcd.
// Worker nodes watch for the updates.
//
// If the state of the API node is synchronized to a new revision,
// then the API node waits until all requests that uses an older revision are completed.
//
// RevisionSyncer used by a request is locked by the RevisionSyncer.Lock and unlocked by the returned UnlockFn.
//
// Worker node wait until all API nodes acknowledge a revision greater or equal
// to the revision of switching the slice to the closing state.
// It is used to indicate that all API nodes already use a new slice.
// Then the old slice is switched from the closing to the uploading state,
// and upload can start in a Worker.
//
// Workflow:
// - The RevisionSyncer runs on all API nodes (because API nodes use cached state to improve speed of the import endpoint).
// - If the API node receives an update of the state, it calls RevisionSyncer.Notify method.
// - If the API node starts an operation (processes a request) that depends on the cached state, it calls RevisionSyncer.Lock method.
// - The Lock return an UnlockFn callback.
// - If the operation (request) is completed, then the UnlockFn is invoked.
// - The RevisionSyncer internally counts how many times a revision is in use.
// - The value is incremented by the Lock method and decremented by the unlock callback.
// - The minimum version that is currently in use is regularly synchronized to the etcd, see sync method.
type RevisionSyncer struct {
	ctx          context.Context
	wg           *sync.WaitGroup
	logger       log.Logger
	stats        StatsSyncer
	targetKey    etcdop.Key
	syncInterval time.Duration

	// session ensures that when an outage occurs,
	// the key is automatically deleted after TTL seconds,
	// so an unavailable API node does not block Worker nodes.
	session *concurrency.Session

	lock *sync.Mutex
	// currentRev is version of the cached state, it is set by the Notify method
	currentRev int64
	// syncedRev is the latest version reported by the syncer to Worker nodes, by the etcd
	syncedRev int64
	// revInUse contains the actual number of uses for each revision
	revInUse map[int64]int
}

type StatsSyncer interface {
	Sync(ctx context.Context) <-chan struct{}
}

type UnlockFn func()

func newSyncer(ctx context.Context, wg *sync.WaitGroup, clk clock.Clock, logger log.Logger, stats StatsSyncer, sess *concurrency.Session, targetKey etcdop.Key, syncInterval time.Duration) (*RevisionSyncer, error) {
	// Create
	s := &RevisionSyncer{
		ctx:          ctx,
		wg:           wg,
		logger:       logger,
		stats:        stats,
		targetKey:    targetKey,
		syncInterval: syncInterval,
		session:      sess,
		lock:         &sync.Mutex{},
		currentRev:   1,
		syncedRev:    0,
		revInUse:     make(map[int64]int),
	}

	// Initial sync
	if err := s.sync(); err != nil {
		return nil, err
	}

	// Periodical sync
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		// If the interval > 0, then sync is not triggered immediately on change, but periodically.
		// Otherwise, trigger is called immediately, see Notify and unlockRevision methods.
		var tickerC <-chan time.Time
		if s.syncInterval > 0 {
			ticker := clk.Ticker(s.syncInterval)
			defer ticker.Stop()
			tickerC = ticker.C
		}

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-tickerC:
				// Sync error is logged
				_ = s.sync()
			}
		}
	}()

	return s, nil
}

// Notify of an update of the state revision.
// Method is called by the API node.
// Value will be synced to the etcd, when no older version is in use.
func (s *RevisionSyncer) Notify(v int64) {
	s.lock.Lock()
	s.currentRev = v
	s.lock.Unlock()

	// Sync immediately, if there is no interval configured.
	// Error is logged.
	if s.syncInterval == 0 {
		_ = s.sync()
	}
}

// StateRev returns current revision of the cached state.
// It is updated by the Notify method.
func (s *RevisionSyncer) StateRev() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.currentRev
}

// SyncedRev returns the last revision reported to Worker nodes.
func (s *RevisionSyncer) SyncedRev() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.syncedRev
}

// MinRevInUse returns minimum locked revision in use by a request.
// Lock method locks the current revision
// and the returned UnlockFn callback unlocks the revision.
func (s *RevisionSyncer) MinRevInUse() int64 {
	min := s.currentRev
	for rev := range s.revInUse {
		if rev < min {
			min = rev
		}
	}
	return min
}

// Lock blocks revision sync until a dependent work is completed.
func (s *RevisionSyncer) Lock() UnlockFn {
	// Increment usage of the revision
	s.lock.Lock()
	currentRev := s.currentRev
	s.revInUse[currentRev]++ // if the map key is missing, zero value is given
	if usedCount := s.revInUse[currentRev]; usedCount == 1 {
		// Log the locked revision on the first use
		s.logger.Infof(`locked revision "%v"`, currentRev)
	}
	s.lock.Unlock()

	// Unlock callback
	return func() {
		s.unlockRevision(currentRev)
	}
}

// unlockRevision decrements version usage.
func (s *RevisionSyncer) unlockRevision(rev int64) {
	s.lock.Lock()
	s.revInUse[rev]--
	if v := s.revInUse[rev]; v == 0 {
		delete(s.revInUse, rev)
		s.logger.Infof(`unlocked revision "%v"`, rev)
	}
	s.lock.Unlock()

	// Sync immediately, if there is no interval configured.
	// Error is logged.
	if s.syncInterval == 0 {
		_ = s.sync()
	}
}

// sync the minimal revision currently in use to the etcd.
// Using this mechanism, worker nodes can safely determine when a revision is synchronized on all API nodes.
func (s *RevisionSyncer) sync() error {
	s.wg.Add(1)
	defer s.wg.Done()

	// Compare local and synced value
	s.lock.Lock()
	minRevInUse, syncedRev := s.MinRevInUse(), s.syncedRev
	s.lock.Unlock()
	if minRevInUse == syncedRev {
		// nop
		return nil
	}

	// Force statistics sync
	<-s.stats.Sync(s.ctx)

	// Update etcd key
	updateOp := s.targetKey.Put(cast.ToString(minRevInUse), etcd.WithLease(s.session.Lease()))
	if err := updateOp.Do(s.ctx, s.session.Client()); err != nil {
		if !errors.Is(err, context.Canceled) {
			s.logger.Errorf(`sync error: %s`, err)
		}
		return err
	}

	// Done, update value
	s.lock.Lock()
	s.syncedRev = minRevInUse
	s.lock.Unlock()

	s.logger.Infof(`reported revision "%v"`, minRevInUse)
	return nil
}