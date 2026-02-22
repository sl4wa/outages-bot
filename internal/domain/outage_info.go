package domain

// OutageInfo is a composite value object containing a period and description.
type OutageInfo struct {
	Period      OutagePeriod
	Description OutageDescription
}

// NewOutageInfo creates a new OutageInfo.
func NewOutageInfo(period OutagePeriod, description OutageDescription) OutageInfo {
	return OutageInfo{Period: period, Description: description}
}

// Equals checks if two OutageInfo values are equal.
func (i OutageInfo) Equals(other OutageInfo) bool {
	return i.Period.Equals(other.Period) && i.Description.Equals(other.Description)
}
