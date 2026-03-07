package domain

import "errors"

var (
	ErrInvalidStreetID       = errors.New("ідентифікатор вулиці має бути додатним")
	ErrEmptyStreetName       = errors.New("назва вулиці не може бути порожньою")
	ErrEmptyBuildings        = errors.New("список будинків не може бути порожнім")
	ErrInvalidDateRange      = errors.New("дата початку має бути раніше або дорівнювати даті завершення")
	ErrInvalidUserStreetID   = errors.New("невірний ідентифікатор вулиці")
	ErrEmptyUserStreetName   = errors.New("назва вулиці не може бути порожньою")
	ErrEmptyBuilding         = errors.New("номер будинку не може бути порожнім")
	ErrInvalidBuildingFormat = errors.New("невірний формат номера будинку, приклад: 13 або 13-А")
)
