package telegram

import (
	"errors"
	"fmt"
	"outages-bot/internal/application/notifier"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NotificationSender sends notifications via the Telegram Bot API.
type NotificationSender struct {
	bot *tgbotapi.BotAPI
}

// NewNotificationSender creates a new NotificationSender.
func NewNotificationSender(bot *tgbotapi.BotAPI) *NotificationSender {
	return &NotificationSender{bot: bot}
}

// Send formats the notification content and sends it to the user via Telegram.
func (s *NotificationSender) Send(userID int64, content notifier.NotificationContent) error {
	text := formatNotification(content)
	msg := tgbotapi.NewMessage(userID, text)
	msg.ParseMode = "HTML"

	_, err := s.bot.Send(msg)
	if err != nil {
		var apiErr *tgbotapi.Error
		if errors.As(err, &apiErr) {
			if apiErr.Code == 403 || strings.Contains(strings.ToLower(apiErr.Message), "forbidden") {
				return notifier.ErrRecipientUnavailable
			}
			return fmt.Errorf("telegram API error: %w", apiErr)
		}
		return fmt.Errorf("telegram send: %w", err)
	}
	return nil
}
