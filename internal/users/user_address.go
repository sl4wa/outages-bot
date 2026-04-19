package users

import (
	"regexp"
	"strings"
)

var buildingPattern = regexp.MustCompile(`^[0-9]+(-[A-ZА-ЯІЇЄҐ])?$`)

// UserAddress represents a user's street address.
type UserAddress struct {
	StreetID   int
	StreetName string
	Building   string
}

// NewUserAddress creates a new UserAddress with validation.
func NewUserAddress(streetID int, streetName, building string) (UserAddress, error) {
	if streetID <= 0 {
		return UserAddress{}, ErrInvalidStreetID
	}

	if strings.TrimSpace(streetName) == "" {
		return UserAddress{}, ErrEmptyStreetName
	}

	if strings.TrimSpace(building) == "" {
		return UserAddress{}, ErrEmptyBuilding
	}

	if !buildingPattern.MatchString(building) {
		return UserAddress{}, ErrInvalidBuildingFormat
	}

	return UserAddress{
		StreetID:   streetID,
		StreetName: streetName,
		Building:   building,
	}, nil
}
