package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutageDescription_CreateWithValue(t *testing.T) {
	d := NewOutageDescription("some comment")
	assert.Equal(t, "some comment", d.Value)
}

func TestOutageDescription_CreateWithEmptyString(t *testing.T) {
	d := NewOutageDescription("")
	assert.Equal(t, "", d.Value)
}

func TestOutageDescription_CreateWithCyrillicText(t *testing.T) {
	d := NewOutageDescription("Планове відключення")
	assert.Equal(t, "Планове відключення", d.Value)
}

func TestOutageDescription_EqualsIdentical(t *testing.T) {
	d1 := NewOutageDescription("test")
	d2 := NewOutageDescription("test")
	assert.True(t, d1.Equals(d2))
}

func TestOutageDescription_EqualsDifferent(t *testing.T) {
	d1 := NewOutageDescription("test1")
	d2 := NewOutageDescription("test2")
	assert.False(t, d1.Equals(d2))
}

func TestOutageDescription_EqualsCaseSensitive(t *testing.T) {
	d1 := NewOutageDescription("Test")
	d2 := NewOutageDescription("test")
	assert.False(t, d1.Equals(d2))
}

func TestOutageDescription_JSONSerialize(t *testing.T) {
	d := NewOutageDescription("some value")
	data, err := json.Marshal(d.Value)
	assert.NoError(t, err)
	assert.Equal(t, `"some value"`, string(data))
}

func TestOutageDescription_JSONSerializeSpecialChars(t *testing.T) {
	d := NewOutageDescription("value with <special> & \"chars\"")
	data, err := json.Marshal(d.Value)
	assert.NoError(t, err)
	// Go's json.Marshal escapes <, >, & by default
	var decoded string
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, d.Value, decoded)
}
