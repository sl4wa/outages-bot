package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
