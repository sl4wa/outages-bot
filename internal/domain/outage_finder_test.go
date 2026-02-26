package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeOutage(t *testing.T, id int, streetID int, buildings []string, desc string) *Outage {
	t.Helper()
	period, err := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	require.NoError(t, err)
	addr, err := NewOutageAddress(streetID, "Street", buildings, "")
	require.NoError(t, err)
	return &Outage{
		ID:          id,
		Period:      period,
		Address:     addr,
		Description: NewOutageDescription(desc),
	}
}

func TestOutageFinder_FindsMatchingOutage(t *testing.T) {
	addr, _ := NewUserAddress(1, "Street", "10")
	user := &User{ID: 1, Address: addr}
	outages := []*Outage{makeOutage(t, 1, 1, []string{"10", "12"}, "test")}

	result := FindOutageForNotification(user, outages)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
}

func TestOutageFinder_NoMatchReturnsNil(t *testing.T) {
	addr, _ := NewUserAddress(1, "Street", "14")
	user := &User{ID: 1, Address: addr}
	outages := []*Outage{makeOutage(t, 1, 1, []string{"10", "12"}, "test")}

	result := FindOutageForNotification(user, outages)
	assert.Nil(t, result)
}

func TestOutageFinder_AlreadyNotifiedReturnsNil(t *testing.T) {
	addr, _ := NewUserAddress(1, "Street", "10")
	outage := makeOutage(t, 1, 1, []string{"10"}, "test")
	info := NewOutageInfo(outage.Period, outage.Description)
	user := &User{ID: 1, Address: addr, OutageInfo: &info}

	result := FindOutageForNotification(user, []*Outage{outage})
	assert.Nil(t, result)
}

func TestOutageFinder_MultipleMatching_ReturnsFirst(t *testing.T) {
	addr, _ := NewUserAddress(1, "Street", "10")
	user := &User{ID: 1, Address: addr}
	o1 := makeOutage(t, 1, 1, []string{"10"}, "first")
	o2 := makeOutage(t, 2, 1, []string{"10"}, "second")

	result := FindOutageForNotification(user, []*Outage{o1, o2})
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
}

func TestOutageFinder_SameOutagesReversed_ReturnsDifferentFirst(t *testing.T) {
	addr, _ := NewUserAddress(1, "Street", "10")
	user := &User{ID: 1, Address: addr}
	o1 := makeOutage(t, 1, 1, []string{"10"}, "first")
	o2 := makeOutage(t, 2, 1, []string{"10"}, "second")

	result := FindOutageForNotification(user, []*Outage{o2, o1})
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
}

func TestOutageFinder_AlreadyNotifiedFirstMatch_SkipsSecondMatch(t *testing.T) {
	addr, _ := NewUserAddress(1, "Street", "10")
	o1 := makeOutage(t, 1, 1, []string{"10"}, "first")
	o2 := makeOutage(t, 2, 1, []string{"10"}, "second")

	// User was notified about outage 1
	info := NewOutageInfo(o1.Period, o1.Description)
	user := &User{ID: 1, Address: addr, OutageInfo: &info}

	// Even though o2 also matches, finding o1 first (already notified) returns nil
	result := FindOutageForNotification(user, []*Outage{o1, o2})
	assert.Nil(t, result)
}

func TestOutageFinder_FirstOutageGone_SecondMatchFires(t *testing.T) {
	addr, _ := NewUserAddress(1, "Street", "10")
	o1 := makeOutage(t, 1, 1, []string{"10"}, "first")
	o2 := makeOutage(t, 2, 1, []string{"10"}, "second")

	// User was notified about outage 1
	info := NewOutageInfo(o1.Period, o1.Description)
	user := &User{ID: 1, Address: addr, OutageInfo: &info}

	// Outage 1 disappeared from the list — only o2 remains.
	// Since o2 has a different description, it's not "already notified" → fires.
	result := FindOutageForNotification(user, []*Outage{o2})
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
}

func TestOutageFinder_EmptyOutageList(t *testing.T) {
	addr, _ := NewUserAddress(1, "Street", "10")
	user := &User{ID: 1, Address: addr}

	result := FindOutageForNotification(user, []*Outage{})
	assert.Nil(t, result)
}

func TestOutageFinder_DifferentStreet_ThenAlreadyNotified(t *testing.T) {
	addr, _ := NewUserAddress(1, "Street", "10")
	// Non-matching outage (different street), followed by matching but already notified
	nonMatching := makeOutage(t, 1, 2, []string{"10"}, "other street")
	matching := makeOutage(t, 2, 1, []string{"10"}, "notified")

	info := NewOutageInfo(matching.Period, matching.Description)
	user := &User{ID: 1, Address: addr, OutageInfo: &info}

	result := FindOutageForNotification(user, []*Outage{nonMatching, matching})
	assert.Nil(t, result)
}
