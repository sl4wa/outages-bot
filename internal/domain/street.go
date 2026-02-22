package domain

import "strings"

// Street represents a street entity.
type Street struct {
	ID   int
	Name string
}

// NameContains checks if the street name contains the query (case-insensitive).
func (s Street) NameContains(query string) bool {
	return strings.Contains(strings.ToLower(s.Name), query)
}

// NameEquals checks if the street name equals the query (case-insensitive).
func (s Street) NameEquals(query string) bool {
	return strings.ToLower(s.Name) == query
}

// StreetRepository defines the interface for street data access.
type StreetRepository interface {
	GetAllStreets() ([]Street, error)
}
