package delete_template

import (
	"context"

	"github.com/keboola/keboola-as-code/internal/pkg/api/client/storageapi"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	deleteTemplate "github.com/keboola/keboola-as-code/internal/pkg/plan/delete-template"
	"github.com/keboola/keboola-as-code/internal/pkg/project"
	"github.com/keboola/keboola-as-code/internal/pkg/utils"
	saveManifest "github.com/keboola/keboola-as-code/pkg/lib/operation/project/local/manifest/save"
)

type Options struct {
	Branch   model.BranchKey
	DryRun   bool
	Instance string
}

type dependencies interface {
	Ctx() context.Context
	Logger() log.Logger
	StorageApi() (*storageapi.Api, error)
}

func Run(projectState *project.State, branch model.BranchKey, instance string, o Options, d dependencies) error {
	logger := d.Logger()

	// Get plan
	plan, err := deleteTemplate.NewPlan(projectState.State(), branch, instance)
	if err != nil {
		return err
	}

	// Log plan
	if !plan.Empty() {
		plan.Log(logger)
	}

	if !plan.Empty() {
		// Dry run?
		if o.DryRun {
			logger.Info("Dry run, nothing changed.")
			return nil
		}

		// Invoke
		if err := plan.Invoke(projectState.Ctx(), projectState.LocalManager()); err != nil {
			return utils.PrefixError(`cannot delete template configs`, err)
		}

		// Remove template instance from metadata
		branchState := projectState.GetOrNil(branch).(*model.BranchState)
		if err := branchState.Local.Metadata.DeleteTemplateUsage(instance); err != nil {
			return utils.PrefixError(`cannot remove template instance metadata`, err)
		}
		saveOp := projectState.LocalManager().NewUnitOfWork(projectState.Ctx())
		saveOp.SaveObject(branchState, branchState.LocalState(), model.NewChangedFields())

		// Save manifest
		if _, err := saveManifest.Run(projectState.ProjectManifest(), projectState.Fs(), d); err != nil {
			return err
		}

		logger.Info(`Delete done.`)
	}

	return nil
}
