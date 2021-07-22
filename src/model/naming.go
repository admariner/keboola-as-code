package model

import (
	"fmt"
	"keboola-as-code/src/utils"
	"path/filepath"
)

const (
	MetaFile          = "meta.json"
	ConfigFile        = "config.json"
	CodeFileName      = `code` // transformation code block name without ext
	blocksDir         = `blocks`
	blockNameTemplate = utils.PathTemplate(`{block_order}-{block_name}`)
	codeNameTemplate  = utils.PathTemplate(`{code_order}-{code_name}`)
)

// Naming of the files
type Naming struct {
	Branch    utils.PathTemplate `json:"branch" validate:"required"`
	Config    utils.PathTemplate `json:"config" validate:"required"`
	ConfigRow utils.PathTemplate `json:"configRow" validate:"required"`
}

func DefaultNaming() *Naming {
	return &Naming{
		Branch:    "{branch_id}-{branch_name}",
		Config:    "{component_type}/{component_id}/{config_id}-{config_name}",
		ConfigRow: "rows/{config_row_id}-{config_row_name}",
	}
}

func (n *Naming) BranchPath(branch *Branch) string {
	return utils.ReplacePlaceholders(string(n.Branch), map[string]interface{}{
		"branch_id":   branch.Id,
		"branch_name": utils.NormalizeName(branch.Name),
	})
}

func (n *Naming) ConfigPath(component *Component, config *Config) string {
	return utils.ReplacePlaceholders(string(n.Config), map[string]interface{}{
		"component_type": component.Type,
		"component_id":   component.Id,
		"config_id":      config.Id,
		"config_name":    utils.NormalizeName(config.Name),
	})
}

func (n *Naming) ConfigRowPath(row *ConfigRow) string {
	return utils.ReplacePlaceholders(string(n.ConfigRow), map[string]interface{}{
		"config_row_id":   row.Id,
		"config_row_name": utils.NormalizeName(row.Name),
	})
}

func (n *Naming) BlocksDir(configDir string) string {
	return filepath.Join(configDir, blocksDir)
}

func (n *Naming) BlocksTmpDir(configDir string) string {
	return filepath.Join(configDir, `.new_`+blocksDir)
}

func (n *Naming) BlockPath(index int, name string) string {
	return utils.ReplacePlaceholders(string(blockNameTemplate), map[string]interface{}{
		"block_order": fmt.Sprintf(`%03d`, index+1),
		"block_name":  utils.NormalizeName(name),
	})
}

func (n *Naming) CodePath(index int, name string) string {
	return utils.ReplacePlaceholders(string(codeNameTemplate), map[string]interface{}{
		"code_order": fmt.Sprintf(`%03d`, index+1),
		"code_name":  utils.NormalizeName(name),
	})
}

func (n *Naming) CodeFileName(componentId string) string {
	return CodeFileName + "." + n.CodeFileExt(componentId)
}

func (n *Naming) CodeFileExt(componentId string) string {
	switch componentId {
	case `keboola.snowflake-transformation`:
		return `sql`
	case `keboola.synapse-transformation`:
		return `sql`
	case `keboola.oracle-transformation`:
		return `sql`
	case `keboola.r-transformation`:
		return `r`
	case `keboola.julia-transformation`:
		return `jl`
	case `keboola.python-spark-transformation`:
		return `py`
	case `keboola.python-transformation`:
		return `py`
	case `keboola.python-transformation-v2`:
		return `py`
	case `keboola.csas-python-transformation-v2`:
		return `py`
	default:
		return `txt`
	}
}
