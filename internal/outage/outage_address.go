package outage

import "strings"

// OutageAddress represents the location of an outage.
type OutageAddress struct {
	StreetID   int
	StreetName string
	Buildings  []string
	City       string
}

// NewOutageAddress creates a new OutageAddress with validation.
func NewOutageAddress(streetID int, streetName string, buildings []string, city string) (OutageAddress, error) {
	if streetID <= 0 {
		return OutageAddress{}, ErrInvalidStreetID
	}

	if strings.TrimSpace(streetName) == "" {
		return OutageAddress{}, ErrEmptyStreetName
	}

	if len(buildings) == 0 {
		return OutageAddress{}, ErrEmptyBuildings
	}

	for _, b := range buildings {
		if strings.TrimSpace(b) == "" {
			return OutageAddress{}, ErrEmptyBuildings
		}
	}

	return OutageAddress{
		StreetID:   streetID,
		StreetName: streetName,
		Buildings:  buildings,
		City:       city,
	}, nil
}
