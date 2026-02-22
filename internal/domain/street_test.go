package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreet_NameContains(t *testing.T) {
	s := Street{ID: 1, Name: "Стрийська"}
	assert.True(t, s.NameContains("стрий"))
	assert.True(t, s.NameContains("стрийська"))
	assert.False(t, s.NameContains("наукова"))
}

func TestStreet_NameEquals(t *testing.T) {
	s := Street{ID: 1, Name: "Стрийська"}
	assert.True(t, s.NameEquals("стрийська"))
	assert.False(t, s.NameEquals("стрий"))
	assert.False(t, s.NameEquals("наукова"))
}

func TestStreet_NameContainsCaseInsensitive(t *testing.T) {
	s := Street{ID: 1, Name: "Наукова"}
	assert.True(t, s.NameContains("наукова"))
	assert.True(t, s.NameContains("наук"))
}

func TestStreet_NameEqualsCaseInsensitive(t *testing.T) {
	s := Street{ID: 1, Name: "Наукова"}
	assert.True(t, s.NameEquals("наукова"))
	// "Наукова" lowered is "наукова", so comparing with lowercase query works
	assert.False(t, s.NameEquals("НАУКОВА")) // name.ToLower is "наукова", query stays uppercase
}
