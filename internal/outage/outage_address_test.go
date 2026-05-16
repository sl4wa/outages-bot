package outage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddress_CreateValid(t *testing.T) {
	addr, err := NewAddress(1, "Стрийська", []string{"10", "12", "14"}, "")
	require.NoError(t, err)
	assert.Equal(t, 1, addr.StreetID)
	assert.Equal(t, "Стрийська", addr.StreetName)
	assert.Equal(t, []string{"10", "12", "14"}, addr.Buildings)
}

func TestAddress_CreateWithCity(t *testing.T) {
	addr, err := NewAddress(1, "Стрийська", []string{"10"}, "Львів")
	require.NoError(t, err)
	assert.Equal(t, "Львів", addr.City)
}

func TestAddress_NonPositiveStreetID(t *testing.T) {
	_, err := NewAddress(0, "Street", []string{"10"}, "")
	assert.True(t, errors.Is(err, ErrInvalidStreetID))

	_, err = NewAddress(-1, "Street", []string{"10"}, "")
	assert.True(t, errors.Is(err, ErrInvalidStreetID))
}

func TestAddress_EmptyStreetName(t *testing.T) {
	_, err := NewAddress(1, "", []string{"10"}, "")
	assert.True(t, errors.Is(err, ErrEmptyStreetName))

	_, err = NewAddress(1, "   ", []string{"10"}, "")
	assert.True(t, errors.Is(err, ErrEmptyStreetName))
}

func TestAddress_EmptyBuildings(t *testing.T) {
	_, err := NewAddress(1, "Street", []string{}, "")
	assert.True(t, errors.Is(err, ErrEmptyBuildings))
}

func TestAddress_EmptyBuildingString(t *testing.T) {
	_, err := NewAddress(1, "Street", []string{"10", ""}, "")
	assert.True(t, errors.Is(err, ErrEmptyBuildings))

	_, err = NewAddress(1, "Street", []string{"  "}, "")
	assert.True(t, errors.Is(err, ErrEmptyBuildings))
}
