package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutage_Create(t *testing.T) {
	period, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	addr, _ := NewOutageAddress(1, "Стрийська", []string{"10", "12"}, "Львів")
	desc := NewOutageDescription("test")
	o := &Outage{ID: 42, Period: period, Address: addr, Description: desc}
	assert.Equal(t, 42, o.ID)
	assert.Equal(t, period, o.Period)
	assert.Equal(t, addr, o.Address)
	assert.Equal(t, desc, o.Description)
}

func TestOutage_AffectsUserAddress_Matching(t *testing.T) {
	period, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	addr, _ := NewOutageAddress(1, "Стрийська", []string{"10", "12"}, "")
	o := &Outage{ID: 1, Period: period, Address: addr, Description: NewOutageDescription("")}
	userAddr, err := NewUserAddress(1, "Стрийська", "10")
	require.NoError(t, err)
	assert.True(t, o.AffectsUserAddress(userAddr))
}

func TestOutage_AffectsUserAddress_DifferentStreet(t *testing.T) {
	period, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	addr, _ := NewOutageAddress(1, "Стрийська", []string{"10"}, "")
	o := &Outage{ID: 1, Period: period, Address: addr, Description: NewOutageDescription("")}
	userAddr, _ := NewUserAddress(2, "Наукова", "10")
	assert.False(t, o.AffectsUserAddress(userAddr))
}

func TestOutage_AffectsUserAddress_BuildingNotInList(t *testing.T) {
	period, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	addr, _ := NewOutageAddress(1, "Стрийська", []string{"10", "12"}, "")
	o := &Outage{ID: 1, Period: period, Address: addr, Description: NewOutageDescription("")}
	userAddr, _ := NewUserAddress(1, "Стрийська", "14")
	assert.False(t, o.AffectsUserAddress(userAddr))
}

func TestOutage_AffectsUserAddress_WithCity(t *testing.T) {
	period, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	addr, _ := NewOutageAddress(1, "Стрийська", []string{"10"}, "Львів")
	o := &Outage{ID: 1, Period: period, Address: addr, Description: NewOutageDescription("")}
	userAddr, _ := NewUserAddress(1, "Стрийська", "10")
	assert.True(t, o.AffectsUserAddress(userAddr))
}
