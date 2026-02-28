package cli

import (
	"bytes"
	"context"
	"errors"
	"log"
	"outages-bot/internal/application"
	"outages-bot/internal/application/notification"
	"outages-bot/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockOutageProvider struct {
	outages []application.OutageDTO
	err     error
}

func (m *mockOutageProvider) FetchOutages(_ context.Context) ([]application.OutageDTO, error) {
	return m.outages, m.err
}

type mockNotifSender struct {
	sent []application.NotificationSenderDTO
	err  error
}

func (m *mockNotifSender) Send(dto application.NotificationSenderDTO) error {
	m.sent = append(m.sent, dto)
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

func TestRunNotifierCommand_Success(t *testing.T) {
	provider := &mockOutageProvider{
		outages: []application.OutageDTO{
			{
				ID:         1,
				StreetID:   1,
				StreetName: "Стрийська",
				Buildings:  []string{"10"},
				Start:      time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
				End:        time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
			},
		},
	}
	sender := &mockNotifSender{}
	userRepo := &mockUserRepo{}

	fetchService := notification.NewOutageFetchService(provider)
	notifService := notification.NewService(sender, userRepo, nil)

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	err := RunNotifierCommand(context.Background(), fetchService, notifService, logger)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Successfully dispatched notifications.")
}

func TestRunNotifierCommand_FetchError(t *testing.T) {
	provider := &mockOutageProvider{err: errors.New("api down")}
	sender := &mockNotifSender{}
	userRepo := &mockUserRepo{}

	fetchService := notification.NewOutageFetchService(provider)
	notifService := notification.NewService(sender, userRepo, nil)

	logger := log.New(&bytes.Buffer{}, "", 0)

	err := RunNotifierCommand(context.Background(), fetchService, notifService, logger)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch outages")
	assert.Contains(t, err.Error(), "api down")
}

