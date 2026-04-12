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
	return NewNotifyUsers(fetchService, sender, repo, &mockOutageRepo{}, logger)
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
		err: ErrRecipientUnavailable,
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
		err: errors.New("send failed"),
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
		err: ErrRecipientUnavailable,
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

func newNotifyUsersWithSnapshot(provider *mockProvider, sender *mockSender, repo *mockUserRepo, snap OutageRepository, logger *log.Logger) *NotifyUsers {
	fetchService := service.NewFetchOutages(provider)
	return NewNotifyUsers(fetchService, sender, repo, snap, logger)
}

func TestNotifyUsers_Snapshot_FirstRun_SavesAndNotifies(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	snap := &mockOutageRepo{} // no prior snapshot
	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10"})}}
	svc := newNotifyUsersWithSnapshot(provider, sender, repo, snap, log.New(io.Discard, "", 0))

	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, snap.saved, "snapshot should have been saved on first run")
	assert.Len(t, sender.sent, 1)
}

func TestNotifyUsers_Snapshot_IdenticalSecondRun_SkipsNotification(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	outages := []service.OutageDTO{makeTestOutage(1, []string{"10"})}
	provider := &mockProvider{outages: outages}
	snap := &mockOutageRepo{}

	svc := newNotifyUsersWithSnapshot(provider, sender, repo, snap, log.New(io.Discard, "", 0))

	// First run — saves snapshot and notifies
	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Len(t, sender.sent, 1)

	// Second run — same data, snapshot unchanged
	sender.sent = nil
	err = svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Empty(t, sender.sent, "no notifications on unchanged snapshot")
}

func TestNotifyUsers_Snapshot_ChangedSnapshot_SavesAndNotifies(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	firstOutages := []service.OutageDTO{makeTestOutage(1, []string{"10"})}
	provider := &mockProvider{outages: firstOutages}
	snap := &mockOutageRepo{}

	svc := newNotifyUsersWithSnapshot(provider, sender, repo, snap, log.New(io.Discard, "", 0))

	// First run
	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Len(t, sender.sent, 1)

	// Change outage data — new building added
	secondOutages := []service.OutageDTO{makeTestOutage(1, []string{"10", "12"})}
	provider.outages = secondOutages

	// Reset user so it can be notified again
	repo.users[100] = &domain.User{ID: 100, Address: addr}
	sender.sent = nil

	err = svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Len(t, sender.sent, 1, "notification sent for changed snapshot")
}

func TestNotifyUsers_Snapshot_SaveFailure_AbortsBeforeNotifications(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[100] = &domain.User{ID: 100, Address: addr}

	snap := &mockOutageRepo{saveErr: errSaveFailed}
	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10"})}}
	svc := newNotifyUsersWithSnapshot(provider, sender, repo, snap, log.New(io.Discard, "", 0))

	err := svc.Handle(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save outage data")
	assert.Empty(t, sender.sent, "no notifications sent after snapshot save failure")
}

func TestNotifyUsers_Snapshot_Hit_LogsSkipAndReturnsSkipped(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	snap := &mockOutageRepo{}
	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10"})}}

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	svc := newNotifyUsersWithSnapshot(provider, sender, repo, snap, logger)

	// First run — populates snapshot
	err := svc.Handle(context.Background())
	require.NoError(t, err)

	// Second run — snapshot hit
	buf.Reset()
	err = svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Outage data unchanged")
}

func TestNotifyUsers_Snapshot_FirstSave_LogsNoPrior(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	snap := &mockOutageRepo{} // no prior snapshot
	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10"})}}

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	svc := newNotifyUsersWithSnapshot(provider, sender, repo, snap, logger)

	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "No prior outage data")
}

func TestNotifyUsers_Snapshot_Miss_LogsChanged(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()
	snap := &mockOutageRepo{}
	provider := &mockProvider{outages: []service.OutageDTO{makeTestOutage(1, []string{"10"})}}

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	svc := newNotifyUsersWithSnapshot(provider, sender, repo, snap, logger)

	// First run — saves initial snapshot
	err := svc.Handle(context.Background())
	require.NoError(t, err)

	// Change outage data so snapshot misses
	provider.outages = []service.OutageDTO{makeTestOutage(1, []string{"10", "12"})}
	buf.Reset()

	err = svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Outage data changed")
}

func TestNotifyUsers_Snapshot_EmptyFetchFirstRun_SavesAndContinues(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()

	snap := &mockOutageRepo{} // no prior snapshot (Load returns nil)
	provider := &mockProvider{outages: nil}
	svc := newNotifyUsersWithSnapshot(provider, sender, repo, snap, log.New(io.Discard, "", 0))

	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, snap.saved, "snapshot should be saved even when API returns empty")
}

func TestNotifyUsers_Snapshot_EmptyAfterNonEmpty_TreatedAsChange(t *testing.T) {
	sender := &mockSender{}
	repo := newMockUserRepo()

	// Pre-populate snapshot with one outage
	firstOutages := []service.OutageDTO{makeTestOutage(1, []string{"10"})}
	provider := &mockProvider{outages: firstOutages}
	snap := &mockOutageRepo{}

	svc := newNotifyUsersWithSnapshot(provider, sender, repo, snap, log.New(io.Discard, "", 0))
	err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, snap.saved)

	// Now fetch returns empty
	provider.outages = nil
	sender.sent = nil

	err = svc.Handle(context.Background())
	require.NoError(t, err)
	// Snapshot should be updated to empty
	assert.Empty(t, snap.outages)
}
