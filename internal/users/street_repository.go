package users

// StreetRepository defines the interface for street data access.
type StreetRepository interface {
	GetAllStreets() []Street
}
