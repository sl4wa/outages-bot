package outage

import "time"

// Period represents the time period of an outage.
type Period struct {
	StartDate time.Time
	EndDate   time.Time
}

// NewPeriod creates a new Period, returning an error if start is after end.
func NewPeriod(startDate, endDate time.Time) (Period, error) {
	if startDate.After(endDate) {
		return Period{}, ErrInvalidDateRange
	}
	return Period{StartDate: startDate, EndDate: endDate}, nil
}

// Equals checks if two Period values are equal by comparing Unix timestamps.
func (p Period) Equals(other Period) bool {
	return p.StartDate.Unix() == other.StartDate.Unix() &&
		p.EndDate.Unix() == other.EndDate.Unix()
}
