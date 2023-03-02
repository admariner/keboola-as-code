package table

import (
	"strings"
	"time"

	"github.com/keboola/go-client/pkg/keboola"
	"github.com/spf13/cobra"
	"github.com/umisama/go-regexpcache"

	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/helpmsg"
	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/options"
	common "github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
	"github.com/keboola/keboola-as-code/pkg/lib/operation/project/remote/table/preview"
)

func PreviewCommand(p dependencies.Provider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `preview [table]`,
		Short: helpmsg.Read(`remote/table/preview/short`),
		Long:  helpmsg.Read(`remote/table/preview/long`),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (cmdErr error) {
			// Ask for host and token if needed
			baseDeps := p.BaseDependencies()
			if err := baseDeps.Dialogs().AskHostAndToken(baseDeps); err != nil {
				return err
			}

			// Get dependencies
			d, err := p.DependenciesForRemoteCommand(common.WithoutMasterToken())
			if err != nil {
				return err
			}

			// Send cmd successful/failed event
			defer d.EventSender().SendCmdEvent(d.CommandCtx(), time.Now(), &cmdErr, "remote-table-preview")

			allTables, err := d.KeboolaProjectAPI().ListTablesRequest(keboola.WithColumns()).Send(d.CommandCtx())
			if err != nil {
				return err
			}

			table, err := d.Dialogs().AskTable(d.Options(), *allTables)
			if err != nil {
				return err
			}

			options, err := parsePreviewOptions(d.Options(), table.ID)
			if err != nil {
				return err
			}

			return preview.Run(d.CommandCtx(), options, d)
		},
	}

	cmd.Flags().StringP("storage-api-host", "H", "", "storage API host, eg. \"connection.keboola.com\"")
	cmd.Flags().String("changed-since", "", "only export rows imported after this date")
	cmd.Flags().String("changed-until", "", "only export rows imported before this date")
	cmd.Flags().StringSlice("columns", []string{}, "comma-separated list of columns to export")
	cmd.Flags().Uint("limit", 100, "limit the number of exported rows")
	cmd.Flags().String("where", "", "filter columns by value")
	cmd.Flags().String("order", "", "order by one or more columns")
	cmd.Flags().String("format", "pretty", "order by one or more columns")
	cmd.Flags().String("out", "", "export table a file")

	return cmd
}

// nolint: gochecknoglobals
var previewFormats = map[string]bool{
	"pretty": true,
	"csv":    true,
	"json":   true,
}

func parsePreviewOptions(options *options.Options, tableID keboola.TableID) (preview.PreviewOptions, error) {
	o := preview.PreviewOptions{TableID: tableID}

	o.ChangedSince = options.GetString("changed-since")
	o.ChangedUntil = options.GetString("changed-until")
	o.Columns = options.GetStringSlice("columns")
	o.Limit = options.GetUint("limit")
	o.Out = options.GetString("out")

	e := errors.NewMultiError()

	whereString := options.GetString("where")
	if len(whereString) > 0 {
		for _, s := range strings.Split(whereString, ";") {
			whereFilter, err := parseWhereFilter(s)
			if err != nil {
				e.Append(err)
				continue
			}
			o.WhereFilters = append(o.WhereFilters, whereFilter)
		}
	}

	orderString := options.GetString("order")
	if len(orderString) > 0 {
		for _, s := range strings.Split(orderString, ",") {
			columnOrder, err := parseColumnOrder(s)
			if err != nil {
				e.Append(err)
				continue
			}
			o.Order = append(o.Order, columnOrder)
		}
	}

	format := options.GetString("format")
	if !previewFormats[format] {
		return preview.PreviewOptions{}, errors.Errorf(`invalid output format "%s"`, format)
	}
	o.Format = format

	return o, nil
}

func parseWhereFilter(s string) (preview.WhereFilter, error) {
	m := regexpcache.MustCompile(`^(\w+)(=|!=|>=|<=)([^=!<>]*)$`).FindStringSubmatch(s)
	if m == nil {
		return preview.WhereFilter{}, errors.Errorf(`invalid where filter "%s"`, s)
	}

	column := m[1]
	operator, err := keboola.ParseCompareOp(m[2])
	if err != nil {
		return preview.WhereFilter{}, err
	}
	values := strings.Split(m[3], ",")

	return preview.WhereFilter{
		Column:   column,
		Operator: operator,
		Values:   values,
	}, nil
}

func parseColumnOrder(s string) (preview.ColumnOrder, error) {
	m := regexpcache.MustCompile(`(\w+)(?:=(asc|desc))?`).FindStringSubmatch(s)
	if m == nil {
		return preview.ColumnOrder{}, errors.Errorf(`invalid column order "%s"`, s)
	}

	column := m[1]
	orderString := m[2]
	if len(orderString) == 0 {
		orderString = "asc"
	}
	order, err := keboola.ParseColumnOrder(orderString)
	if err != nil {
		return preview.ColumnOrder{}, err
	}

	return preview.ColumnOrder{Column: column, Order: order}, nil
}
