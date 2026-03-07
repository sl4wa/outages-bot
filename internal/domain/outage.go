package domain

// Outage represents a power outage event.
type Outage struct {
	ID          int
	Period      OutagePeriod
	Address     OutageAddress
	Description OutageDescription
}
