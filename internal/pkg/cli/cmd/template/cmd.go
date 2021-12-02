package template

import (
	"github.com/spf13/cobra"

	"github.com/keboola/keboola-as-code/internal/pkg/cli/cmd/template/repository"
	"github.com/keboola/keboola-as-code/internal/pkg/cli/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/cli/helpmsg"
)

func Commands(d dependencies.Provider) *cobra.Command {
	cmd := &cobra.Command{
		Use:  `template`,
		Long: helpmsg.Read(`template/long`),
	}
	cmd.AddCommand(
		UseCommand(d),
		DescribeCommand(d),
		EditCommand(d),
		repository.Commands(d),
	)
	return cmd
}
