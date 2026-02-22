package notification

import (
	"bytes"
	"errors"
	"io"
	"log"
	"outages-bot/internal/application"
	"outages-bot/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSender struct {
	sent []application.NotificationSenderDTO
	err  error
}

func (m *mockSender) Send(dto application.NotificationSenderDTO) error {
	m.sent = append(m.sent, dto)
	if m.err != nil {
		return m.err
	}
	return nil
}

type mockUserRepo struct {
	users     map[int64]*domain.User
	saved     []*domain.User
	removed   []int64
	saveErr   error
	removeErr error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[int64]*domain.User)}
}

func (m *mockUserRepo) FindAll() ([]*domain.User, error) {
	var result []*domain.User
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, nil
}

func (m *mockUserRepo) Find(chatID int64) (*domain.User, error) {
	return m.users[chatID], nil
}

func (m *mockUserRepo) Save(user *domain.User) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.saved = append(m.saved, user)
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Remove(chatID int64) (bool, error) {
	m.removed = append(m.removed, chatID)
	if m.removeErr != nil {
		return false, m.removeErr
	}
	delete(m.users, chatID)
	return true, nil
}

func makeTestOutage(streetID int, buildings []string) *domain.Outage {
	period, _ := domain.NewOutagePeriod(
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	)
	addr, _ := domain.NewOutageAddress(streetID, "Стрийська", buildings, "Львів")
	return &domain.Outage{
		ID:          1,
		Period:      period,
		Address:     addr,
		Description: domain.NewOutageDescription("test"),
	}
}

func TestNotification_MatchingSendsAndSaves(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	svc := NewNotificationService(sender, repo, log.New(io.Discard, "", 0))
	outages := []*domain.Outage{makeTestOutage(1, []string{"10", "12"})}

	count, err := svc.Handle(outages)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Len(t, sender.sent, 1)
	assert.Equal(t, int64(100), sender.sent[0].UserID)
	assert.Len(t, repo.saved, 1)
}

func TestNotification_BlockedUserRemoved(t *testing.T) {
	sender := &mockSender{
		err: &application.NotificationSendError{UserID: 100, Code: 403, Message: "Forbidden"},
	}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	svc := NewNotificationService(sender, repo, log.New(io.Discard, "", 0))
	outages := []*domain.Outage{makeTestOutage(1, []string{"10"})}

	svc.Handle(outages)
	assert.Contains(t, repo.removed, int64(100))
}

func TestNotification_NonMatchingDoesNothing(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(2, "Наукова", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	svc := NewNotificationService(sender, repo, log.New(io.Discard, "", 0))
	outages := []*domain.Outage{makeTestOutage(1, []string{"10"})}

	svc.Handle(outages)
	assert.Empty(t, sender.sent)
	assert.Empty(t, repo.saved)
}

func TestNotification_DedupSecondRunSkips(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	svc := NewNotificationService(sender, repo, log.New(io.Discard, "", 0))
	outages := []*domain.Outage{makeTestOutage(1, []string{"10"})}

	// First run
	svc.Handle(outages)
	assert.Len(t, sender.sent, 1)

	// Second run — same outage, same user (now with outageInfo set)
	sender.sent = nil
	svc.Handle(outages)
	assert.Empty(t, sender.sent)
}

func TestNotification_NonBlockingSendError_UserNotRemoved(t *testing.T) {
	sender := &mockSender{
		err: &application.NotificationSendError{UserID: 100, Code: 500, Message: "Internal Server Error"},
	}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	svc := NewNotificationService(sender, repo, log.New(io.Discard, "", 0))
	outages := []*domain.Outage{makeTestOutage(1, []string{"10"})}

	svc.Handle(outages)
	assert.Empty(t, repo.removed)
	assert.Empty(t, repo.saved) // Not saved because send failed
}

func TestNotification_SaveError_LogsAndContinues(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	repo.saveErr = errors.New("disk full")
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	svc := NewNotificationService(sender, repo, logger)
	outages := []*domain.Outage{makeTestOutage(1, []string{"10"})}

	count, err := svc.Handle(outages)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Len(t, sender.sent, 1)
	assert.Contains(t, buf.String(), "failed to save user 100")
}

func TestNotification_RemoveError_LogsAndContinues(t *testing.T) {
	sender := &mockSender{
		err: &application.NotificationSendError{UserID: 100, Code: 403, Message: "Forbidden"},
	}
	repo := newMockUserRepo()
	repo.removeErr = errors.New("disk full")
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	svc := NewNotificationService(sender, repo, logger)
	outages := []*domain.Outage{makeTestOutage(1, []string{"10"})}

	count, err := svc.Handle(outages)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Contains(t, buf.String(), "failed to remove blocked user 100")
}
