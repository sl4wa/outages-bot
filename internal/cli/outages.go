package cli

import (
	"context"
	"fmt"
	"io"
	"outages-bot/internal/outage"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

// RunOutagesCommand fetches and prints outages in a table.
func RunOutagesCommand(ctx context.Context, provider outage.RawProvider, w io.Writer) error {
	rows, err := provider.FetchOutages(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch outages: %w", err)
	}

	if len(rows) == 0 {
		fmt.Fprintln(w, "No outages found.")
		return nil
	}

	cfg := tablewriter.NewConfigBuilder().
		WithHeaderAutoFormat(tw.Off).
		WithRowAutoWrap(tw.WrapNormal).
		ForColumn(1).WithMaxWidth(30).Build().
		ForColumn(2).WithMaxWidth(40).Build().
		ForColumn(4).WithMaxWidth(40).Build().
		Build()

	table := tablewriter.NewTable(w,
		tablewriter.WithConfig(cfg),
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{})),
	)
	table.Header([]string{"StreetID", "Street", "Buildings", "Period", "Comment"})

	for _, o := range rows {
		buildings := strings.Join(o.Buildings, ", ")
		period := PeriodFormatter(o.Start, o.End)
		table.Append([]string{fmt.Sprintf("%d", o.StreetID), o.StreetName, buildings, period, o.Comment})
	}

	return table.Render()
}
