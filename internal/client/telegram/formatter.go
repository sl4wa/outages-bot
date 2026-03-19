package telegram

import (
	"fmt"
	"outages-bot/internal/application/notifier"
	"strings"
)

func formatNotification(c notifier.NotificationContent) string {
	return fmt.Sprintf(
		"Поточні відключення:\nМісто: %s\nВулиця: %s\n<b>%s – %s</b>\nКоментар: %s\nБудинки: %s",
		c.City,
		c.StreetName,
		c.Start.Format("2006-01-02 15:04"),
		c.End.Format("2006-01-02 15:04"),
		c.Comment,
		strings.Join(c.Buildings, ", "),
	)
}
