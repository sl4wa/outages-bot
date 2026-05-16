package telegram

import (
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ErrRecipientUnavailable = errors.New("recipient unavailable")

func SendHTML(bot *tgbotapi.BotAPI, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true

	_, err := bot.Send(msg)
	if err != nil {
		return NormalizeError(err)
	}
	return nil
}

func NormalizeError(err error) error {
	if err == nil {
		return nil
	}
	var apiErr *tgbotapi.Error
	if errors.As(err, &apiErr) {
		if apiErr.Code == 403 || strings.Contains(strings.ToLower(apiErr.Message), "forbidden") {
			return fmt.Errorf("%w: %w", ErrRecipientUnavailable, apiErr)
		}
		return fmt.Errorf("telegram API error: %w", apiErr)
	}
	return fmt.Errorf("telegram send: %w", err)
}
