package table

import (
	"github.com/spf13/cobra"

	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/helpmsg"
)

func Commands(p dependencies.Provider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `table`,
		Short: helpmsg.Read(`remote/table/short`),
		Long:  helpmsg.Read(`remote/table/long`),
	}
	cmd.AddCommand(
		DetailCommand(p),
		ImportCommand(p),
		PreviewCommand(p),
		UnloadCommand(p),
		UploadCommand(p),
		DownloadCommand(p),
	)

	return cmd
}
