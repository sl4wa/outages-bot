package telegram

import (
	"errors"
	"fmt"
	"github.com/sl4wa/outages-bot/internal/outage/notifier"
	sharedtelegram "github.com/sl4wa/outages-bot/internal/shared/telegram"

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
func (s *NotificationSender) Send(userID int64, content notifier.Content) error {
	text := formatNotification(content)
	err := sharedtelegram.SendHTML(s.bot, userID, text)
	if err != nil {
		if errors.Is(err, sharedtelegram.ErrRecipientUnavailable) {
			return notifier.ErrRecipientUnavailable
		}
		return fmt.Errorf("telegram send: %w", err)
	}
	return nil
}
