package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestUser(t *testing.T) *User {
	t.Helper()
	addr, err := NewUserAddress(1, "Стрийська", "10")
	require.NoError(t, err)
	return &User{ID: 12345, Address: addr, OutageInfo: nil}
}

func newTestOutage(t *testing.T) *Outage {
	t.Helper()
	period, err := NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	require.NoError(t, err)
	addr, err := NewOutageAddress(1, "Стрийська", []string{"10", "12"}, "Львів")
	require.NoError(t, err)
	return &Outage{
		ID:          1,
		Period:      period,
		Address:     addr,
		Description: NewOutageDescription("Планове відключення"),
	}
}

func TestUser_Create(t *testing.T) {
	user := newTestUser(t)
	assert.Equal(t, int64(12345), user.ID)
	assert.Equal(t, "Стрийська", user.Address.StreetName)
	assert.Nil(t, user.OutageInfo)
}

func TestUser_WithNotifiedOutage(t *testing.T) {
	user := newTestUser(t)
	outage := newTestOutage(t)
	updated := user.WithNotifiedOutage(outage)
	assert.NotNil(t, updated.OutageInfo)
	assert.Equal(t, outage.Period, updated.OutageInfo.Period)
	assert.Equal(t, outage.Description, updated.OutageInfo.Description)
	// Original user unchanged
	assert.Nil(t, user.OutageInfo)
}

func TestUser_IsAlreadyNotifiedAbout_NoOutageInfo(t *testing.T) {
	user := newTestUser(t)
	info := NewOutageInfo(
		OutagePeriod{
			StartDate: time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		},
		NewOutageDescription("test"),
	)
	assert.False(t, user.IsAlreadyNotifiedAbout(info))
}

func TestUser_IsAlreadyNotifiedAbout_SameOutage(t *testing.T) {
	user := newTestUser(t)
	outage := newTestOutage(t)
	updated := user.WithNotifiedOutage(outage)
	info := NewOutageInfo(outage.Period, outage.Description)
	assert.True(t, updated.IsAlreadyNotifiedAbout(info))
}

func TestUser_IsAlreadyNotifiedAbout_DifferentOutage(t *testing.T) {
	user := newTestUser(t)
	outage := newTestOutage(t)
	updated := user.WithNotifiedOutage(outage)
	differentInfo := NewOutageInfo(
		OutagePeriod{
			StartDate: time.Date(2024, 2, 1, 8, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 2, 1, 16, 0, 0, 0, time.UTC),
		},
		NewOutageDescription("different"),
	)
	assert.False(t, updated.IsAlreadyNotifiedAbout(differentInfo))
}

func TestUser_WithNotifiedOutage_PreservesID(t *testing.T) {
	user := newTestUser(t)
	outage := newTestOutage(t)
	updated := user.WithNotifiedOutage(outage)
	assert.Equal(t, user.ID, updated.ID)
	assert.Equal(t, user.Address, updated.Address)
}

func TestUser_WithNotifiedOutage_PreservesAddress(t *testing.T) {
	user := newTestUser(t)
	outage := newTestOutage(t)
	updated := user.WithNotifiedOutage(outage)
	assert.Equal(t, user.Address.StreetID, updated.Address.StreetID)
	assert.Equal(t, user.Address.Building, updated.Address.Building)
}
