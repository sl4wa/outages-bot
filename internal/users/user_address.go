package users

import (
	"regexp"
	"strings"
)

var buildingPattern = regexp.MustCompile(`^[0-9]+(-[A-ZА-ЯІЇЄҐ])?$`)

// Address represents a user's street address.
type Address struct {
	StreetID   int
	StreetName string
	Building   string
}

// NewAddress creates a new Address with validation.
func NewAddress(streetID int, streetName, building string) (Address, error) {
	if streetID <= 0 {
		return Address{}, ErrInvalidStreetID
	}

	if strings.TrimSpace(streetName) == "" {
		return Address{}, ErrEmptyStreetName
	}

	if strings.TrimSpace(building) == "" {
		return Address{}, ErrEmptyBuilding
	}

	if !buildingPattern.MatchString(building) {
		return Address{}, ErrInvalidBuildingFormat
	}

	return Address{
		StreetID:   streetID,
		StreetName: streetName,
		Building:   building,
	}, nil
}
