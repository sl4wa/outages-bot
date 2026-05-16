package telegram

import (
	"fmt"
	"outages-bot/internal/users"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// InfoProvider retrieves user info from the Telegram API.
type InfoProvider struct {
	bot *tgbotapi.BotAPI
}

// NewInfoProvider creates a new InfoProvider.
func NewInfoProvider(bot *tgbotapi.BotAPI) *InfoProvider {
	return &InfoProvider{bot: bot}
}

// GetUserInfo retrieves user info for the given chat ID.
func (p *InfoProvider) GetUserInfo(chatID int64) (users.Info, error) {
	chatConfig := tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	}

	chat, err := p.bot.GetChat(chatConfig)
	if err != nil {
		return users.Info{}, fmt.Errorf("failed to get user info for chat %d: %w", chatID, err)
	}

	return users.Info{
		ChatID:    chat.ID,
		Username:  chat.UserName,
		FirstName: chat.FirstName,
		LastName:  chat.LastName,
	}, nil
}
