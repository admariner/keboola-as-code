package rename

import (
	"github.com/spf13/cobra"

	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/helpmsg"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/configmap"
	renameOp "github.com/keboola/keboola-as-code/pkg/lib/operation/project/local/template/rename"
	loadState "github.com/keboola/keboola-as-code/pkg/lib/operation/state/load"
)

type Flags struct {
	Branch   configmap.Value[string] `configKey:"branch" configShorthand:"b" configUsage:"branch ID or name"`
	Instance configmap.Value[string] `configKey:"instance" configShorthand:"i" configUsage:"instance ID of the template to delete"`
	NewName  configmap.Value[string] `configKey:"new-name" configShorthand:"n" configUsage:"new name of the template instance"`
}

func DefaultFlags() Flags {
	return Flags{}
}

func Command(p dependencies.Provider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `rename`,
		Short: helpmsg.Read(`local/template/rename/short`),
		Long:  helpmsg.Read(`local/template/rename/long`),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command must be used in project directory
			prj, d, err := p.LocalProject(cmd.Context(), false)
			if err != nil {
				return err
			}

			// flags
			f := Flags{}
			if err = p.BaseScope().ConfigBinder().Bind(cmd.Flags(), args, &f); err != nil {
				return err
			}

			// Load project state
			projectState, err := prj.LoadState(loadState.LocalOperationOptions(), d)
			if err != nil {
				return err
			}

			// Ask
			renameOpts, err := AskRenameInstance(projectState, d.Dialogs(), f)
			if err != nil {
				return err
			}

			// Rename template instance
			return renameOp.Run(cmd.Context(), projectState, renameOpts, d)
		},
	}

	configmap.MustGenerateFlags(cmd.Flags(), DefaultFlags())

	return cmd
}