package users

import (
	"outages-bot/internal/outage"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockListUserRepo struct {
	users []*User
}

func (m *mockListUserRepo) FindAll() []*User             { return m.users }
func (m *mockListUserRepo) Find(_ int64) (*User, error)  { return nil, nil }
func (m *mockListUserRepo) Save(_ *User) error           { return nil }
func (m *mockListUserRepo) Remove(_ int64) (bool, error) { return false, nil }

func makeAddr(t *testing.T) UserAddress {
	t.Helper()
	addr, err := NewUserAddress(1, "Стрийська", "10")
	require.NoError(t, err)
	return addr
}

func makeOutageInfo(t *testing.T, start time.Time) *OutageInfo {
	t.Helper()
	period, err := outage.NewOutagePeriod(start, start.Add(8*time.Hour))
	require.NoError(t, err)
	desc := outage.NewOutageDescription("test")
	info := NewOutageInfo(period, desc)
	return &info
}

func TestListUsers_Empty(t *testing.T) {
	repo := &mockListUserRepo{users: []*User{}}

	users := ListUsers(repo)
	assert.Empty(t, users)
}

func TestListUsers_SortedByOutageStartDescending(t *testing.T) {
	early := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	late := time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC)
	addr := makeAddr(t)

	repo := &mockListUserRepo{users: []*User{
		{ID: 1, Address: addr, OutageInfo: makeOutageInfo(t, early)},
		{ID: 2, Address: addr, OutageInfo: makeOutageInfo(t, late)},
	}}

	users := ListUsers(repo)
	require.Len(t, users, 2)
	assert.Equal(t, int64(2), users[0].ID) // late first
	assert.Equal(t, int64(1), users[1].ID) // early second
}

func TestListUsers_WithoutOutageSortedToEnd(t *testing.T) {
	ts := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	addr := makeAddr(t)

	repo := &mockListUserRepo{users: []*User{
		{ID: 1, Address: addr}, // no outage
		{ID: 2, Address: addr, OutageInfo: makeOutageInfo(t, ts)}, // has outage
	}}

	users := ListUsers(repo)
	require.Len(t, users, 2)
	assert.Equal(t, int64(2), users[0].ID) // with outage first
	assert.Equal(t, int64(1), users[1].ID) // without outage last
}

func TestListUsers_AllWithoutOutage(t *testing.T) {
	addr := makeAddr(t)

	repo := &mockListUserRepo{users: []*User{
		{ID: 1, Address: addr},
		{ID: 2, Address: addr},
	}}

	users := ListUsers(repo)
	assert.Len(t, users, 2)
}
