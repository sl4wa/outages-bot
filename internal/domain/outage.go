package domain

import (
	"slices"
)

// Outage represents a power outage event.
type Outage struct {
	ID          int
	Period      OutagePeriod
	Address     OutageAddress
	Description OutageDescription
}

// OutagesEqual reports whether a and b contain the same outages in the same order.
// Comparison is positional and field-by-field.
func OutagesEqual(a, b []*Outage) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		oa, ob := a[i], b[i]
		if oa.Address.StreetID != ob.Address.StreetID ||
			oa.Address.City != ob.Address.City ||
			oa.Address.StreetName != ob.Address.StreetName ||
			!slices.Equal(oa.Address.Buildings, ob.Address.Buildings) ||
			oa.Period.StartDate.Unix() != ob.Period.StartDate.Unix() ||
			oa.Period.EndDate.Unix() != ob.Period.EndDate.Unix() ||
			oa.Description.Value != ob.Description.Value {
			return false
		}
	}
	return true
}
