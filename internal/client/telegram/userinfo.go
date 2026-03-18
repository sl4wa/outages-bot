package telegram

import (
	"fmt"
	"outages-bot/internal/application/users"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UserInfoProvider retrieves user info from the Telegram API.
type UserInfoProvider struct {
	bot *tgbotapi.BotAPI
}

// NewUserInfoProvider creates a new UserInfoProvider.
func NewUserInfoProvider(bot *tgbotapi.BotAPI) *UserInfoProvider {
	return &UserInfoProvider{bot: bot}
}

// GetUserInfo retrieves user info for the given chat ID.
func (p *UserInfoProvider) GetUserInfo(chatID int64) (users.UserInfoDTO, error) {
	chatConfig := tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	}

	chat, err := p.bot.GetChat(chatConfig)
	if err != nil {
		return users.UserInfoDTO{}, fmt.Errorf("failed to get user info for chat %d: %w", chatID, err)
	}

	return users.UserInfoDTO{
		ChatID:    chat.ID,
		Username:  chat.UserName,
		FirstName: chat.FirstName,
		LastName:  chat.LastName,
	}, nil
}
