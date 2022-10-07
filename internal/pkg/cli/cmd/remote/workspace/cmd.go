package workspace

import (
	"github.com/spf13/cobra"

	"github.com/keboola/keboola-as-code/internal/pkg/cli/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/cli/helpmsg"
)

func Commands(p dependencies.Provider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `workspace`,
		Short: helpmsg.Read(`remote/workspace/short`),
		Long:  helpmsg.Read(`remote/workspace/long`),
	}
	cmd.AddCommand(
		CreateCommand(p),
		ListCommand(p),
		DeleteCommand(p),
	)
	return cmd
}