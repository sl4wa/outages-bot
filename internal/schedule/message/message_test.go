package message

import (
	"testing"
	"time"

	"github.com/sl4wa/outages-bot/internal/schedule/schedule"

	"github.com/stretchr/testify/assert"
)

func TestChangedGroups(t *testing.T) {
	oldText := "Група 4.1. Відключення з 09:00 до 11:00.\nГрупа 5.1. Відключення з 09:00 до 11:00.\nГрупа 6.1. Електроенергія є."
	newText := "Група 4.1. Відключення з 09:00 до 10:00.\nГрупа 5.1. Відключення з 09:00 до 10:00.\nГрупа 6.1. Електроенергія є."

	assert.Equal(t, []string{"4.1", "5.1"}, ChangedGroups(oldText, newText))
	assert.Equal(t, []string{"4.1", "5.1", "6.1"}, ChangedGroups("", newText))
}

func TestFormatBlockFutureDateMode(t *testing.T) {
	date := schedule.NormalizeDate(time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC))
	parsed := schedule.ParsedSchedule{
		ScheduleDate:  date,
		InfoTimestamp: "18:00 13.02.2026",
		Groups: []schedule.GroupEntry{
			{ID: "1.1", Outages: []schedule.TimeInterval{{From: schedule.TimeOfDay{Minutes: 9 * 60}, To: schedule.TimeOfDay{Minutes: 11 * 60}}}},
		},
	}

	assert.Equal(t,
		"<b>Графік відключень на завтра</b>\n<i>Станом на 18:00 13.02.2026</i>\n\n<b>Відключення:</b>\n<b>1.1</b>: 09:00-11:00",
		FormatBlock(parsed, nil, date.AddDate(0, 0, -1), schedule.TimeOfDay{Minutes: 18 * 60}),
	)
}

func TestFormatBlockTodayCurrentAndLater(t *testing.T) {
	date := schedule.NormalizeDate(time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC))
	parsed := schedule.ParsedSchedule{
		ScheduleDate:  date,
		InfoTimestamp: "09:30 13.02.2026",
		Groups: []schedule.GroupEntry{
			{ID: "1.2", Outages: []schedule.TimeInterval{{From: schedule.TimeOfDay{Minutes: 19 * 60}, To: schedule.TimeOfDay{Minutes: 22 * 60}}}},
			{ID: "2.1", Outages: []schedule.TimeInterval{
				{From: schedule.TimeOfDay{Minutes: 7 * 60}, To: schedule.TimeOfDay{Minutes: 10 * 60}},
				{From: schedule.TimeOfDay{Minutes: 17 * 60}, To: schedule.TimeOfDay{Minutes: 20 * 60}},
			}},
		},
	}

	expected := "<b>Графік відключень на сьогодні</b>\n" +
		"<i>Станом на 09:30 13.02.2026</i>\n" +
		"\n<b>Зараз без світла:</b>\n" +
		"<b>2.1</b>: 07:00-10:00\n" +
		"\n<b>Далі сьогодні:</b>\n" +
		"<b>1.2</b>: 19:00-22:00\n" +
		"<b>2.1</b>: 17:00-20:00"

	assert.Equal(t, expected, FormatBlock(parsed, nil, date, schedule.TimeOfDay{Minutes: 9*60 + 30}))
}

func TestFormatBlockTodayAllDayInterval(t *testing.T) {
	date := schedule.NormalizeDate(time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC))
	parsed := schedule.ParsedSchedule{
		ScheduleDate:  date,
		InfoTimestamp: "09:30 13.02.2026",
		Groups: []schedule.GroupEntry{
			{ID: "3.1", Outages: []schedule.TimeInterval{
				{From: schedule.TimeOfDay{Minutes: 0}, To: schedule.TimeOfDay{Minutes: 24 * 60}},
			}},
		},
	}

	expected := "<b>Графік відключень на сьогодні</b>\n" +
		"<i>Станом на 09:30 13.02.2026</i>\n" +
		"\n<b>Зараз без світла:</b>\n" +
		"<b>3.1</b>: 00:00-24:00"

	assert.Equal(t, expected, FormatBlock(parsed, nil, date, schedule.TimeOfDay{Minutes: 12 * 60}))
}

func TestEscapeHTML(t *testing.T) {
	assert.Equal(t, "a &amp; b &lt;tag&gt;", EscapeHTML("a & b <tag>"))
}
