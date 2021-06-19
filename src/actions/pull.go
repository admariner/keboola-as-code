package actions

import (
	"fmt"
	"keboola-as-code/src/diff"
)

func Pull(diffResults *diff.Results) *Actions {
	actions := &Actions{}
	for _, item := range diffResults.Results {
		switch item.State() {
		case diff.ResultEqual:
			// nop
		case diff.ResultNotEqual:
			actions.Add(item, ActionSaveLocal)
		case diff.ResultOnlyInLocal:
			actions.Add(item, ActionDeleteLocal)
		case diff.ResultOnlyInRemote:
			actions.Add(item, ActionSaveLocal)
		case diff.ResultNotSet:
			panic(fmt.Errorf("diff was not generated"))
		}
	}
	return actions
}
