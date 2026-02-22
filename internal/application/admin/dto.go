package admin

import (
	"fmt"
	"time"
)

const (
	dateFormat     = "02.01.2006"
	timeFormat     = "15:04"
	dateTimeFormat = dateFormat + " " + timeFormat
)

// PeriodFormatter formats outage periods for CLI display.
func PeriodFormatter(start, end time.Time) string {
	if start.Format("2006-01-02") == end.Format("2006-01-02") {
		return fmt.Sprintf("%s - %s", start.Format(dateTimeFormat), end.Format(timeFormat))
	}
	return fmt.Sprintf("%s - %s", start.Format(dateTimeFormat), end.Format(dateTimeFormat))
}
