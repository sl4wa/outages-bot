package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutageAddress_CreateValid(t *testing.T) {
	addr, err := NewOutageAddress(1, "Стрийська", []string{"10", "12", "14"}, "")
	require.NoError(t, err)
	assert.Equal(t, 1, addr.StreetID)
	assert.Equal(t, "Стрийська", addr.StreetName)
	assert.Equal(t, []string{"10", "12", "14"}, addr.Buildings)
}

func TestOutageAddress_CreateWithCity(t *testing.T) {
	addr, err := NewOutageAddress(1, "Стрийська", []string{"10"}, "Львів")
	require.NoError(t, err)
	assert.Equal(t, "Львів", addr.City)
}

func TestOutageAddress_CoversUserAddress_Matching(t *testing.T) {
	outageAddr, _ := NewOutageAddress(1, "Стрийська", []string{"10", "12", "14"}, "")
	userAddr, _ := NewUserAddress(1, "Стрийська", "12")
	assert.True(t, outageAddr.CoversUserAddress(userAddr))
}

func TestOutageAddress_CoversUserAddress_DifferentStreet(t *testing.T) {
	outageAddr, _ := NewOutageAddress(1, "Стрийська", []string{"10", "12"}, "")
	userAddr, _ := NewUserAddress(2, "Наукова", "12")
	assert.False(t, outageAddr.CoversUserAddress(userAddr))
}

func TestOutageAddress_CoversUserAddress_BuildingNotInList(t *testing.T) {
	outageAddr, _ := NewOutageAddress(1, "Стрийська", []string{"10", "12"}, "")
	userAddr, _ := NewUserAddress(1, "Стрийська", "14")
	assert.False(t, outageAddr.CoversUserAddress(userAddr))
}

func TestOutageAddress_CoversUserAddress_DataProvider(t *testing.T) {
	tests := []struct {
		name      string
		streetID  int
		buildings []string
		userSID   int
		userBldg  string
		expected  bool
	}{
		{"same street matching building", 1, []string{"10", "12"}, 1, "10", true},
		{"same street non-matching building", 1, []string{"10", "12"}, 1, "14", false},
		{"different street", 1, []string{"10"}, 2, "10", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outageAddr, _ := NewOutageAddress(tt.streetID, "Street", tt.buildings, "")
			userAddr, _ := NewUserAddress(tt.userSID, "Street", tt.userBldg)
			assert.Equal(t, tt.expected, outageAddr.CoversUserAddress(userAddr))
		})
	}
}

func TestOutageAddress_NonPositiveStreetID(t *testing.T) {
	_, err := NewOutageAddress(0, "Street", []string{"10"}, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "street ID must be positive")

	_, err = NewOutageAddress(-1, "Street", []string{"10"}, "")
	assert.Error(t, err)
}

func TestOutageAddress_EmptyStreetName(t *testing.T) {
	_, err := NewOutageAddress(1, "", []string{"10"}, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "street name cannot be empty")

	_, err = NewOutageAddress(1, "   ", []string{"10"}, "")
	assert.Error(t, err)
}

func TestOutageAddress_EmptyBuildings(t *testing.T) {
	_, err := NewOutageAddress(1, "Street", []string{}, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "buildings must be non-empty strings")
}

func TestOutageAddress_EmptyBuildingString(t *testing.T) {
	_, err := NewOutageAddress(1, "Street", []string{"10", ""}, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "buildings must be non-empty strings")

	_, err = NewOutageAddress(1, "Street", []string{"  "}, "")
	assert.Error(t, err)
}
