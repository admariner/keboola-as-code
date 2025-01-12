// nolint: dupl
package local

import (
	"github.com/spf13/cobra"

	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/helpmsg"
	"github.com/keboola/keboola-as-code/pkg/lib/operation/project/local/encrypt"
	loadState "github.com/keboola/keboola-as-code/pkg/lib/operation/state/load"
)

func EncryptCommand(p dependencies.Provider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: helpmsg.Read(`local/encrypt/short`),
		Long:  helpmsg.Read(`local/encrypt/long`),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command must be used in project directory
			prj, d, err := p.LocalProject(false)
			if err != nil {
				return err
			}

			// Load project state
			projectState, err := prj.LoadState(loadState.LocalOperationOptions(), d)
			if err != nil {
				return err
			}

			// Options
			options := encrypt.Options{
				DryRun:   d.Options().GetBool(`dry-run`),
				LogEmpty: true,
			}

			// Encrypt
			return encrypt.Run(d.CommandCtx(), projectState, options, d)
		},
	}
	cmd.Flags().Bool("dry-run", false, "print what needs to be done")
	return cmd
}
