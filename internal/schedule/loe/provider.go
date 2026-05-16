package loe

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sl4wa/outages-bot/internal/schedule/schedule"
)

type PayloadLoader func(context.Context) (string, error)

type Provider struct {
	LoadPayload       PayloadLoader
	OnPayloadAccepted func() error
}

func (p Provider) GetSchedules(ctx context.Context) ([]schedule.Snapshot, error) {
	if p.LoadPayload == nil {
		return nil, fmt.Errorf("payload loader is required")
	}
	payload, err := p.LoadPayload(ctx)
	if err != nil {
		return nil, err
	}
	result, err := parseSchedules(payload)
	if err != nil {
		return nil, err
	}
	if p.OnPayloadAccepted != nil {
		if err := p.OnPayloadAccepted(); err != nil {
			return nil, err
		}
	}
	return result, nil
}

type menuObject struct {
	HydraMember []menuMember `json:"hydra:member"`
}

type menuMember struct {
	MenuItems []menuItem `json:"menuItems"`
}

type menuItem struct {
	ID          json.RawMessage `json:"id"`
	RawHTML     string          `json:"rawHtml"`
	Description string          `json:"description"`
	Children    []menuItem      `json:"children"`
}

func parseSchedules(payload string) ([]schedule.Snapshot, error) {
	if strings.TrimSpace(payload) == "" {
		return nil, fmt.Errorf("outage schedule payload is blank")
	}

	var members []menuMember
	var root menuObject
	if err := json.Unmarshal([]byte(payload), &root); err == nil && root.HydraMember != nil {
		members = root.HydraMember
	} else {
		if arrErr := json.Unmarshal([]byte(payload), &members); arrErr != nil {
			return nil, fmt.Errorf("outage schedule payload is not valid JSON: %.200s: %w", payload, arrErr)
		}
	}
	if members == nil {
		return nil, fmt.Errorf("outage schedule payload is missing hydra:member")
	}

	var result []schedule.Snapshot
	for _, member := range members {
		for _, item := range member.MenuItems {
			result = appendParsed(result, item)
			for _, child := range item.Children {
				result = appendParsed(result, child)
			}
		}
	}
	return result, nil
}

func appendParsed(result []schedule.Snapshot, item menuItem) []schedule.Snapshot {
	source := firstScheduleText(item.RawHTML, item.Description)
	if source == "" {
		return result
	}
	if parsed, ok := parseRawHTMLBlock(itemID(item.ID), source); ok {
		result = append(result, parsed)
	}
	return result
}

func firstScheduleText(rawHTML, description string) string {
	if strings.TrimSpace(rawHTML) != "" {
		return rawHTML
	}
	if strings.TrimSpace(description) != "" {
		return description
	}
	return ""
}

func parseRawHTMLBlock(sourceID *int, rawHTML string) (schedule.Snapshot, bool) {
	cleaned := cleanHTMLText(strings.TrimSpace(rawHTML))
	if cleaned == "" {
		return schedule.Snapshot{}, false
	}
	scheduleDate, ok := extractScheduleDate(cleaned)
	if !ok {
		return schedule.Snapshot{}, false
	}
	text := addBlankLineAfterHeaders(cleaned)
	return schedule.Snapshot{
		ScheduleDate: scheduleDate,
		UpdatedAt:    extractUpdatedAt(text),
		Text:         text,
		SourceID:     sourceID,
	}, true
}

var (
	headerLineRegex      = regexp.MustCompile(`(?m)^(Графік погодинних відключень на .+|Інформація станом на .+)$`)
	dateInTextRegex      = regexp.MustCompile(`(\d{2}\.\d{2}\.\d{4})`)
	updateTimestampRegex = regexp.MustCompile(`(\d{2}:\d{2})\s+(\d{2}\.\d{2}\.\d{4})`)
)

func addBlankLineAfterHeaders(text string) string {
	return headerLineRegex.ReplaceAllString(text, "$1\n")
}

func extractScheduleDate(text string) (time.Time, bool) {
	match := dateInTextRegex.FindStringSubmatch(text)
	if len(match) != 2 {
		return time.Time{}, false
	}
	parsed, err := schedule.ParseProviderDate(match[1])
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

func extractUpdatedAt(text string) *time.Time {
	match := updateTimestampRegex.FindStringSubmatch(text)
	if len(match) != 3 {
		return nil
	}
	date, err := schedule.ParseProviderDate(match[2])
	if err != nil {
		return nil
	}
	tod, err := schedule.ParseProviderTime(match[1])
	if err != nil {
		return nil
	}
	updated := time.Date(date.Year(), date.Month(), date.Day(), tod.Minutes/60, tod.Minutes%60, 0, 0, time.UTC)
	return &updated
}

func itemID(raw json.RawMessage) *int {
	if len(raw) == 0 {
		return nil
	}
	var n int
	if err := json.Unmarshal(raw, &n); err == nil {
		return &n
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if parsed, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			return &parsed
		}
	}
	return nil
}

var (
	paragraphEndRegex    = regexp.MustCompile(`(?i)</p>`)
	blockEndRegex        = regexp.MustCompile(`(?i)</div>`)
	breakTagRegex        = regexp.MustCompile(`(?i)<br\s*/?>`)
	htmlTagRegex         = regexp.MustCompile(`<[^>]+>`)
	inlineWhitespace     = regexp.MustCompile(`[\t\v\f\r ]+`)
	excessiveNewlineExpr = regexp.MustCompile(`\n{2,}`)
)

func cleanHTMLText(html string) string {
	text := strings.ReplaceAll(html, "\r\n", "\n")
	text = paragraphEndRegex.ReplaceAllString(text, "\n")
	text = blockEndRegex.ReplaceAllString(text, "\n")
	text = breakTagRegex.ReplaceAllString(text, "\n")
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = htmlTagRegex.ReplaceAllString(text, "")

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(inlineWhitespace.ReplaceAllString(line, " "))
	}
	text = strings.Join(lines, "\n")
	text = excessiveNewlineExpr.ReplaceAllString(text, "\n")
	return strings.TrimSpace(text)
}
