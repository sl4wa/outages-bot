package notification

import (
	"outages-bot/internal/application"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatNotification_Standard(t *testing.T) {
	dto := application.NotificationSenderDTO{
		UserID:     1,
		City:       "Львів",
		StreetName: "Стрийська",
		Buildings:  []string{"10", "12", "14"},
		Start:      time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
		End:        time.Date(2024, 1, 15, 16, 0, 0, 0, time.UTC),
		Comment:    "Планове відключення",
	}

	result := FormatNotification(dto)
	expected := "Поточні відключення:\nМісто: Львів\nВулиця: Стрийська\n<b>2024-01-15 08:00 – 2024-01-15 16:00</b>\nКоментар: Планове відключення\nБудинки: 10, 12, 14"
	assert.Equal(t, expected, result)
}

func TestFormatNotification_SpecialCharsInStreetName(t *testing.T) {
	dto := application.NotificationSenderDTO{
		UserID:     1,
		City:       "City",
		StreetName: "Street <test> & name",
		Buildings:  []string{"1"},
		Start:      time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		End:        time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		Comment:    "comment",
	}

	result := FormatNotification(dto)
	assert.Contains(t, result, "Street <test> & name")
}

func TestFormatNotification_SpecialCharsInComment(t *testing.T) {
	dto := application.NotificationSenderDTO{
		UserID:     1,
		City:       "City",
		StreetName: "Street",
		Buildings:  []string{"1"},
		Start:      time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		End:        time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		Comment:    "test <b>bold</b> & stuff",
	}

	result := FormatNotification(dto)
	assert.Contains(t, result, "test <b>bold</b> & stuff")
}

func TestFormatNotification_SpecialCharsInBuildings(t *testing.T) {
	dto := application.NotificationSenderDTO{
		UserID:     1,
		City:       "City",
		StreetName: "Street",
		Buildings:  []string{"<1>", "2&3"},
		Start:      time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		End:        time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		Comment:    "comment",
	}

	result := FormatNotification(dto)
	assert.Contains(t, result, "<1>, 2&3")
}

func TestFormatNotification_BoldTagsPreserved(t *testing.T) {
	dto := application.NotificationSenderDTO{
		UserID:     1,
		City:       "City",
		StreetName: "Street",
		Buildings:  []string{"1"},
		Start:      time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		End:        time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
		Comment:    "comment",
	}

	result := FormatNotification(dto)
	assert.Contains(t, result, "<b>")
	assert.Contains(t, result, "</b>")
}
