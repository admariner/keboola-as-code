package manifest

import (
	"github.com/iancoleman/strcase"
	"github.com/spf13/cast"
	"keboola-as-code/src/model"
	"regexp"
	"strings"
)

// LocalNaming of the files
type LocalNaming struct {
	Branch    PathTemplate `json:"branch" validate:"required"`
	Config    PathTemplate `json:"config" validate:"required"`
	ConfigRow PathTemplate `json:"configRow" validate:"required"`
}

func DefaultNaming() *LocalNaming {
	return &LocalNaming{
		Branch:    "{branch_id}-{branch_name}",
		Config:    "{component_type}/{component_id}/{config_id}-{config_name}",
		ConfigRow: "rows/{config_row_id}-{config_row_name}",
	}
}

func (n *LocalNaming) BranchPath(branch *model.Branch) string {
	return n.replace(string(n.Branch), map[string]interface{}{
		"branch_id":   branch.Id,
		"branch_name": n.normalizeName(branch.Name),
	})
}

func (n *LocalNaming) ConfigPath(component *model.Component, config *model.Config) string {
	return n.replace(string(n.Config), map[string]interface{}{
		"component_type": component.Type,
		"component_id":   component.Id,
		"config_id":      config.Id,
		"config_name":    n.normalizeName(config.Name),
	})
}

func (n *LocalNaming) ConfigRowPath(row *model.ConfigRow) string {
	return n.replace(string(n.ConfigRow), map[string]interface{}{
		"config_row_id":   row.Id,
		"config_row_name": n.normalizeName(row.Name),
	})
}

func (n *LocalNaming) normalizeName(name string) string {
	str := regexp.
		MustCompile(`[^a-zA-Z0-9]+`).
		ReplaceAllString(strcase.ToDelimited(name, '-'), "-")
	return strings.Trim(str, "-")
}

func (n *LocalNaming) replace(path string, placeholders map[string]interface{}) string {
	for key, value := range placeholders {
		path = strings.ReplaceAll(path, "{"+key+"}", cast.ToString(value))
	}
	return path
}