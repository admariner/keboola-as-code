package actions

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"keboola-as-code/src/api"
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
