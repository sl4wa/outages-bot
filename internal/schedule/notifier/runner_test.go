package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sl4wa/outages-bot/internal/schedule/schedule"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunnerReturnsNilAndPreservesStateWhenNoScheduleFound(t *testing.T) {
	store := &fakeStore{state: map[time.Time]string{date(2026, 2, 12): "old"}}
	result, err := testRunner(fakeProvider{items: []schedule.Snapshot{{ScheduleDate: date(2026, 2, 12), Text: "old"}}}, store, &fakeNotifier{}).
		Execute(context.Background(), []time.Time{date(2026, 2, 13)})

	require.NoError(t, err)
	assert.Nil(t, result)
	assert.Nil(t, store.saved)
}

func TestRunnerDoesNotNotifyWhenProviderHasNoSchedulesAndNoSavedState(t *testing.T) {
	store := &fakeStore{}
	notifier := &fakeNotifier{}

	result, err := testRunner(fakeProvider{}, store, notifier).
		Execute(context.Background(), []time.Time{date(2026, 2, 13)})

	require.NoError(t, err)
	assert.Nil(t, result)
	assert.Nil(t, store.saved)
	assert.False(t, notifier.called)
}

func TestRunnerExportsLatestAndNotifiesOnChange(t *testing.T) {
	updatedOlder := time.Date(2026, 2, 13, 8, 0, 0, 0, time.UTC)
	updatedNewer := time.Date(2026, 2, 13, 10, 30, 0, 0, time.UTC)
	store := &fakeStore{}
	notifier := &fakeNotifier{}

	result, err := testRunner(fakeProvider{items: []schedule.Snapshot{
		{ScheduleDate: date(2026, 2, 13), UpdatedAt: &updatedOlder, Text: "older"},
		{ScheduleDate: date(2026, 2, 13), UpdatedAt: &updatedNewer, Text: "newer"},
	}}, store, notifier).Execute(context.Background(), []time.Time{date(2026, 2, 13)})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Changed)
	assert.True(t, result.Notified)
	assert.Equal(t, "newer", result.Text)
	assert.Equal(t, map[time.Time]string{date(2026, 2, 13): "newer"}, store.saved)
	assert.Equal(t, "<b>Графік відключень на сьогодні</b>\n\nnewer", notifier.message)
}

func TestRunnerSkipsSaveAndNotificationWhenUnchanged(t *testing.T) {
	store := &fakeStore{state: map[time.Time]string{date(2026, 2, 13): "same"}}
	notifier := &fakeNotifier{}

	result, err := testRunner(fakeProvider{items: []schedule.Snapshot{{ScheduleDate: date(2026, 2, 13), Text: "same"}}}, store, notifier).
		Execute(context.Background(), []time.Time{date(2026, 2, 13)})

	require.NoError(t, err)
	assert.False(t, result.Changed)
	assert.False(t, result.Notified)
	assert.False(t, notifier.called)
	assert.Nil(t, store.saved)
}

func TestRunnerSavesBeforeNotifyingAndPropagatesNotifyError(t *testing.T) {
	store := &fakeStore{state: map[time.Time]string{date(2026, 2, 13): "old"}}
	notifier := &fakeNotifier{err: errors.New("telegram failed")}

	_, err := testRunner(fakeProvider{items: []schedule.Snapshot{{ScheduleDate: date(2026, 2, 13), Text: "new"}}}, store, notifier).
		Execute(context.Background(), []time.Time{date(2026, 2, 13)})

	require.Error(t, err)
	assert.Equal(t, map[time.Time]string{date(2026, 2, 13): "new"}, store.saved,
		"state must be saved before notify, so a notify failure does not re-broadcast on retry")
	assert.True(t, notifier.called)
}

func TestRunnerDoesNotNotifyWhenSaveFails(t *testing.T) {
	store := &fakeStore{
		state:   map[time.Time]string{date(2026, 2, 13): "old"},
		saveErr: errors.New("disk full"),
	}
	notifier := &fakeNotifier{}

	_, err := testRunner(fakeProvider{items: []schedule.Snapshot{{ScheduleDate: date(2026, 2, 13), Text: "new"}}}, store, notifier).
		Execute(context.Background(), []time.Time{date(2026, 2, 13)})

	require.Error(t, err)
	assert.False(t, notifier.called, "notify must not be called when save fails")
}

func TestRunnerExportsOnlyTomorrowWhenTodayMissing(t *testing.T) {
	store := &fakeStore{}
	today := date(2026, 2, 13)
	tomorrow := date(2026, 2, 14)

	result, err := testRunner(fakeProvider{items: []schedule.Snapshot{{ScheduleDate: tomorrow, Text: "tomorrow"}}}, store, &fakeNotifier{}).
		Execute(context.Background(), []time.Time{today, tomorrow})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Schedules, 1)
	assert.Equal(t, tomorrow, result.Schedules[0].ScheduleDate)
	assert.Equal(t, map[time.Time]string{tomorrow: "tomorrow"}, store.saved)
}

type fakeProvider struct {
	items []schedule.Snapshot
}

func (p fakeProvider) GetSchedules(context.Context) ([]schedule.Snapshot, error) {
	return p.items, nil
}

type fakeStore struct {
	state   map[time.Time]string
	saved   map[time.Time]string
	saveErr error
}

func (s *fakeStore) Load() (map[time.Time]string, error) {
	if s.saved != nil {
		return s.saved, nil
	}
	if s.state == nil {
		return map[time.Time]string{}, nil
	}
	return s.state, nil
}

func (s *fakeStore) Save(state map[time.Time]string) error {
	if s.saveErr != nil {
		return s.saveErr
	}
	s.saved = state
	return nil
}

type fakeNotifier struct {
	called  bool
	message string
	err     error
}

func (n *fakeNotifier) Notify(ctx context.Context, message string) (bool, error) {
	_ = ctx
	n.called = true
	n.message = message
	return n.err == nil, n.err
}

func testRunner(provider ScheduleProvider, store ScheduleStateStore, notifier NotificationService) Runner {
	return Runner{
		Provider: provider,
		Store:    store,
		Notifier: notifier,
		Clock: func() time.Time {
			return time.Date(2026, 2, 13, 12, 0, 0, 0, time.UTC)
		},
		Zone: time.UTC,
	}
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
