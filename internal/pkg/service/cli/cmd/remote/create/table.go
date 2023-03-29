package create

import (
	"time"

	"github.com/keboola/go-client/pkg/keboola"
	"github.com/spf13/cobra"

	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/helpmsg"
)

func TableCommand(p dependencies.Provider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "table [table]",
		Short: helpmsg.Read(`remote/create/table/short`),
		Long:  helpmsg.Read(`remote/create/table/long`),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (cmdErr error) {
			// Ask for host and token if needed
			baseDeps := p.BaseDependencies()
			if err := baseDeps.Dialogs().AskHostAndToken(baseDeps); err != nil {
				return err
			}

			d, err := p.DependenciesForRemoteCommand()
			if err != nil {
				return err
			}

			// Options
			var allBuckets []*keboola.Bucket
			if len(args) == 0 && !d.Options().IsSet("bucket") {
				// Get buckets list for dialog select only if needed
				allBucketsPtr, err := d.KeboolaProjectAPI().ListBucketsRequest().Send(d.CommandCtx())
				if err != nil {
					return err
				}
				allBuckets = *allBucketsPtr
			}
			_, err = d.Dialogs().AskCreateTable(args, d, allBuckets)
			if err != nil {
				return err
			}

			defer d.EventSender().SendCmdEvent(d.CommandCtx(), time.Now(), &cmdErr, "remote-create-table")

			return nil
		},
	}

	cmd.Flags().SortFlags = true
	cmd.Flags().StringP("storage-api-host", "H", "", "if command is run outside the project directory")
	cmd.Flags().String("bucket", "", "bucket ID (required if the tableId argument is empty)")
	cmd.Flags().String("name", "", "name of the table (required if the tableId argument is empty)")
	cmd.Flags().String("columns", "", "comma-separated list of column names")
	cmd.Flags().String("primary-key", "", "columns used as primary key, comma-separated")

	return cmd
}
