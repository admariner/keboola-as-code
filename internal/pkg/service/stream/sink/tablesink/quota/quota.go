// Package quota provides limitation of a sink buffered data in local disks.
// This prevents one client from wasting all of our disk space in the Stream API cluster.
package quota

import (
	"context"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/c2h5oh/datasize"

	commonErrors "github.com/keboola/keboola-as-code/internal/pkg/service/common/errors"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/config"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/definition/key"
	statsCache "github.com/keboola/keboola-as-code/internal/pkg/service/stream/sink/tablesink/statistics/cache"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

const (
	// MinErrorLogInterval defines the minimum interval between logged quota errors, per receiver and API node.
	// It prevents repeating errors from flooding the log.
	MinErrorLogInterval = 5 * time.Minute
)

type Checker struct {
	clock         clock.Clock
	config        config.Config
	cachedL2Stats *statsCache.L2

	// nextLogAt prevents errors from flooding the log
	nextLogAtLock *sync.RWMutex
	nextLogAt     map[key.SinkKey]time.Time
}

type dependencies interface {
	Config() config.Config
	Clock() clock.Clock
	StatisticsL2Cache() *statsCache.L2
}

func New(d dependencies) *Checker {
	return &Checker{
		clock:         d.Clock(),
		config:        d.Config(),
		cachedL2Stats: d.StatisticsL2Cache(),
		nextLogAtLock: &sync.RWMutex{},
		nextLogAt:     make(map[key.SinkKey]time.Time),
	}
}

// Check checks whether the size of records that one sink can buffer has not been exceeded.
// The method compares used disk space, on all disks, with the provided quota value.
func (c *Checker) Check(ctx context.Context, k key.SinkKey, quota datasize.ByteSize) error {
	stats, err := c.cachedL2Stats.SinkStats(ctx, k)
	if err != nil {
		return err
	}

	if diskUsage := stats.Local.CompressedSize; diskUsage > quota {
		return commonErrors.NewInsufficientStorageError(c.shouldLogError(k), errors.Errorf(
			`full storage buffer for the sink "%s", buffered "%s", quota "%s"`,
			k.String(), diskUsage.HumanReadable(), quota.HumanReadable(),
		))
	}

	return nil
}

// shouldLogError method determines if the quota error should be logged.
func (c *Checker) shouldLogError(k key.SinkKey) bool {
	now := c.clock.Now()

	c.nextLogAtLock.RLock()
	logTime := c.nextLogAt[k] // first time it returns zero time
	c.nextLogAtLock.RUnlock()

	if logTime.Before(now) {
		c.nextLogAtLock.Lock()
		c.nextLogAt[k] = now.Add(MinErrorLogInterval)
		c.nextLogAtLock.Unlock()
		return true
	}
	return false
}
