package admin

import (
	"context"
	"errors"
	"outages-bot/internal/application"
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

func TestListOutages_DelegatesToProvider(t *testing.T) {
	outages := []application.OutageDTO{
		{
			ID:         1,
			StreetID:   100,
			StreetName: "Стрийська",
			Buildings:  []string{"10"},
			Start:      time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
			End:        time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		},
	}
	provider := &mockOutageProvider{outages: outages}

	result, err := ListOutages(context.Background(), provider)
	require.NoError(t, err)
	assert.Equal(t, outages, result)
}

func TestListOutages_ProviderError(t *testing.T) {
	provider := &mockOutageProvider{err: errors.New("connection refused")}

	result, err := ListOutages(context.Background(), provider)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "connection refused", err.Error())
}
