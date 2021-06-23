package plan

import (
	"fmt"
	"keboola-as-code/src/diff"
)

func Pull(diffResults *diff.Results) *Plan {
	recipe := &Plan{Name: "pull"}
	for _, result := range diffResults.Results {
		switch result.State {
		case diff.ResultEqual:
			// nop
		case diff.ResultNotEqual:
			recipe.Add(result, ActionSaveLocal)
		case diff.ResultOnlyInLocal:
			recipe.Add(result, ActionDeleteLocal)
		case diff.ResultOnlyInRemote:
			recipe.Add(result, ActionSaveLocal)
		case diff.ResultNotSet:
			panic(fmt.Errorf("diff was not generated"))
		}
	}
	return recipe
}
