package diff

import (
	"fmt"
	"go.uber.org/zap"
)

type ActionType int

const (
	ActionSaveLocal ActionType = iota
	ActionSaveRemote
	ActionDeleteLocal
	ActionDeleteRemote
)

type Action struct {
	Result
	Type ActionType
}

type Actions struct {
	Actions []*Action
}

func (a *Action) String() string {
	return a.StringPrefix() + " " + a.LocalPath()
}

func (a *Action) StringPrefix() string {
	switch a.Result.State() {
	case ResultNotSet:
		return "? "
	case ResultNotEqual:
		return "CH"
	case ResultEqual:
		return "= "
	default:
		if a.Type == ActionSaveLocal || a.Type == ActionSaveRemote {
			return "+ "
		} else {
			return "- "
		}
	}
}

func (a *Actions) Add(r Result, t ActionType) {
	a.Actions = append(a.Actions, &Action{r, t})
}

func (a *Actions) Log(logger *zap.SugaredLogger) *Actions {
	for range a.Actions {

	}
	return a
}

func (r *Results) PullActions() *Actions {
	actions := &Actions{}
	for _, item := range r.Results {
		switch item.State() {
		case ResultEqual:
			// nop
		case ResultNotEqual:
			actions.Add(item, ActionSaveLocal)
		case ResultOnlyInLocal:
			actions.Add(item, ActionDeleteLocal)
		case ResultOnlyInRemote:
			actions.Add(item, ActionSaveLocal)
		case ResultNotSet:
			panic(fmt.Errorf("diff was not generated"))
		}
	}
	return actions
}
