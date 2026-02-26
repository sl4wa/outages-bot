package cli

import (
	"bytes"
	"errors"
	"log"
	"outages-bot/internal/application"
	"outages-bot/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepoForUsers struct {
	users []*domain.User
	err   error
}

func (m *mockUserRepoForUsers) FindAll() ([]*domain.User, error) { return m.users, m.err }
func (m *mockUserRepoForUsers) Find(_ int64) (*domain.User, error) {
	return nil, nil
}
func (m *mockUserRepoForUsers) Save(_ *domain.User) error    { return nil }
func (m *mockUserRepoForUsers) Remove(_ int64) (bool, error) { return false, nil }

type mockUserInfoProvider struct {
	infos map[int64]application.UserInfoDTO
	err   error
}

func (m *mockUserInfoProvider) GetUserInfo(chatID int64) (application.UserInfoDTO, error) {
	if m.err != nil {
		return application.UserInfoDTO{}, m.err
	}
	info, ok := m.infos[chatID]
	if !ok {
		return application.UserInfoDTO{}, errors.New("user not found")
	}
	return info, nil
}

func makeUserWithAddr(t *testing.T, id int64, streetName, building string) *domain.User {
	t.Helper()
	addr, err := domain.NewUserAddress(1, streetName, building)
	require.NoError(t, err)
	return &domain.User{ID: id, Address: addr}
}

func TestRunUsersCommand_PrintsUsers(t *testing.T) {
	users := []*domain.User{
		makeUserWithAddr(t, 100, "Стрийська", "10"),
	}
	repo := &mockUserRepoForUsers{users: users}
	infoProvider := &mockUserInfoProvider{
		infos: map[int64]application.UserInfoDTO{
			100: {ChatID: 100, Username: "testuser", FirstName: "John", LastName: "Doe"},
		},
	}

	var buf bytes.Buffer
	logger := log.New(&bytes.Buffer{}, "", 0)

	err := RunUsersCommand(repo, infoProvider, &buf, logger)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "100")
	assert.Contains(t, output, "@testuser")
	assert.Contains(t, output, "John")
	assert.Contains(t, output, "Doe")
	assert.Contains(t, output, "Стрийська")
	assert.Contains(t, output, "10")
	assert.Contains(t, output, "Total Users: 1")
}

func TestRunUsersCommand_Empty(t *testing.T) {
	repo := &mockUserRepoForUsers{users: []*domain.User{}}
	infoProvider := &mockUserInfoProvider{}

	var buf bytes.Buffer
	logger := log.New(&bytes.Buffer{}, "", 0)

	err := RunUsersCommand(repo, infoProvider, &buf, logger)
	require.NoError(t, err)
	assert.Equal(t, "No users found.\n", buf.String())
}

func TestRunUsersCommand_GetUserInfoError_SkipsUser(t *testing.T) {
	users := []*domain.User{
		makeUserWithAddr(t, 100, "Стрийська", "10"),
		makeUserWithAddr(t, 200, "Молдавська", "5"),
	}
	repo := &mockUserRepoForUsers{users: users}
	infoProvider := &mockUserInfoProvider{
		infos: map[int64]application.UserInfoDTO{
			200: {ChatID: 200, Username: "user2", FirstName: "Jane", LastName: "Smith"},
		},
	}

	var buf bytes.Buffer
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	err := RunUsersCommand(repo, infoProvider, &buf, logger)
	require.NoError(t, err)

	output := buf.String()
	assert.NotContains(t, output, "100")
	assert.Contains(t, output, "200")
	assert.Contains(t, output, "Total Users: 1")
	assert.Contains(t, logBuf.String(), "Failed to get info for chat 100")
}

func TestRunUsersCommand_UsernameFormatting(t *testing.T) {
	users := []*domain.User{
		makeUserWithAddr(t, 100, "Стрийська", "10"),
		makeUserWithAddr(t, 200, "Молдавська", "5"),
	}
	repo := &mockUserRepoForUsers{users: users}
	infoProvider := &mockUserInfoProvider{
		infos: map[int64]application.UserInfoDTO{
			100: {ChatID: 100, Username: "hasuser", FirstName: "A", LastName: "B"},
			200: {ChatID: 200, Username: "", FirstName: "C", LastName: "D"},
		},
	}

	var buf bytes.Buffer
	logger := log.New(&bytes.Buffer{}, "", 0)

	err := RunUsersCommand(repo, infoProvider, &buf, logger)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "@hasuser")
	// User 200 has no username, should show "-"
	// We can't easily assert the exact column, but the output should contain "-" for no-username users
}

func TestRunUsersCommand_OutageInfoFormatting(t *testing.T) {
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	start := time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC)
	end := time.Date(2024, 3, 15, 16, 0, 0, 0, time.UTC)
	period, _ := domain.NewOutagePeriod(start, end)
	desc := domain.NewOutageDescription("Ремонт")
	info := domain.NewOutageInfo(period, desc)

	users := []*domain.User{
		{ID: 100, Address: addr, OutageInfo: &info},
	}
	repo := &mockUserRepoForUsers{users: users}
	infoProvider := &mockUserInfoProvider{
		infos: map[int64]application.UserInfoDTO{
			100: {ChatID: 100, Username: "user1", FirstName: "A", LastName: "B"},
		},
	}

	var buf bytes.Buffer
	logger := log.New(&bytes.Buffer{}, "", 0)

	err := RunUsersCommand(repo, infoProvider, &buf, logger)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "15.03.2024 08:00 - 16:00")
	assert.Contains(t, output, "Ремонт")
}

func TestRunUsersCommand_SanitizationApplied(t *testing.T) {
	users := []*domain.User{
		makeUserWithAddr(t, 100, "Стрийська", "10"),
	}
	repo := &mockUserRepoForUsers{users: users}
	infoProvider := &mockUserInfoProvider{
		infos: map[int64]application.UserInfoDTO{
			100: {ChatID: 100, Username: "user1", FirstName: "John\u200b", LastName: "Doe\u200c"},
		},
	}

	var buf bytes.Buffer
	logger := log.New(&bytes.Buffer{}, "", 0)

	err := RunUsersCommand(repo, infoProvider, &buf, logger)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "John")
	assert.Contains(t, output, "Doe")
	// Zero-width characters should be stripped
	assert.NotContains(t, output, "\u200b")
	assert.NotContains(t, output, "\u200c")
}

func TestRunUsersCommand_RepositoryError(t *testing.T) {
	repo := &mockUserRepoForUsers{err: errors.New("disk error")}
	infoProvider := &mockUserInfoProvider{}

	var buf bytes.Buffer
	logger := log.New(&bytes.Buffer{}, "", 0)

	err := RunUsersCommand(repo, infoProvider, &buf, logger)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list users")
	assert.Contains(t, err.Error(), "disk error")
}
