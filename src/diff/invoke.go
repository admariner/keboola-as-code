package diff

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"keboola-as-code/src/api"
	"keboola-as-code/src/client"
	"keboola-as-code/src/utils"
)

func (a *Actions) Invoke(ctx context.Context, api *api.StorageApi, logger *zap.SugaredLogger) error {
	errors := &utils.Error{}
	workers, _ := errgroup.WithContext(ctx)
	pool := api.NewPool()
	for _, action := range a.Actions {
		switch action.Type {
		case ActionSaveLocal:
			if err := action.SaveLocal(logger, workers); err != nil {
				errors.Add(err)
			}
		case ActionSaveRemote:
			if err := action.SaveRemote(logger, workers, pool, api); err != nil {
				errors.Add(err)
			}
		case ActionDeleteLocal:
			if err := action.DeleteLocal(logger, workers); err != nil {
				errors.Add(err)
			}
		case ActionDeleteRemote:
			if err := action.DeleteRemote(logger, workers, pool, api); err != nil {
				errors.Add(err)
			}
		default:
			panic(fmt.Errorf("unexpected action type"))
		}
	}

	if err := pool.StartAndWait(); err != nil {
		errors.Add(err)
	}

	if err := workers.Wait(); err != nil {
		errors.Add(err)
	}

	if errors.Len() > 0 {
		return fmt.Errorf("pull failed: %s", errors)
	}

	return nil
}

func (b *BranchDiff) SaveLocal(logger *zap.SugaredLogger, workers *errgroup.Group) error {
	return fmt.Errorf("TODO")
}

func (c *ConfigDiff) SaveLocal(logger *zap.SugaredLogger, workers *errgroup.Group) error {
	return fmt.Errorf("TODO")
}

func (r *ConfigRowDiff) SaveLocal(logger *zap.SugaredLogger, workers *errgroup.Group) error {
	return fmt.Errorf("TODO")
}

func (b *BranchDiff) SaveRemote(logger *zap.SugaredLogger, workers *errgroup.Group, pool *client.Pool, a *api.StorageApi) error {
	return fmt.Errorf("TODO")
}

func (c *ConfigDiff) SaveRemote(logger *zap.SugaredLogger, workers *errgroup.Group, pool *client.Pool, a *api.StorageApi) error {
	return fmt.Errorf("TODO")
}

func (r *ConfigRowDiff) SaveRemote(logger *zap.SugaredLogger, workers *errgroup.Group, pool *client.Pool, a *api.StorageApi) error {
	return fmt.Errorf("TODO")
}

func (b *BranchDiff) DeleteRemote(logger *zap.SugaredLogger, workers *errgroup.Group, pool *client.Pool, a *api.StorageApi) error {
	return fmt.Errorf("TODO")
}

func (c *ConfigDiff) DeleteRemote(logger *zap.SugaredLogger, workers *errgroup.Group, pool *client.Pool, a *api.StorageApi) error {
	return fmt.Errorf("TODO")
}

func (r *ConfigRowDiff) DeleteRemote(logger *zap.SugaredLogger, workers *errgroup.Group, pool *client.Pool, a *api.StorageApi) error {
	return fmt.Errorf("TODO")
}

func (b *BranchDiff) DeleteLocal(logger *zap.SugaredLogger, workers *errgroup.Group) error {
	return fmt.Errorf("TODO")
}

func (c *ConfigDiff) DeleteLocal(logger *zap.SugaredLogger, workers *errgroup.Group) error {
	return fmt.Errorf("TODO")
}

func (r *ConfigRowDiff) DeleteLocal(logger *zap.SugaredLogger, workers *errgroup.Group) error {
	return fmt.Errorf("TODO")
}
