package telegram

import (
	"log"
	"outages-bot/internal/subscription"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotRunner adapts Telegram updates to the subscription workflow.
type BotRunner struct {
	bot      *tgbotapi.BotAPI
	workflow *subscription.Workflow
	logger   *log.Logger
}

// BotRunnerConfig holds configuration for BotRunner.
type BotRunnerConfig struct {
	Bot      *tgbotapi.BotAPI
	Workflow *subscription.Workflow
	Logger   *log.Logger
}

// NewBotRunner creates a new BotRunner with the given configuration.
func NewBotRunner(cfg BotRunnerConfig) *BotRunner {
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}

	br := &BotRunner{
		bot:      cfg.Bot,
		workflow: cfg.Workflow,
		logger:   cfg.Logger,
	}

	return br
}

// Run starts the bot polling loop.
func (br *BotRunner) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := br.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		br.handleMessage(update.Message)
	}
}

// HandleMessage processes a single message (exported for testing).
func (br *BotRunner) HandleMessage(msg *tgbotapi.Message) {
	br.handleMessage(msg)
}

func (br *BotRunner) handleMessage(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	cmd := subscription.Command{Kind: subscription.CommandText, Text: msg.Text}

	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			cmd = subscription.Command{Kind: subscription.CommandStart}
		case "stop":
			cmd = subscription.Command{Kind: subscription.CommandStop}
		case "subscription":
			cmd = subscription.Command{Kind: subscription.CommandSubscription}
		}
	}

	response := br.workflow.Handle(chatID, cmd)
	br.sendResponse(chatID, response)
}

func (br *BotRunner) sendMessage(chatID int64, text string, markup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	if markup != nil {
		msg.ReplyMarkup = markup
	}
	if _, err := br.bot.Send(msg); err != nil {
		br.logger.Printf("failed to send message to %d: %v", chatID, err)
	}
}

func (br *BotRunner) sendResponse(chatID int64, response subscription.Response) {
	if response.Text == "" && len(response.StreetOptions) == 0 {
		return
	}
	if response.Err != nil {
		br.logger.Printf("subscription error for user %d: %v", chatID, response.Err)
	}

	var markup interface{}
	if len(response.StreetOptions) > 0 {
		markup = streetOptionsKeyboard(response.StreetOptions)
	}
	br.sendMessage(chatID, response.Text, markup)
}

func streetOptionsKeyboard(names []string) interface{} {
	rows := make([][]tgbotapi.KeyboardButton, len(names))
	for i, name := range names {
		rows[i] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(name))
	}
	markup := tgbotapi.NewReplyKeyboard(rows...)
	markup.ResizeKeyboard = true
	markup.OneTimeKeyboard = true
	return markup
}

// GetState returns the conversation state for testing.
func (br *BotRunner) GetState(chatID int64) *subscription.State {
	return br.workflow.GetState(chatID)
}
