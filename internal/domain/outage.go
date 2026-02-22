package domain

// Outage represents a power outage event.
type Outage struct {
	ID          int
	Period      OutagePeriod
	Address     OutageAddress
	Description OutageDescription
}

// AffectsUserAddress checks if this outage affects the given user address.
func (o *Outage) AffectsUserAddress(userAddr UserAddress) bool {
	return o.Address.CoversUserAddress(userAddr)
}
