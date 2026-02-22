package domain

import (
	"errors"
	"time"
)

// OutagePeriod represents the time period of an outage.
type OutagePeriod struct {
	StartDate time.Time
	EndDate   time.Time
}

// NewOutagePeriod creates a new OutagePeriod, returning an error if start is after end.
func NewOutagePeriod(startDate, endDate time.Time) (OutagePeriod, error) {
	if startDate.After(endDate) {
		return OutagePeriod{}, errors.New("start date must be before or equal to end date")
	}
	return OutagePeriod{StartDate: startDate, EndDate: endDate}, nil
}

// Equals checks if two OutagePeriod values are equal by comparing Unix timestamps.
func (p OutagePeriod) Equals(other OutagePeriod) bool {
	return p.StartDate.Unix() == other.StartDate.Unix() &&
		p.EndDate.Unix() == other.EndDate.Unix()
}
