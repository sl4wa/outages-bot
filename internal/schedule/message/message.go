package message

import (
	"html"
	"sort"
	"strings"
	"time"

	"github.com/sl4wa/outages-bot/internal/schedule/schedule"
)

func EscapeHTML(text string) string {
	return html.EscapeString(text)
}

func ChangedGroups(oldText, newText string) []string {
	oldGroups := groupLines(oldText)
	newGroups := groupLines(newText)
	var changed []string
	for key, newValue := range newGroups {
		if oldGroups[key] != newValue {
			changed = append(changed, key)
		}
	}
	for key := range oldGroups {
		if _, ok := newGroups[key]; !ok {
			changed = append(changed, key)
		}
	}
	sort.Strings(changed)
	return changed
}

func FormatMessage(schedules []schedule.Snapshot, previousState map[time.Time]string, today time.Time, now schedule.TimeOfDay) string {
	blocks := make([]string, 0, len(schedules))
	for _, item := range schedules {
		oldText, hadOld := previousState[schedule.NormalizeDate(item.ScheduleDate)]
		var changed []string
		if hadOld && item.Text != oldText {
			changed = ChangedGroups(oldText, item.Text)
		}
		blocks = append(blocks, FormatBlock(schedule.ParseText(item), changed, today, now))
	}
	return strings.Join(blocks, "\n\n")
}

func FormatBlock(parsed schedule.ParsedSchedule, changedGroupIDs []string, today time.Time, now schedule.TimeOfDay) string {
	today = schedule.NormalizeDate(today)
	date := schedule.NormalizeDate(parsed.ScheduleDate)
	dateLabel := schedule.FormatProviderDate(date)
	switch {
	case date.Equal(today):
		dateLabel = "сьогодні"
	case date.Equal(today.AddDate(0, 0, 1)):
		dateLabel = "завтра"
	}

	parts := []string{"<b>Графік відключень на " + dateLabel + "</b>"}
	if parsed.InfoTimestamp != "" {
		parts = append(parts, "<i>Станом на "+EscapeHTML(parsed.InfoTimestamp)+"</i>")
	}
	if len(changedGroupIDs) > 0 {
		parts = append(parts, "<b>Змінено:</b> "+EscapeHTML(strings.Join(changedGroupIDs, ", ")))
	}

	if len(parsed.Groups) == 0 {
		bodyLines := filteredBodyLines(parsed.RawText)
		if len(bodyLines) > 0 {
			parts = append(parts, "")
			for _, line := range bodyLines {
				parts = append(parts, EscapeHTML(line))
			}
		}
		return strings.Join(parts, "\n")
	}

	var withOutages, withoutOutages []schedule.GroupEntry
	for _, group := range parsed.Groups {
		if len(group.Outages) == 0 {
			withoutOutages = append(withoutOutages, group)
		} else {
			withOutages = append(withOutages, group)
		}
	}

	extraLines := extraScheduleLines(parsed.RawText)
	if date.Equal(today) {
		current := currentlyOff(withOutages, now)
		if len(current) > 0 {
			parts = append(parts, "", "<b>Зараз без світла:</b>")
			for _, entry := range current {
				parts = append(parts, "<b>"+entry.ID+"</b>: "+formatIntervals(entry.Outages))
			}
		}

		later := laterToday(withOutages, now)
		if len(later) > 0 {
			parts = append(parts, "", "<b>Далі сьогодні:</b>")
			for _, entry := range later {
				parts = append(parts, "<b>"+entry.ID+"</b>: "+formatIntervals(entry.Outages))
			}
		}
	} else if len(withOutages) > 0 {
		parts = append(parts, "", "<b>Відключення:</b>")
		for _, entry := range withOutages {
			parts = append(parts, "<b>"+entry.ID+"</b>: "+formatIntervals(entry.Outages))
		}
	}

	if len(withoutOutages) > 0 {
		ids := make([]string, 0, len(withoutOutages))
		for _, entry := range withoutOutages {
			ids = append(ids, entry.ID)
		}
		parts = append(parts, "", "<b>Без відключень:</b> "+strings.Join(ids, ", "))
	}

	if len(extraLines) > 0 {
		parts = append(parts, "")
		for _, line := range extraLines {
			parts = append(parts, EscapeHTML(line))
		}
	}

	return strings.Join(parts, "\n")
}

func groupLines(text string) map[string]string {
	result := make(map[string]string)
	for _, match := range schedule.GroupLineRegex.FindAllStringSubmatch(text, -1) {
		result[match[2]] = match[1]
	}
	return result
}

func isCurrentlyInInterval(interval schedule.TimeInterval, current schedule.TimeOfDay) bool {
	from := interval.From.Minutes
	to := interval.To.Minutes
	now := current.Minutes
	if from <= to {
		return now >= from && now < to
	}
	return !(now >= to && now < from)
}

func currentlyOff(entries []schedule.GroupEntry, now schedule.TimeOfDay) []schedule.GroupEntry {
	var result []schedule.GroupEntry
	for _, entry := range entries {
		var intervals []schedule.TimeInterval
		for _, interval := range entry.Outages {
			if isCurrentlyInInterval(interval, now) {
				intervals = append(intervals, interval)
			}
		}
		if len(intervals) > 0 {
			result = append(result, schedule.GroupEntry{ID: entry.ID, Outages: intervals})
		}
	}
	return result
}

func laterToday(entries []schedule.GroupEntry, now schedule.TimeOfDay) []schedule.GroupEntry {
	var result []schedule.GroupEntry
	for _, entry := range entries {
		var intervals []schedule.TimeInterval
		for _, interval := range entry.Outages {
			if interval.From.Minutes > now.Minutes && !isCurrentlyInInterval(interval, now) {
				intervals = append(intervals, interval)
			}
		}
		if len(intervals) > 0 {
			result = append(result, schedule.GroupEntry{ID: entry.ID, Outages: intervals})
		}
	}
	return result
}

func formatIntervals(intervals []schedule.TimeInterval) string {
	parts := make([]string, 0, len(intervals))
	for _, interval := range intervals {
		parts = append(parts, interval.From.Format()+"-"+interval.To.Format())
	}
	return strings.Join(parts, ", ")
}

func filteredBodyLines(text string) []string {
	var result []string
	for _, line := range strings.Split(text, "\n") {
		if line == "" || strings.HasPrefix(line, "Графік погодинних") || strings.HasPrefix(line, "Інформація станом") {
			continue
		}
		result = append(result, line)
	}
	return result
}

func extraScheduleLines(text string) []string {
	var result []string
	for _, line := range strings.Split(text, "\n") {
		if line == "" ||
			strings.HasPrefix(line, "Графік погодинних") ||
			strings.HasPrefix(line, "Інформація станом") ||
			schedule.GroupLineRegex.MatchString(line) {
			continue
		}
		result = append(result, line)
	}
	return result
}
