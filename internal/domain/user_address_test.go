package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAddress_ValidBuildings(t *testing.T) {
	tests := []struct {
		name     string
		building string
	}{
		{"simple number", "13"},
		{"three digit", "196"},
		{"large number", "271"},
		{"with latin suffix", "13-A"},
		{"with cyrillic А", "196-А"},
		{"with cyrillic Б", "271-Б"},
		{"with cyrillic І", "10-І"},
		{"with cyrillic Ї", "5-Ї"},
		{"with cyrillic Є", "7-Є"},
		{"with cyrillic Ґ", "3-Ґ"},
		{"single digit", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := NewUserAddress(1, "Test Street", tt.building)
			require.NoError(t, err)
			assert.Equal(t, tt.building, addr.Building)
		})
	}
}

func TestUserAddress_InvalidBuildings(t *testing.T) {
	tests := []struct {
		name     string
		building string
	}{
		{"lowercase letter", "13-a"},
		{"cyrillic lowercase", "13-а"},
		{"slash format", "13/1"},
		{"with space", "13 A"},
		{"leading zero", "013-a"}, // lowercase suffix is invalid
		{"letters only", "ABC"},
		{"special chars", "13!"},
		{"double dash", "13--A"},
		{"trailing dash", "13-"},
		{"leading dash", "-13"},
		{"two letters", "13-AB"},
		{"decimal", "13.1"},
		{"fraction", "1/2"},
		{"zero building", "0-"},
		{"negative style", "-1"},
		{"with dot suffix", "13.A"},
		{"with comma", "13,A"},
		{"unicode other", "13-Ω"},
		{"just dash", "-"},
		{"alphanumeric mix", "A13"},
		{"number dash number", "13-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUserAddress(1, "Test Street", tt.building)
			assert.Error(t, err)
		})
	}
}

func TestUserAddress_EmptyBuilding(t *testing.T) {
	_, err := NewUserAddress(1, "Test Street", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Невірний формат номера будинку")
}

func TestUserAddress_WhitespaceBuilding(t *testing.T) {
	_, err := NewUserAddress(1, "Test Street", "   ")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Невірний формат номера будинку")
}

func TestUserAddress_ZeroStreetID(t *testing.T) {
	_, err := NewUserAddress(0, "Test Street", "13")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Невірний ідентифікатор вулиці")
}

func TestUserAddress_NegativeStreetID(t *testing.T) {
	_, err := NewUserAddress(-1, "Test Street", "13")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Невірний ідентифікатор вулиці")
}

func TestUserAddress_EmptyStreetName(t *testing.T) {
	_, err := NewUserAddress(1, "", "13")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Назва вулиці не може бути порожньою")
}

func TestUserAddress_WhitespaceStreetName(t *testing.T) {
	_, err := NewUserAddress(1, "   ", "13")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Назва вулиці не може бути порожньою")
}

func TestUserAddress_ValidComplete(t *testing.T) {
	addr, err := NewUserAddress(123, "Стрийська", "10-А")
	require.NoError(t, err)
	assert.Equal(t, 123, addr.StreetID)
	assert.Equal(t, "Стрийська", addr.StreetName)
	assert.Equal(t, "10-А", addr.Building)
}
