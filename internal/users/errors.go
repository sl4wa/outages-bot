package users

import "errors"

var (
	ErrInvalidStreetID       = errors.New("невірний ідентифікатор вулиці")
	ErrEmptyStreetName       = errors.New("назва вулиці не може бути порожньою")
	ErrEmptyBuilding         = errors.New("номер будинку не може бути порожнім")
	ErrInvalidBuildingFormat = errors.New("невірний формат номера будинку, приклад: 13 або 13-А")
)
