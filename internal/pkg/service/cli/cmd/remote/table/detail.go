package table

import (
	"time"

	"github.com/keboola/go-client/pkg/keboola"
	"github.com/spf13/cobra"

	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/helpmsg"
	"github.com/keboola/keboola-as-code/pkg/lib/operation/project/remote/table/detail"
)

func DetailCommand(p dependencies.Provider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `detail [table]`,
		Short: helpmsg.Read(`remote/table/detail/short`),
		Long:  helpmsg.Read(`remote/table/detail/long`),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (cmdErr error) {
			// Get dependencies
			d, err := p.RemoteCommandScope(dependencies.WithoutMasterToken())
			if err != nil {
				return err
			}

			// Ask options
			var tableID keboola.TableID
			if len(args) == 0 {
				tableID, _, err = askTable(d, false)
				if err != nil {
					return err
				}
			} else {
				id, err := keboola.ParseTableID(args[0])
				if err != nil {
					return err
				}
				tableID = id
			}

			// Send cmd successful/failed event
			defer d.EventSender().SendCmdEvent(d.CommandCtx(), time.Now(), &cmdErr, "remote-table-detail")

			return detail.Run(d.CommandCtx(), tableID, d)
		},
	}

	cmd.Flags().StringP("storage-api-host", "H", "", "storage API host, eg. \"connection.keboola.com\"")

	return cmd
}
