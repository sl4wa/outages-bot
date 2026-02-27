package domain

import (
	"errors"
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
		return UserAddress{}, errors.New("Невірний ідентифікатор вулиці")
	}

	if strings.TrimSpace(streetName) == "" {
		return UserAddress{}, errors.New("Назва вулиці не може бути порожньою")
	}

	if strings.TrimSpace(building) == "" {
		return UserAddress{}, errors.New("Невірний формат номера будинку")
	}

	if !buildingPattern.MatchString(building) {
		return UserAddress{}, errors.New("Невірний формат номера будинку. Приклад: 13 або 13-А")
	}

	return UserAddress{
		StreetID:   streetID,
		StreetName: streetName,
		Building:   building,
	}, nil
}
