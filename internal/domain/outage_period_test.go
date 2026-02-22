package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutagePeriod_CreateWithValidDates(t *testing.T) {
	start := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC)
	p, err := NewOutagePeriod(start, end)
	require.NoError(t, err)
	assert.Equal(t, start, p.StartDate)
	assert.Equal(t, end, p.EndDate)
}

func TestOutagePeriod_EqualDatesAllowed(t *testing.T) {
	ts := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	p, err := NewOutagePeriod(ts, ts)
	require.NoError(t, err)
	assert.Equal(t, ts, p.StartDate)
	assert.Equal(t, ts, p.EndDate)
}

func TestOutagePeriod_StartAfterEndReturnsError(t *testing.T) {
	start := time.Date(2024, 1, 2, 8, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	_, err := NewOutagePeriod(start, end)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "start date must be before or equal to end date")
}

func TestOutagePeriod_EqualsIdentical(t *testing.T) {
	start := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC)
	p1, _ := NewOutagePeriod(start, end)
	p2, _ := NewOutagePeriod(start, end)
	assert.True(t, p1.Equals(p2))
}

func TestOutagePeriod_EqualsValueNotIdentity(t *testing.T) {
	s1 := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	e1 := time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC)
	s2 := time.Date(2024, 1, 1, 8, 0, 0, 0, time.FixedZone("", 0))
	e2 := time.Date(2024, 1, 1, 16, 0, 0, 0, time.FixedZone("", 0))
	p1, _ := NewOutagePeriod(s1, e1)
	p2, _ := NewOutagePeriod(s2, e2)
	assert.True(t, p1.Equals(p2))
}

func TestOutagePeriod_DifferentStart(t *testing.T) {
	p1, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	p2, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	assert.False(t, p1.Equals(p2))
}

func TestOutagePeriod_DifferentEnd(t *testing.T) {
	p1, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	p2, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC),
	)
	assert.False(t, p1.Equals(p2))
}

func TestOutagePeriod_CompletelyDifferent(t *testing.T) {
	p1, _ := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	p2, _ := NewOutagePeriod(
		time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 1, 18, 0, 0, 0, time.UTC),
	)
	assert.False(t, p1.Equals(p2))
}
