package dialog

import (
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/state"
	"github.com/keboola/keboola-as-code/internal/pkg/template"
	"github.com/keboola/keboola-as-code/internal/pkg/template/context/upgrade"
	"github.com/keboola/keboola-as-code/internal/pkg/template/input"
	upgradeTemplate "github.com/keboola/keboola-as-code/pkg/lib/operation/project/local/template/upgrade"
)

type upgradeTmplDialog struct {
	*Dialogs
	groups input.StepsGroupsExt
	out    upgradeTemplate.Options
}

type upgradeTmplDeps interface {
	Logger() log.Logger
}

// AskUpgradeTemplateOptions - dialog for updating a template to a new version.
func (p *Dialogs) AskUpgradeTemplateOptions(d upgradeTmplDeps, projectState *state.State, branchKey model.BranchKey, instance model.TemplateInstance, groups template.StepsGroups) (upgradeTemplate.Options, error) {
	groupsExt := upgrade.ExportInputsValues(d.Logger().DebugWriter(), projectState, branchKey, instance.InstanceID, groups)
	dialog := &upgradeTmplDialog{Dialogs: p, groups: groupsExt}
	dialog.out.Branch = branchKey
	dialog.out.Instance = instance
	return dialog.ask()
}

func (d *upgradeTmplDialog) ask() (upgradeTemplate.Options, error) {
	// User inputs
	if v, _, err := d.askUseTemplateInputs(d.groups, false); err != nil {
		return d.out, err
	} else {
		d.out.Inputs = v
	}

	return d.out, nil
}
