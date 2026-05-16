package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectLatestForDateUsesSourceIDTiebreaker(t *testing.T) {
	date := NormalizeDate(time.Date(2026, 2, 18, 0, 0, 0, 0, time.UTC))
	updatedAt := time.Date(2026, 2, 18, 7, 5, 0, 0, time.UTC)
	firstID, secondID := 1004, 1005
	first := Snapshot{ScheduleDate: date, UpdatedAt: &updatedAt, Text: "first", SourceID: &firstID}
	second := Snapshot{ScheduleDate: date, UpdatedAt: &updatedAt, Text: "second", SourceID: &secondID}

	selected, ok := SelectLatestForDate([]Snapshot{first, second}, date)

	require.True(t, ok)
	assert.Equal(t, second, selected)
}

func TestParseTextExtractsFieldsAnd24HourTime(t *testing.T) {
	date := NormalizeDate(time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC))
	text := "Графік погодинних відключень на 13.02.2026\n\nІнформація станом на 10:00 13.02.2026\n\n" +
		"Група 1.1. Електроенергія є.\nГрупа 2.1. Електроенергії немає з 22:00 до 24:00."

	parsed := ParseText(Snapshot{ScheduleDate: date, Text: text})

	assert.Equal(t, date, parsed.ScheduleDate)
	assert.Equal(t, "10:00 13.02.2026", parsed.InfoTimestamp)
	require.NotNil(t, parsed.InfoTime)
	assert.Equal(t, TimeOfDay{Minutes: 10 * 60}, *parsed.InfoTime)
	require.Len(t, parsed.Groups, 2)
	assert.Empty(t, parsed.Groups[0].Outages)
	assert.Equal(t, []TimeInterval{{From: TimeOfDay{Minutes: 22 * 60}, To: TimeOfDay{Minutes: 24 * 60}}}, parsed.Groups[1].Outages)
}

func TestParseTextMissingInfoFields(t *testing.T) {
	parsed := ParseText(Snapshot{
		ScheduleDate: NormalizeDate(time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)),
		Text:         "Група 1.1. Електроенергія є.",
	})

	assert.Empty(t, parsed.InfoTimestamp)
	assert.Nil(t, parsed.InfoTime)
	assert.Nil(t, parsed.InfoDate)
}
