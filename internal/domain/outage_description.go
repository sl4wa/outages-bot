package domain

// OutageDescription represents the comment/description of an outage.
type OutageDescription struct {
	Value string
}

// NewOutageDescription creates a new OutageDescription.
func NewOutageDescription(value string) OutageDescription {
	return OutageDescription{Value: value}
}

// Equals checks if two OutageDescription values are equal (case-sensitive).
func (d OutageDescription) Equals(other OutageDescription) bool {
	return d.Value == other.Value
}
