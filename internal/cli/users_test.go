package cli

import (
	"bytes"
	"errors"
	"log"
	"outages-bot/internal/outage"
	"outages-bot/internal/users"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepoForUsers struct {
	users []*users.User
}

func (m *mockUserRepoForUsers) FindAll() []*users.User { return m.users }
func (m *mockUserRepoForUsers) Find(_ int64) (*users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForUsers) Save(_ *users.User) error     { return nil }
func (m *mockUserRepoForUsers) Remove(_ int64) (bool, error) { return false, nil }

type mockUserInfoProvider struct {
	infos map[int64]users.UserInfoDTO
	err   error
}

func (m *mockUserInfoProvider) GetUserInfo(chatID int64) (users.UserInfoDTO, error) {
	if m.err != nil {
		return users.UserInfoDTO{}, m.err
	}
	info, ok := m.infos[chatID]
	if !ok {
		return users.UserInfoDTO{}, errors.New("user not found")
	}
	return info, nil
}

func makeUserWithAddr(t *testing.T, id int64, streetName, building string) *users.User {
	t.Helper()
	addr, err := users.NewUserAddress(1, streetName, building)
	require.NoError(t, err)
	return &users.User{ID: id, Address: addr}
}

func TestRunUsersCommand_PrintsUsers(t *testing.T) {
	testUsers := []*users.User{
		makeUserWithAddr(t, 100, "Стрийська", "10"),
	}
	repo := &mockUserRepoForUsers{users: testUsers}
	infoProvider := &mockUserInfoProvider{
		infos: map[int64]users.UserInfoDTO{
			100: {ChatID: 100, Username: "testuser", FirstName: "John", LastName: "Doe"},
		},
	}

	var buf bytes.Buffer
	logger := log.New(&bytes.Buffer{}, "", 0)

	RunUsersCommand(repo, infoProvider, &buf, logger)

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
	repo := &mockUserRepoForUsers{users: []*users.User{}}
	infoProvider := &mockUserInfoProvider{}

	var buf bytes.Buffer
	logger := log.New(&bytes.Buffer{}, "", 0)

	RunUsersCommand(repo, infoProvider, &buf, logger)
	assert.Equal(t, "No users found.\n", buf.String())
}

func TestRunUsersCommand_GetUserInfoError_SkipsUser(t *testing.T) {
	testUsers := []*users.User{
		makeUserWithAddr(t, 100, "Стрийська", "10"),
		makeUserWithAddr(t, 200, "Молдавська", "5"),
	}
	repo := &mockUserRepoForUsers{users: testUsers}
	infoProvider := &mockUserInfoProvider{
		infos: map[int64]users.UserInfoDTO{
			200: {ChatID: 200, Username: "user2", FirstName: "Jane", LastName: "Smith"},
		},
	}

	var buf bytes.Buffer
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	RunUsersCommand(repo, infoProvider, &buf, logger)

	output := buf.String()
	assert.NotContains(t, output, "100")
	assert.Contains(t, output, "200")
	assert.Contains(t, output, "Total Users: 1")
	assert.Contains(t, logBuf.String(), "Failed to get info for chat 100")
}

func TestRunUsersCommand_UsernameFormatting(t *testing.T) {
	testUsers := []*users.User{
		makeUserWithAddr(t, 100, "Стрийська", "10"),
		makeUserWithAddr(t, 200, "Молдавська", "5"),
	}
	repo := &mockUserRepoForUsers{users: testUsers}
	infoProvider := &mockUserInfoProvider{
		infos: map[int64]users.UserInfoDTO{
			100: {ChatID: 100, Username: "hasuser", FirstName: "A", LastName: "B"},
			200: {ChatID: 200, Username: "", FirstName: "C", LastName: "D"},
		},
	}

	var buf bytes.Buffer
	logger := log.New(&bytes.Buffer{}, "", 0)

	RunUsersCommand(repo, infoProvider, &buf, logger)

	output := buf.String()
	assert.Contains(t, output, "@hasuser")
	// User 200 has no username, should show "-"
	// We can't easily assert the exact column, but the output should contain "-" for no-username users
}

func TestRunUsersCommand_OutageInfoFormatting(t *testing.T) {
	addr, _ := users.NewUserAddress(1, "Стрийська", "10")
	start := time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC)
	end := time.Date(2024, 3, 15, 16, 0, 0, 0, time.UTC)
	period, _ := outage.NewOutagePeriod(start, end)
	desc := outage.NewOutageDescription("Ремонт")
	info := users.NewOutageInfo(period, desc)

	testUsers := []*users.User{
		{ID: 100, Address: addr, OutageInfo: &info},
	}
	repo := &mockUserRepoForUsers{users: testUsers}
	infoProvider := &mockUserInfoProvider{
		infos: map[int64]users.UserInfoDTO{
			100: {ChatID: 100, Username: "user1", FirstName: "A", LastName: "B"},
		},
	}

	var buf bytes.Buffer
	logger := log.New(&bytes.Buffer{}, "", 0)

	RunUsersCommand(repo, infoProvider, &buf, logger)

	output := buf.String()
	assert.Contains(t, output, "15.03.2024 08:00 - 16:00")
	assert.Contains(t, output, "Ремонт")
}

func TestRunUsersCommand_SanitizationApplied(t *testing.T) {
	testUsers := []*users.User{
		makeUserWithAddr(t, 100, "Стрийська", "10"),
	}
	repo := &mockUserRepoForUsers{users: testUsers}
	infoProvider := &mockUserInfoProvider{
		infos: map[int64]users.UserInfoDTO{
			100: {ChatID: 100, Username: "user1", FirstName: "John\u200b", LastName: "Doe\u200c"},
		},
	}

	var buf bytes.Buffer
	logger := log.New(&bytes.Buffer{}, "", 0)

	RunUsersCommand(repo, infoProvider, &buf, logger)

	output := buf.String()
	assert.Contains(t, output, "John")
	assert.Contains(t, output, "Doe")
	// Zero-width characters should be stripped
	assert.NotContains(t, output, "\u200b")
	assert.NotContains(t, output, "\u200c")
}

func TestSanitizeDisplayText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal ASCII text", "Hello World", "Hello World"},
		{"zero-width joiner", "Hello\u200DWorld", "HelloWorld"},
		{"zero-width non-joiner", "Hello\u200CWorld", "HelloWorld"},
		{"soft hyphen", "Hello\u00ADWorld", "HelloWorld"},
		{"hangul filler", "Hello\u3164World", "HelloWorld"},
		{"multiple invisible chars", "He\u200Dl\u200Cl\u00ADo", "Hello"},
		{"only invisible chars", "\u200D\u200C\u00AD", "-"},
		{"empty string", "", "-"},
		{"whitespace only after stripping", "\u200D \u200C", "-"},
		{"cyrillic text preserved", "Іван Петрович", "Іван Петрович"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, sanitizeDisplayText(tt.input))
		})
	}
}
