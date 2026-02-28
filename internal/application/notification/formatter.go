package notification

import (
	"fmt"
	"outages-bot/internal/domain"
	"strings"
	"time"
)

// FormatNotification formats outage data into an HTML message for Telegram.
func FormatNotification(city, streetName string, buildings []string, start, end time.Time, comment string) string {
	return fmt.Sprintf(
		"Поточні відключення:\nМісто: %s\nВулиця: %s\n<b>%s – %s</b>\nКоментар: %s\nБудинки: %s",
		city,
		streetName,
		start.Format("2006-01-02 15:04"),
		end.Format("2006-01-02 15:04"),
		comment,
		strings.Join(buildings, ", "),
	)
}

// formatOutageNotification formats an Outage into an HTML message.
func formatOutageNotification(outage *domain.Outage) string {
	return FormatNotification(
		outage.Address.City,
		outage.Address.StreetName,
		outage.Address.Buildings,
		outage.Period.StartDate,
		outage.Period.EndDate,
		outage.Description.Value,
	)
}
