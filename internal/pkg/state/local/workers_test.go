package local

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

func TestWorkers(t *testing.T) {
	t.Parallel()
	w := NewWorkers(context.Background())

	counter := atomic.NewInt64(0)
	w.AddWorker(func() error {
		counter.Inc()
		return nil
	})
	w.AddWorker(func() error {
		counter.Inc()
		return nil
	})

	// Not stared
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, int64(0), counter.Load())

	// Start and wait
	assert.NoError(t, w.StartAndWait())
	assert.Equal(t, int64(2), counter.Load())

	// Cannot be reused
	assert.PanicsWithError(t, `invoked local.Workers cannot be reused`, func() {
		w.StartAndWait()
	})
}

func TestWorkersErrors(t *testing.T) {
	t.Parallel()
	w := NewWorkers(context.Background())

	w.AddWorker(func() error {
		return errors.New(`first`)
	})
	w.AddWorker(func() error {
		return errors.New(`second`)
	})
	w.AddWorker(func() error {
		return nil
	})
	w.AddWorker(func() error {
		return errors.New(`third`)
	})
	w.AddWorker(func() error {
		return nil
	})

	// The order of errors is the same as the workers were defined
	err := w.StartAndWait()
	assert.Error(t, err)
	assert.Equal(t, "- first\n- second\n- third", err.Error())
}

func TestLocalUnitOfWork_workersFor(t *testing.T) {
	t.Parallel()
	manager := newTestLocalManager(t, nil)
	uow := manager.NewUnitOfWork(context.Background())

	lock := &sync.Mutex{}
	var order []int

	for _, level := range []int{3, 2, 4, 1} {
		level := level
		uow.workersFor(level).AddWorker(func() error {
			lock.Lock()
			defer lock.Unlock()
			order = append(order, level)
			return nil
		})
		uow.workersFor(level).AddWorker(func() error {
			lock.Lock()
			defer lock.Unlock()
			order = append(order, level)
			return nil
		})
	}

	// Not started
	time.Sleep(10 * time.Millisecond)
	assert.Empty(t, order)

	// Invoke
	assert.NoError(t, uow.Invoke())
	assert.Equal(t, []int{1, 1, 2, 2, 3, 3, 4, 4}, order)

	// Cannot be reused
	assert.PanicsWithError(t, `invoked local.UnitOfWork cannot be reused`, func() {
		uow.Invoke()
	})
}
