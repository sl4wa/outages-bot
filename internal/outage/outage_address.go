package outage

import "strings"

// Address represents the location of an outage.
type Address struct {
	StreetID   int
	StreetName string
	Buildings  []string
	City       string
}

// NewAddress creates a new Address with validation.
func NewAddress(streetID int, streetName string, buildings []string, city string) (Address, error) {
	if streetID <= 0 {
		return Address{}, ErrInvalidStreetID
	}

	if strings.TrimSpace(streetName) == "" {
		return Address{}, ErrEmptyStreetName
	}

	if len(buildings) == 0 {
		return Address{}, ErrEmptyBuildings
	}

	for _, b := range buildings {
		if strings.TrimSpace(b) == "" {
			return Address{}, ErrEmptyBuildings
		}
	}

	return Address{
		StreetID:   streetID,
		StreetName: streetName,
		Buildings:  buildings,
		City:       city,
	}, nil
}
