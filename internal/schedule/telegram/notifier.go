package telegram

import (
	"context"
	"log"
	"strings"

	sharedtelegram "github.com/sl4wa/outages-bot/internal/shared/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender interface {
	SendHTML(ctx context.Context, chatID int64, text string) error
}

type BotSender struct {
	Bot *tgbotapi.BotAPI
}

func (s BotSender) SendHTML(ctx context.Context, chatID int64, text string) error {
	_ = ctx
	return sharedtelegram.SendHTML(s.Bot, chatID, text)
}

type SubscriberStore interface {
	ChatIDs() ([]int64, error)
}

type UserNotifier struct {
	Sender      Sender
	Subscribers SubscriberStore
	Logger      *log.Logger
}

func (n UserNotifier) Notify(ctx context.Context, message string) (bool, error) {
	if strings.TrimSpace(message) == "" {
		return false, nil
	}
	chatIDs, err := n.Subscribers.ChatIDs()
	if err != nil {
		n.logger().Printf("failed to enumerate subscribers: %v", err)
		return false, err
	}

	notified := false
	for _, chatID := range chatIDs {
		if err := n.Sender.SendHTML(ctx, chatID, message); err != nil {
			n.logger().Printf("failed to send to user %d: %v", chatID, err)
			continue
		}
		notified = true
	}
	return notified, nil
}

func (n UserNotifier) logger() *log.Logger {
	if n.Logger != nil {
		return n.Logger
	}
	return log.Default()
}
