package outage

import "errors"

var (
	ErrInvalidStreetID  = errors.New("ідентифікатор вулиці має бути додатним")
	ErrEmptyStreetName  = errors.New("назва вулиці не може бути порожньою")
	ErrEmptyBuildings   = errors.New("список будинків не може бути порожнім")
	ErrInvalidDateRange = errors.New("дата початку має бути раніше або дорівнювати даті завершення")
)
