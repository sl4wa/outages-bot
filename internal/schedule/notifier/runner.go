package notifier

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/sl4wa/outages-bot/internal/schedule/message"
	"github.com/sl4wa/outages-bot/internal/schedule/schedule"
)

type ScheduleProvider interface {
	GetSchedules(ctx context.Context) ([]schedule.Snapshot, error)
}

type ScheduleStateStore interface {
	Load() (map[time.Time]string, error)
	Save(state map[time.Time]string) error
}

type NotificationService interface {
	Notify(ctx context.Context, text string) (bool, error)
}

type Runner struct {
	Provider ScheduleProvider
	Store    ScheduleStateStore
	Notifier NotificationService
	Clock    func() time.Time
	Zone     *time.Location
}

type ExportResult struct {
	Schedules    []schedule.Snapshot
	Text         string
	Changed      bool
	Notified     bool
	ChangedDates []time.Time
}

func (r Runner) Execute(ctx context.Context, dates []time.Time) (*ExportResult, error) {
	if len(dates) == 0 {
		return nil, fmt.Errorf("dates must not be empty")
	}
	if r.Provider == nil {
		return nil, fmt.Errorf("schedule provider is required")
	}
	if r.Store == nil {
		return nil, fmt.Errorf("state store is required")
	}
	if r.Notifier == nil {
		return nil, fmt.Errorf("notification service is required")
	}
	clock := r.Clock
	if clock == nil {
		clock = time.Now
	}
	zone := r.Zone
	if zone == nil {
		zone = time.UTC
	}

	normalizedDates := normalizeDates(dates)
	schedules, err := r.Provider.GetSchedules(ctx)
	if err != nil {
		return nil, err
	}

	var selected []schedule.Snapshot
	for _, date := range normalizedDates {
		if item, ok := schedule.SelectLatestForDate(schedules, date); ok {
			selected = append(selected, item)
		}
	}
	if len(selected) == 0 {
		return nil, nil
	}

	currentState, err := r.Store.Load()
	if err != nil {
		return nil, err
	}
	newState := make(map[time.Time]string, len(normalizedDates))
	for _, date := range normalizedDates {
		if text, ok := currentState[date]; ok {
			newState[date] = text
		}
	}

	var changedDates []time.Time
	for _, item := range selected {
		date := schedule.NormalizeDate(item.ScheduleDate)
		if newState[date] != item.Text {
			changedDates = append(changedDates, date)
		}
		newState[date] = item.Text
	}

	changed := len(changedDates) > 0
	if changed {
		if err := r.Store.Save(newState); err != nil {
			return nil, err
		}
	}
	notified := false
	if changed {
		now := clock().In(zone)
		text := message.FormatMessage(selected, currentState, schedule.NormalizeDate(now), schedule.TimeOfDayFromTime(now))
		notified, err = r.Notifier.Notify(ctx, text)
		if err != nil {
			return nil, err
		}
	}

	parts := make([]string, 0, len(selected))
	for _, item := range selected {
		parts = append(parts, item.Text)
	}

	return &ExportResult{
		Schedules:    selected,
		Text:         strings.Join(parts, "\n\n"),
		Changed:      changed,
		Notified:     notified,
		ChangedDates: changedDates,
	}, nil
}

func normalizeDates(dates []time.Time) []time.Time {
	seen := make(map[time.Time]struct{}, len(dates))
	for _, date := range dates {
		seen[schedule.NormalizeDate(date)] = struct{}{}
	}
	result := make([]time.Time, 0, len(seen))
	for date := range seen {
		result = append(result, date)
	}
	slices.SortFunc(result, func(a, b time.Time) int { return a.Compare(b) })
	return result
}
