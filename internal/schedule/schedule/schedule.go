package schedule

import (
	"regexp"
	"sort"
	"strconv"
	"time"
)

type Snapshot struct {
	ScheduleDate time.Time
	UpdatedAt    *time.Time
	Text         string
	SourceID     *int
}

type TimeOfDay struct {
	Minutes int
}

type TimeInterval struct {
	From TimeOfDay
	To   TimeOfDay
}

type GroupEntry struct {
	ID      string
	Outages []TimeInterval
}

type ParsedSchedule struct {
	ScheduleDate  time.Time
	InfoTimestamp string
	InfoTime      *TimeOfDay
	InfoDate      *time.Time
	Groups        []GroupEntry
	RawText       string
}

var (
	GroupLineRegex      = regexp.MustCompile(`(?m)^(Група (\d+\.\d+)\..*)$`)
	infoExtractRegex    = regexp.MustCompile(`(?m)^Інформація станом на (\d{2}:\d{2})\s+(\d{2}\.\d{2}\.\d{4})$`)
	outageIntervalRegex = regexp.MustCompile(`з (\d{2}:\d{2}) до (\d{2}:\d{2})`)
	groupIDPartsRegex   = regexp.MustCompile(`^(\d+)\.(\d+)$`)
	providerDateLayout  = "02.01.2006"
	providerTimeLayout  = "15:04"
	stateDateLayout     = "2006-01-02"
)

func NormalizeDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func ParseProviderDate(value string) (time.Time, error) {
	t, err := time.ParseInLocation(providerDateLayout, value, time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	return NormalizeDate(t), nil
}

func FormatProviderDate(t time.Time) string {
	return NormalizeDate(t).Format(providerDateLayout)
}

func ParseStateDate(value string) (time.Time, error) {
	t, err := time.ParseInLocation(stateDateLayout, value, time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	return NormalizeDate(t), nil
}

func FormatStateDate(t time.Time) string {
	return NormalizeDate(t).Format(stateDateLayout)
}

func ParseProviderTime(value string) (TimeOfDay, error) {
	if value == "24:00" {
		return TimeOfDay{Minutes: 24 * 60}, nil
	}
	t, err := time.Parse(providerTimeLayout, value)
	if err != nil {
		return TimeOfDay{}, err
	}
	return TimeOfDay{Minutes: t.Hour()*60 + t.Minute()}, nil
}

func TimeOfDayFromTime(t time.Time) TimeOfDay {
	return TimeOfDay{Minutes: t.Hour()*60 + t.Minute()}
}

func (t TimeOfDay) Format() string {
	if t.Minutes == 24*60 {
		return "24:00"
	}
	minutes := t.Minutes % (24 * 60)
	if minutes < 0 {
		minutes += 24 * 60
	}
	return time.Date(2000, 1, 1, minutes/60, minutes%60, 0, 0, time.UTC).Format(providerTimeLayout)
}

func SelectLatestForDate(schedules []Snapshot, date time.Time) (Snapshot, bool) {
	date = NormalizeDate(date)
	var candidates []Snapshot
	for _, item := range schedules {
		if NormalizeDate(item.ScheduleDate).Equal(date) {
			candidates = append(candidates, item)
		}
	}
	if len(candidates) == 0 {
		return Snapshot{}, false
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		var left, right time.Time
		if candidates[i].UpdatedAt != nil {
			left = *candidates[i].UpdatedAt
		}
		if candidates[j].UpdatedAt != nil {
			right = *candidates[j].UpdatedAt
		}
		if !left.Equal(right) {
			return left.Before(right)
		}
		leftID, rightID := 0, 0
		if candidates[i].SourceID != nil {
			leftID = *candidates[i].SourceID
		}
		if candidates[j].SourceID != nil {
			rightID = *candidates[j].SourceID
		}
		return leftID < rightID
	})
	return candidates[len(candidates)-1], true
}

func ParseText(snapshot Snapshot) ParsedSchedule {
	var infoTimestamp string
	var infoTime *TimeOfDay
	var infoDate *time.Time
	if match := infoExtractRegex.FindStringSubmatch(snapshot.Text); len(match) == 3 {
		infoTimestamp = match[1] + " " + match[2]
		if parsedTime, err := ParseProviderTime(match[1]); err == nil {
			infoTime = &parsedTime
		}
		if parsedDate, err := ParseProviderDate(match[2]); err == nil {
			infoDate = &parsedDate
		}
	}

	matches := GroupLineRegex.FindAllStringSubmatch(snapshot.Text, -1)
	groups := make([]GroupEntry, 0, len(matches))
	for _, match := range matches {
		entry := GroupEntry{ID: match[2]}
		for _, intervalMatch := range outageIntervalRegex.FindAllStringSubmatch(match[1], -1) {
			from, fromErr := ParseProviderTime(intervalMatch[1])
			to, toErr := ParseProviderTime(intervalMatch[2])
			if fromErr == nil && toErr == nil {
				entry.Outages = append(entry.Outages, TimeInterval{From: from, To: to})
			}
		}
		groups = append(groups, entry)
	}

	sort.SliceStable(groups, func(i, j int) bool {
		ia, ib := splitGroupID(groups[i].ID)
		ja, jb := splitGroupID(groups[j].ID)
		if ia != ja {
			return ia < ja
		}
		return ib < jb
	})

	return ParsedSchedule{
		ScheduleDate:  NormalizeDate(snapshot.ScheduleDate),
		InfoTimestamp: infoTimestamp,
		InfoTime:      infoTime,
		InfoDate:      infoDate,
		Groups:        groups,
		RawText:       snapshot.Text,
	}
}

func splitGroupID(value string) (int, int) {
	match := groupIDPartsRegex.FindStringSubmatch(value)
	if len(match) != 3 {
		return 0, 0
	}
	a, _ := strconv.Atoi(match[1])
	b, _ := strconv.Atoi(match[2])
	return a, b
}
