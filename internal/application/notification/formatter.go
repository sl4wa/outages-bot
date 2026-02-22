package notification

import (
	"fmt"
	"outages-bot/internal/application"
	"strings"
)

// FormatNotification formats a notification DTO into an HTML message for Telegram.
func FormatNotification(dto application.NotificationSenderDTO) string {
	buildings := strings.Join(dto.Buildings, ", ")

	return fmt.Sprintf(
		"Поточні відключення:\nМісто: %s\nВулиця: %s\n<b>%s – %s</b>\nКоментар: %s\nБудинки: %s",
		dto.City,
		dto.StreetName,
		dto.Start.Format("2006-01-02 15:04"),
		dto.End.Format("2006-01-02 15:04"),
		dto.Comment,
		buildings,
	)
}
