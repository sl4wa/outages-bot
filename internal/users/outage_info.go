package users

import "outages-bot/internal/outage"

// OutageInfo is a composite value object containing a period and description.
type OutageInfo struct {
	Period      outage.OutagePeriod
	Description outage.OutageDescription
}

// NewOutageInfo creates a new OutageInfo.
func NewOutageInfo(period outage.OutagePeriod, description outage.OutageDescription) OutageInfo {
	return OutageInfo{Period: period, Description: description}
}

// Equals checks if two OutageInfo values are equal.
func (i OutageInfo) Equals(other OutageInfo) bool {
	return i.Period.Equals(other.Period) && i.Description.Equals(other.Description)
}
