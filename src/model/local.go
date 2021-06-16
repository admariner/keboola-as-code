package model

import (
	"fmt"
	"keboola-as-code/src/utils"
)

func LoadLocalState(projectDir string, metadataDir string) (*State, *utils.Error) {
	state := NewState()

	// Load manifest
	manifest, err := LoadManifest(projectDir, metadataDir)
	if err != nil {
		state.AddError(err)
		return state, state.Error()
	}

	// Add branches
	var branchById = make(map[int]*ManifestBranch)
	for _, b := range manifest.Branches {
		branch, err := b.ToModel(projectDir)
		if err == nil {
			branchById[b.Id] = b
			state.AddBranch(branch)
		} else {
			state.AddError(err)
		}
	}

	// Add configs
	for _, c := range manifest.Configs {
		if branch, ok := branchById[c.BranchId]; ok {
			config, err := c.ToModel(branch, projectDir)
			if err == nil {
				state.AddConfig(config)
			} else {
				state.AddError(err)
			}
		} else {
			state.AddError(fmt.Errorf("branch \"%d\" not found - referenced from the config \"%s:%s\" in \"%s\"", c.BranchId, c.ComponentId, c.Id, manifest.path))
		}
	}

	return state, state.Error()
}