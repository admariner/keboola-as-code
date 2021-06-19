package actions

import (
	"go.uber.org/zap"
	"keboola-as-code/src/diff"
)

type ActionType int

const (
	ActionSaveLocal ActionType = iota
	ActionSaveRemote
	ActionDeleteLocal
	ActionDeleteRemote
)

type Action struct {
	diff.Result
	Type ActionType
}

//SaveRemote(logger *zap.SugaredLogger, workers *errgroup.Group, pool *client.Pool, a *api.StorageApi) error
//DeleteRemote(logger *zap.SugaredLogger, workers *errgroup.Group, pool *client.Pool, a *api.StorageApi) error
//DeleteLocal(logger *zap.SugaredLogger, workers *errgroup.Group) error
//SaveLocal(logger *zap.SugaredLogger, workers *errgroup.Group) error

type Actions struct {
	Actions []*Action
}

func (a *Action) String() string {
	return a.StringPrefix() + " " + a.LocalPath()
}

func (a *Action) StringPrefix() string {
	switch a.Result.State() {
	case diff.ResultNotSet:
		return "? "
	case diff.ResultNotEqual:
		return "CH"
	case diff.ResultEqual:
		return "= "
	default:
		if a.Type == ActionSaveLocal || a.Type == ActionSaveRemote {
			return "+ "
		} else {
			return "- "
		}
	}
}

func (a *Actions) Add(r diff.Result, t ActionType) {
	a.Actions = append(a.Actions, &Action{r, t})
}

func (a *Actions) Log(logger *zap.SugaredLogger) *Actions {
	logger.Debugf("Planned actions:")
	for _, action := range a.Actions {
		logger.Debugf(action.String())
	}
	return a
}
