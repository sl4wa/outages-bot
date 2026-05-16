package outage

// Description represents the comment/description of an outage.
type Description struct {
	Value string
}

// NewDescription creates a new Description.
func NewDescription(value string) Description {
	return Description{Value: value}
}

// Equals checks if two Description values are equal (case-sensitive).
func (d Description) Equals(other Description) bool {
	return d.Value == other.Value
}
