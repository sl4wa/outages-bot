package notification

import (
	"context"
	"errors"
	"outages-bot/internal/application"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProvider struct {
	outages []application.OutageDTO
	err     error
}

func (m *mockProvider) FetchOutages(_ context.Context) ([]application.OutageDTO, error) {
	return m.outages, m.err
}

func TestFetchService_ReturnsOutages(t *testing.T) {
	provider := &mockProvider{
		outages: []application.OutageDTO{
			{
				ID:         1,
				Start:      time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
				End:        time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
				City:       "Львів",
				StreetID:   1,
				StreetName: "Стрийська",
				Buildings:  []string{"10", "12"},
				Comment:    "Планове",
			},
		},
	}
	svc := NewOutageFetchService(provider)
	outages, err := svc.Handle(context.Background())
	require.NoError(t, err)
	require.Len(t, outages, 1)
	assert.Equal(t, 1, outages[0].ID)
	assert.Equal(t, "Стрийська", outages[0].Address.StreetName)
	assert.Equal(t, "Планове", outages[0].Description.Value)
}

func TestFetchService_EmptyProvider(t *testing.T) {
	provider := &mockProvider{outages: []application.OutageDTO{}}
	svc := NewOutageFetchService(provider)
	outages, err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Empty(t, outages)
}

func TestFetchService_ProviderError(t *testing.T) {
	provider := &mockProvider{err: errors.New("network error")}
	svc := NewOutageFetchService(provider)
	_, err := svc.Handle(context.Background())
	assert.Error(t, err)
}

func TestFetchService_SkipsInvalidPeriod(t *testing.T) {
	provider := &mockProvider{
		outages: []application.OutageDTO{
			{
				ID:         1,
				Start:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				End:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // end before start
				StreetID:   1,
				StreetName: "S",
				Buildings:  []string{"1"},
			},
		},
	}
	svc := NewOutageFetchService(provider)
	outages, err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Empty(t, outages)
}

func TestFetchService_SkipsInvalidAddress(t *testing.T) {
	provider := &mockProvider{
		outages: []application.OutageDTO{
			{
				ID:         1,
				Start:      time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
				End:        time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
				StreetID:   0, // invalid
				StreetName: "S",
				Buildings:  []string{"1"},
			},
		},
	}
	svc := NewOutageFetchService(provider)
	outages, err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Empty(t, outages)
}

func TestFetchService_MultipleOutages(t *testing.T) {
	provider := &mockProvider{
		outages: []application.OutageDTO{
			{ID: 1, Start: time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC), End: time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC), StreetID: 1, StreetName: "A", Buildings: []string{"1"}},
			{ID: 2, Start: time.Date(2024, 1, 2, 8, 0, 0, 0, time.UTC), End: time.Date(2024, 1, 2, 16, 0, 0, 0, time.UTC), StreetID: 2, StreetName: "B", Buildings: []string{"2"}},
		},
	}
	svc := NewOutageFetchService(provider)
	outages, err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Len(t, outages, 2)
}

func TestFetchService_PreservesCityAndComment(t *testing.T) {
	provider := &mockProvider{
		outages: []application.OutageDTO{
			{ID: 1, Start: time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC), End: time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC), City: "Львів", StreetID: 1, StreetName: "S", Buildings: []string{"1"}, Comment: "test comment"},
		},
	}
	svc := NewOutageFetchService(provider)
	outages, err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Львів", outages[0].Address.City)
	assert.Equal(t, "test comment", outages[0].Description.Value)
}

func TestFetchService_NilProvider(t *testing.T) {
	provider := &mockProvider{outages: nil}
	svc := NewOutageFetchService(provider)
	outages, err := svc.Handle(context.Background())
	require.NoError(t, err)
	assert.Empty(t, outages)
}
