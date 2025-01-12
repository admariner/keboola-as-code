// nolint: dupl
package local

import (
	"github.com/spf13/cobra"

	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/helpmsg"
	"github.com/keboola/keboola-as-code/pkg/lib/operation/project/local/persist"
	loadState "github.com/keboola/keboola-as-code/pkg/lib/operation/state/load"
)

func PersistCommand(p dependencies.Provider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "persist",
		Short: helpmsg.Read(`local/persist/short`),
		Long:  helpmsg.Read(`local/persist/long`),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Command must be used in project directory
			prj, d, err := p.LocalProject(false)
			if err != nil {
				return err
			}

			// Load project state
			projectState, err := prj.LoadState(loadState.PersistOptions(), d)
			if err != nil {
				return err
			}

			// Options
			options := persist.Options{
				DryRun:            d.Options().GetBool(`dry-run`),
				LogUntrackedPaths: true,
			}

			// Persist
			return persist.Run(d.CommandCtx(), projectState, options, d)
		},
	}

	// Flags
	cmd.Flags().Bool("dry-run", false, "print what needs to be done")

	return cmd
}
