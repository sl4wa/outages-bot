package notifier

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"outages-bot/internal/application/service"
	"outages-bot/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sentNotification struct {
	UserID  int64
	Content NotificationContent
}

type mockSender struct {
	sent []sentNotification
	err  error
}

func (m *mockSender) Send(userID int64, content NotificationContent) error {
	m.sent = append(m.sent, sentNotification{UserID: userID, Content: content})
	if m.err != nil {
		return m.err
	}
	return nil
}

type mockProvider struct {
	outages []service.OutageDTO
	err     error
}

func (m *mockProvider) FetchOutages(_ context.Context) ([]service.OutageDTO, error) {
	return m.outages, m.err
}

func makeTestOutage(streetID int, buildings []string) service.OutageDTO {
	return service.OutageDTO{
		ID:         1,
		Start:      time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		End:        time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		City:       "Львів",
		StreetID:   streetID,
		StreetName: "Стрийська",
		Buildings:  buildings,
		Comment:    "test",
	}
}

func newNotifyUsers(provider *mockProvider, sender *mockSender, repo *mockUserRepo, logger *log.Logger) *NotifyUsers {
	fetchService := service.NewFetchOutages(provider)
	return NewNotifyUsers(fetchService, sender, repo, logger)
}

func TestNotifyUsers_MatchingSendsAndSaves(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10", "12"})}}
	svc := newNotifyUsers(provider, sender, repo, log.New(io.Discard, "", 0))

	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Len(t, sender.sent, 1)
	assert.Equal(t, int64(100), sender.sent[0].UserID)
	assert.Len(t, repo.saved, 1)
}

func TestNotifyUsers_BlockedUserRemoved(t *testing.T) {
	sender := &mockSender{
		err: &NotificationSendError{UserID: 100, Code: 403, Message: "Forbidden"},
	}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10"})}}
	svc := newNotifyUsers(provider, sender, repo, log.New(io.Discard, "", 0))

	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Contains(t, repo.removed, int64(100))
}

func TestNotifyUsers_NonMatchingDoesNothing(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(2, "Наукова", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10"})}}
	svc := newNotifyUsers(provider, sender, repo, log.New(io.Discard, "", 0))

	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Empty(t, sender.sent)
	assert.Empty(t, repo.saved)
}

func TestNotifyUsers_DedupSecondRunSkips(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10"})}}
	svc := newNotifyUsers(provider, sender, repo, log.New(io.Discard, "", 0))

	// First run
	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Len(t, sender.sent, 1)

	// Second run — same outage, same user (now with outageInfo set)
	sender.sent = nil
	err = svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Empty(t, sender.sent)
}

func TestNotifyUsers_NonBlockingSendError_UserNotRemoved(t *testing.T) {
	sender := &mockSender{
		err: &NotificationSendError{UserID: 100, Code: 500, Message: "Internal Server Error"},
	}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10"})}}
	svc := newNotifyUsers(provider, sender, repo, log.New(io.Discard, "", 0))

	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Empty(t, repo.removed)
	assert.Empty(t, repo.saved) // Not saved because send failed
}

func TestNotifyUsers_SaveError_LogsAndContinues(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	repo.saveErr = errors.New("disk full")
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10"})}}
	svc := newNotifyUsers(provider, sender, repo, logger)

	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Len(t, sender.sent, 1)
	assert.Contains(t, buf.String(), "failed to save user 100")
}

func TestNotifyUsers_RemoveError_LogsAndContinues(t *testing.T) {
	sender := &mockSender{
		err: &NotificationSendError{UserID: 100, Code: 403, Message: "Forbidden"},
	}
	repo := newMockUserRepo()
	repo.removeErr = errors.New("disk full")
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10"})}}
	svc := newNotifyUsers(provider, sender, repo, logger)

	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "failed to remove blocked user 100")
}

func TestNotifyUsers_FetchError(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()

	provider := &mockProvider{err: errors.New("api down")}
	svc := newNotifyUsers(provider, sender, repo, log.New(io.Discard, "", 0))

	err := svc.Handle(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch outages")
	assert.Contains(t, err.Error(), "api down")
}
