package service

import (
	"context"
	"sync"

	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

func (s *service) cleanup(ctx context.Context, wg *sync.WaitGroup) <-chan error {
	logger := s.deps.Logger().AddPrefix("[cleanup]")

	initDone := make(chan error)
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := s.deps.Clock().Ticker(s.config.TasksCleanupInterval)
		defer ticker.Stop()

		logger.Infof("ready")
		close(initDone) // no error expected

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := s.tasks.Cleanup(); err != nil && !errors.Is(err, context.Canceled) {
					logger.Error(err)
				}
			}
		}
	}()

	return initDone
}
