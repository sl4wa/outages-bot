package telegram

import (
	"fmt"
	"outages-bot/internal/application"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramUserInfoProvider retrieves user info from the Telegram API.
type TelegramUserInfoProvider struct {
	bot *tgbotapi.BotAPI
}

// NewTelegramUserInfoProvider creates a new TelegramUserInfoProvider.
func NewTelegramUserInfoProvider(bot *tgbotapi.BotAPI) *TelegramUserInfoProvider {
	return &TelegramUserInfoProvider{bot: bot}
}

// GetUserInfo retrieves user info for the given chat ID.
func (p *TelegramUserInfoProvider) GetUserInfo(chatID int64) (application.UserInfoDTO, error) {
	chatConfig := tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	}

	chat, err := p.bot.GetChat(chatConfig)
	if err != nil {
		return application.UserInfoDTO{}, fmt.Errorf("failed to get user info for chat %d: %w", chatID, err)
	}

	return application.UserInfoDTO{
		ChatID:    chat.ID,
		Username:  chat.UserName,
		FirstName: chat.FirstName,
		LastName:  chat.LastName,
	}, nil
}
