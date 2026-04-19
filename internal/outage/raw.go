package outage

import (
	"context"
	"time"
)

// RawProvider fetches outage rows from an external source.
type RawProvider interface {
	FetchOutages(ctx context.Context) ([]RawOutage, error)
}

// RawOutage is the raw outage row used by the CLI and normalization pipeline.
type RawOutage struct {
	ID         int
	Start      time.Time
	End        time.Time
	City       string
	StreetID   int
	StreetName string
	Buildings  []string
	Comment    string
}
