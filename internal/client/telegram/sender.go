package telegram

import (
	"outages-bot/internal/application/notifier"

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
		code := 0
		message := err.Error()

		// Extract HTTP status code from tgbotapi.Error
		if apiErr, ok := err.(*tgbotapi.Error); ok {
			code = apiErr.Code
			message = apiErr.Message
		}

		return &notifier.NotificationSendError{
			UserID:  userID,
			Code:    code,
			Message: message,
		}
	}
	return nil
}
