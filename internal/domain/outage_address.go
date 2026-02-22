package domain

import (
	"errors"
	"strings"
)

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
		return OutageAddress{}, errors.New("street ID must be positive")
	}

	if strings.TrimSpace(streetName) == "" {
		return OutageAddress{}, errors.New("street name cannot be empty")
	}

	if len(buildings) == 0 {
		return OutageAddress{}, errors.New("buildings must be non-empty strings")
	}

	for _, b := range buildings {
		if strings.TrimSpace(b) == "" {
			return OutageAddress{}, errors.New("buildings must be non-empty strings")
		}
	}

	return OutageAddress{
		StreetID:   streetID,
		StreetName: streetName,
		Buildings:  buildings,
		City:       city,
	}, nil
}

// CoversUserAddress checks if this outage address covers the given user address.
func (a OutageAddress) CoversUserAddress(userAddr UserAddress) bool {
	if a.StreetID != userAddr.StreetID {
		return false
	}
	for _, b := range a.Buildings {
		if b == userAddr.Building {
			return true
		}
	}
	return false
}
