package notification

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatNotification_Standard(t *testing.T) {
	result := FormatNotification(
		"Львів", "Стрийська", []string{"10", "12", "14"},
		time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 16, 0, 0, 0, time.UTC),
		"Планове відключення",
	)
	expected := "Поточні відключення:\nМісто: Львів\nВулиця: Стрийська\n<b>2024-01-15 08:00 – 2024-01-15 16:00</b>\nКоментар: Планове відключення\nБудинки: 10, 12, 14"
	assert.Equal(t, expected, result)
}

func TestFormatNotification_SpecialCharsInStreetName(t *testing.T) {
	result := FormatNotification(
		"City", "Street <test> & name", []string{"1"},
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		"comment",
	)
	assert.Contains(t, result, "Street <test> & name")
}

func TestFormatNotification_SpecialCharsInComment(t *testing.T) {
	result := FormatNotification(
		"City", "Street", []string{"1"},
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		"test <b>bold</b> & stuff",
	)
	assert.Contains(t, result, "test <b>bold</b> & stuff")
}

func TestFormatNotification_SpecialCharsInBuildings(t *testing.T) {
	result := FormatNotification(
		"City", "Street", []string{"<1>", "2&3"},
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		"comment",
	)
	assert.Contains(t, result, "<1>, 2&3")
}

func TestFormatNotification_BoldTagsPreserved(t *testing.T) {
	result := FormatNotification(
		"City", "Street", []string{"1"},
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		"comment",
	)
	assert.Contains(t, result, "<b>")
	assert.Contains(t, result, "</b>")
}
