package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sl4wa/outages-bot/internal/schedule/loe"
	"github.com/sl4wa/outages-bot/internal/schedule/notifier"
	"github.com/sl4wa/outages-bot/internal/schedule/persistence"
	"github.com/sl4wa/outages-bot/internal/schedule/schedule"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFullSchedulePipelinePersistsFixtureState(t *testing.T) {
	payload, err := os.ReadFile("testdata/schedule_response.json")
	require.NoError(t, err)
	statePath := filepath.Join(t.TempDir(), "outages.csv")
	store := persistence.NewCSVStateStore(statePath)
	runner := notifier.Runner{
		Provider: loe.Provider{LoadPayload: func(context.Context) (string, error) { return string(payload), nil }},
		Store:    store,
		Notifier: scheduleNoopNotifier{},
		Clock:    func() time.Time { return time.Date(2026, 2, 18, 12, 0, 0, 0, time.UTC) },
		Zone:     time.UTC,
	}

	result, err := runner.Execute(context.Background(), []time.Time{time.Date(2026, 2, 18, 0, 0, 0, 0, time.UTC)})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Schedules, 1)
	assert.Equal(t, time.Date(2026, 2, 18, 0, 0, 0, 0, time.UTC), result.Schedules[0].ScheduleDate)
	require.NotNil(t, result.Schedules[0].UpdatedAt)
	assert.Equal(t, time.Date(2026, 2, 18, 9, 31, 0, 0, time.UTC), *result.Schedules[0].UpdatedAt)
	require.NotNil(t, result.Schedules[0].SourceID)
	assert.Equal(t, 1006, *result.Schedules[0].SourceID)
	loaded, err := store.Load()
	require.NoError(t, err)
	assert.Contains(t, loaded[schedule.NormalizeDate(time.Date(2026, 2, 18, 0, 0, 0, 0, time.UTC))], "Інформація станом на 09:31 18.02.2026")
}

func TestFullSchedulePipelineReturnsNilForMissingDate(t *testing.T) {
	payload, err := os.ReadFile("testdata/schedule_response.json")
	require.NoError(t, err)
	statePath := filepath.Join(t.TempDir(), "missing.csv")
	runner := notifier.Runner{
		Provider: loe.Provider{LoadPayload: func(context.Context) (string, error) { return string(payload), nil }},
		Store:    persistence.NewCSVStateStore(statePath),
		Notifier: scheduleNoopNotifier{},
	}

	result, err := runner.Execute(context.Background(), []time.Time{time.Date(2026, 3, 17, 0, 0, 0, 0, time.UTC)})

	require.NoError(t, err)
	assert.Nil(t, result)
	assert.NoFileExists(t, statePath)
}

type scheduleNoopNotifier struct{}

func (scheduleNoopNotifier) Notify(context.Context, string) (bool, error) {
	return false, nil
}
