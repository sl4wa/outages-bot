package cli

import (
	"bytes"
	"context"
	"errors"
	"log"
	"outages-bot/internal/application/notifier"
	"outages-bot/internal/application/service"
	"outages-bot/internal/domain"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockOutageProvider struct {
	outages []service.OutageDTO
	err     error
}

func (m *mockOutageProvider) FetchOutages(_ context.Context) ([]service.OutageDTO, error) {
	return m.outages, m.err
}

type mockNotifSender struct {
	sent []notifier.NotificationContent
	err  error
}

func (m *mockNotifSender) Send(userID int64, content notifier.NotificationContent) error {
	m.sent = append(m.sent, content)
	return m.err
}

type mockUserRepo struct {
	users []*domain.User
}

func (m *mockUserRepo) FindAll() []*domain.User { return m.users }
func (m *mockUserRepo) Find(_ int64) (*domain.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Save(_ *domain.User) error    { return nil }
func (m *mockUserRepo) Remove(_ int64) (bool, error) { return false, nil }

type mockOutageRepo struct {
	outages []*domain.Outage
	saveErr error
	loadErr error
}

func (m *mockOutageRepo) Load() ([]*domain.Outage, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	return m.outages, nil
}

func (m *mockOutageRepo) Save(outages []*domain.Outage) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.outages = outages
	return nil
}

func makeTestDTO() service.OutageDTO {
	return service.OutageDTO{
		ID:         1,
		StreetID:   1,
		StreetName: "Стрийська",
		Buildings:  []string{"10"},
		Start:      time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		End:        time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
	}
}

func TestRunNotifierCommand_NormalRun_LogsProcessedSummary(t *testing.T) {
	provider := &mockOutageProvider{outages: []service.OutageDTO{makeTestDTO()}}
	sender := &mockNotifSender{}
	userRepo := &mockUserRepo{}

	fetchService := service.NewFetchOutages(provider)
	notifyUsers := notifier.NewNotifyUsers(fetchService, sender, userRepo, &mockOutageRepo{}, nil)

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	err := RunNotifierCommand(context.Background(), notifyUsers, logger)
	require.NoError(t, err)
	assert.Empty(t, buf.String())
}

func TestRunNotifierCommand_SnapshotHit_NotifierLogsSkip_CLISilent(t *testing.T) {
	provider := &mockOutageProvider{outages: []service.OutageDTO{makeTestDTO()}}
	sender := &mockNotifSender{}
	userRepo := &mockUserRepo{}
	snap := &mockOutageRepo{}

	// Use a single shared logger — mirrors production where both layers share log.Default().
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	fetchService := service.NewFetchOutages(provider)
	notifyUsers := notifier.NewNotifyUsers(fetchService, sender, userRepo, snap, logger)

	// First run — saves snapshot; clear log output afterwards.
	require.NoError(t, RunNotifierCommand(context.Background(), notifyUsers, logger))
	buf.Reset()

	// Second run — snapshot unchanged: notifier logs the skip once, CLI adds nothing.
	err := RunNotifierCommand(context.Background(), notifyUsers, logger)
	require.NoError(t, err)
	const skipMsg = "Outage data unchanged; checker/notifier logic skipped."
	assert.Equal(t, 1, strings.Count(buf.String(), skipMsg), "skip message must appear exactly once")
}

func TestRunNotifierCommand_FetchError_ReturnsError(t *testing.T) {
	provider := &mockOutageProvider{err: errors.New("api down")}
	sender := &mockNotifSender{}
	userRepo := &mockUserRepo{}

	fetchService := service.NewFetchOutages(provider)
	notifyUsers := notifier.NewNotifyUsers(fetchService, sender, userRepo, &mockOutageRepo{}, nil)

	logger := log.New(&bytes.Buffer{}, "", 0)

	err := RunNotifierCommand(context.Background(), notifyUsers, logger)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch outages")
	assert.Contains(t, err.Error(), "api down")
}
