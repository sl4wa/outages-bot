package telegram

import (
	"outages-bot/internal/application"
	"outages-bot/internal/application/notification"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramNotificationSender sends notifications via the Telegram Bot API.
type TelegramNotificationSender struct {
	bot *tgbotapi.BotAPI
}

// NewTelegramNotificationSender creates a new TelegramNotificationSender.
func NewTelegramNotificationSender(bot *tgbotapi.BotAPI) *TelegramNotificationSender {
	return &TelegramNotificationSender{bot: bot}
}

// Send sends a notification to the user via Telegram.
func (s *TelegramNotificationSender) Send(dto application.NotificationSenderDTO) error {
	text := notification.FormatNotification(dto)
	msg := tgbotapi.NewMessage(dto.UserID, text)
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

		return &application.NotificationSendError{
			UserID:  dto.UserID,
			Code:    code,
			Message: message,
		}
	}
	return nil
}
