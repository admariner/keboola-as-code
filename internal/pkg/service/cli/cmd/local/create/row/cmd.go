package row

import (
	"github.com/spf13/cobra"

	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/helpmsg"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/configmap"
	createRow "github.com/keboola/keboola-as-code/pkg/lib/operation/project/local/create/row"
	loadState "github.com/keboola/keboola-as-code/pkg/lib/operation/state/load"
)

type Flags struct {
	Branch configmap.Value[string] `configKey:"branch" configShorthand:"b" configUsage:"branch ID or name"`
	Config configmap.Value[string] `configKey:"config" configShorthand:"c" configUsage:"config name or ID"`
	Name   configmap.Value[string] `configKey:"name" configShorthand:"n" configUsage:"name of the new config row"`
}

func DefaultFlags() Flags {
	return Flags{}
}

// nolint: dupl
func Command(p dependencies.Provider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "row",
		Short: helpmsg.Read(`local/create/row/short`),
		Long:  helpmsg.Read(`local/create/row/long`),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command must be used in project directory
			prj, d, err := p.LocalProject(cmd.Context(), false)
			if err != nil {
				return err
			}

			// Load project state
			projectState, err := prj.LoadState(loadState.LocalOperationOptions(), d)
			if err != nil {
				return err
			}

			// flags
			f := Flags{}
			if err = p.BaseScope().ConfigBinder().Bind(cmd.Flags(), args, &f); err != nil {
				return err
			}

			// Options
			options, err := AskCreateRow(projectState, d.Dialogs(), f)
			if err != nil {
				return err
			}

			// Create row
			return createRow.Run(cmd.Context(), projectState, options, d)
		},
	}

	configmap.MustGenerateFlags(cmd.Flags(), DefaultFlags())

	return cmd
}