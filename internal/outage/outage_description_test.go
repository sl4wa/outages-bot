package outage

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescription_CreateWithValue(t *testing.T) {
	d := NewDescription("some comment")
	assert.Equal(t, "some comment", d.Value)
}

func TestDescription_CreateWithEmptyString(t *testing.T) {
	d := NewDescription("")
	assert.Equal(t, "", d.Value)
}

func TestDescription_CreateWithCyrillicText(t *testing.T) {
	d := NewDescription("Планове відключення")
	assert.Equal(t, "Планове відключення", d.Value)
}

func TestDescription_EqualsIdentical(t *testing.T) {
	d1 := NewDescription("test")
	d2 := NewDescription("test")
	assert.True(t, d1.Equals(d2))
}

func TestDescription_EqualsDifferent(t *testing.T) {
	d1 := NewDescription("test1")
	d2 := NewDescription("test2")
	assert.False(t, d1.Equals(d2))
}

func TestDescription_EqualsCaseSensitive(t *testing.T) {
	d1 := NewDescription("Test")
	d2 := NewDescription("test")
	assert.False(t, d1.Equals(d2))
}

func TestDescription_JSONSerialize(t *testing.T) {
	d := NewDescription("some value")
	data, err := json.Marshal(d.Value)
	assert.NoError(t, err)
	assert.Equal(t, `"some value"`, string(data))
}

func TestDescription_JSONSerializeSpecialChars(t *testing.T) {
	d := NewDescription("value with <special> & \"chars\"")
	data, err := json.Marshal(d.Value)
	assert.NoError(t, err)
	// Go's json.Marshal escapes <, >, & by default
	var decoded string
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, d.Value, decoded)
}
