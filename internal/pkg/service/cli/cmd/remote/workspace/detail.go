package workspace

import (
	"time"

	"github.com/keboola/go-client/pkg/sandboxesapi"
	"github.com/spf13/cobra"

	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/helpmsg"
	"github.com/keboola/keboola-as-code/pkg/lib/operation/project/remote/workspace/detail"
)

func DetailCommand(p dependencies.Provider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `detail`,
		Short: helpmsg.Read(`remote/workspace/detail/short`),
		Long:  helpmsg.Read(`remote/workspace/detail/long`),
		RunE: func(cmd *cobra.Command, args []string) (cmdErr error) {
			start := time.Now()

			// Ask for host and token if needed
			baseDeps := p.BaseDependencies()
			if err := baseDeps.Dialogs().AskHostAndToken(baseDeps); err != nil {
				return err
			}

			d, err := p.DependenciesForRemoteCommand()
			if err != nil {
				return err
			}

			defer func() { d.EventSender().SendCmdEvent(d.CommandCtx(), start, cmdErr, "remote-detail-workspace") }()

			id, err := d.Dialogs().AskWorkspaceID(d.Options())
			if err != nil {
				return err
			}

			err = detail.Run(d.CommandCtx(), d, sandboxesapi.ConfigID(id))
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP("storage-api-host", "H", "", "storage API host, eg. \"connection.keboola.com\"")
	cmd.Flags().StringP("workspace-id", "W", "", "id of the workspace to fetch")

	return cmd
}