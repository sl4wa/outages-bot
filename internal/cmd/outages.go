package cli

import (
	"context"
	"fmt"
	"io"
	"outages-bot/internal/application"
	"outages-bot/internal/application/admin"
	"strings"
	"text/tabwriter"
)

// RunOutagesCommand fetches and prints outages in a table.
func RunOutagesCommand(ctx context.Context, provider application.OutageProvider, w io.Writer) error {
	outages, err := admin.ListOutages(ctx, provider)
	if err != nil {
		return fmt.Errorf("failed to fetch outages: %w", err)
	}

	if len(outages) == 0 {
		fmt.Fprintln(w, "No outages found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "StreetID\tStreet\tBuildings\tPeriod\tComment")

	for _, o := range outages {
		buildings := strings.Join(o.Buildings, ", ")
		period := admin.PeriodFormatter(o.Start, o.End)
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\n", o.StreetID, o.StreetName, buildings, period, o.Comment)
	}

	return tw.Flush()
}
