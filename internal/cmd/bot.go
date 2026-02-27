package cli

import (
	"context"
	"log"
	"outages-bot/internal/application/subscription"
	"outages-bot/internal/domain"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Conversation steps.
const (
	StepNone             = 0
	StepSearchStreet     = 1
	StepSaveSubscription = 2
)

// ConversationState holds the state of a user's conversation.
type ConversationState struct {
	Step               int
	SelectedStreetID   int
	SelectedStreetName string
	LastActivity       time.Time
}

// BotRunner manages the Telegram bot with conversation state machine.
type BotRunner struct {
	bot                     *tgbotapi.BotAPI
	searchStreetService     *subscription.SearchStreetService
	showSubscriptionService *subscription.ShowSubscriptionService
	saveSubscriptionService *subscription.SaveSubscriptionService
	userRepo                domain.UserRepository
	logger                  *log.Logger
	clock                   func() time.Time
	ttl                     time.Duration

	conversations map[int64]*ConversationState
	mu            sync.RWMutex

	cancel     context.CancelFunc
	cancelOnce sync.Once
}

// BotRunnerConfig holds configuration for BotRunner.
type BotRunnerConfig struct {
	Bot                     *tgbotapi.BotAPI
	SearchStreetService     *subscription.SearchStreetService
	ShowSubscriptionService *subscription.ShowSubscriptionService
	SaveSubscriptionService *subscription.SaveSubscriptionService
	UserRepo                domain.UserRepository
	Logger                  *log.Logger
	Clock                   func() time.Time
	CleanupTicker           <-chan time.Time
	TTL                     time.Duration
}

// NewBotRunner creates a new BotRunner with the given configuration.
func NewBotRunner(cfg BotRunnerConfig) *BotRunner {
	if cfg.Clock == nil {
		cfg.Clock = time.Now
	}
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}
	if cfg.TTL == 0 {
		cfg.TTL = 30 * time.Minute
	}

	ctx, cancel := context.WithCancel(context.Background())
	br := &BotRunner{
		bot:                     cfg.Bot,
		searchStreetService:     cfg.SearchStreetService,
		showSubscriptionService: cfg.ShowSubscriptionService,
		saveSubscriptionService: cfg.SaveSubscriptionService,
		userRepo:                cfg.UserRepo,
		logger:                  cfg.Logger,
		clock:                   cfg.Clock,
		ttl:                     cfg.TTL,
		conversations:           make(map[int64]*ConversationState),
		cancel:                  cancel,
	}

	if cfg.CleanupTicker != nil {
		go br.runCleanup(ctx, cfg.CleanupTicker)
	} else {
		ticker := time.NewTicker(10 * time.Minute)
		go func() {
			br.runCleanup(ctx, ticker.C)
			ticker.Stop()
		}()
	}

	return br
}

// Close stops the cleanup goroutine.
func (br *BotRunner) Close() {
	br.cancelOnce.Do(func() {
		br.cancel()
	})
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
	text := msg.Text

	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			br.handleStart(chatID)
			return
		case "stop":
			br.handleStop(chatID)
			return
		case "subscription":
			br.handleSubscription(chatID)
			return
		}
	}

	br.mu.RLock()
	state, exists := br.conversations[chatID]
	br.mu.RUnlock()

	if !exists || state == nil {
		return
	}

	switch state.Step {
	case StepSearchStreet:
		br.handleSearchStreet(chatID, text)
	case StepSaveSubscription:
		br.handleSaveSubscription(chatID, text)
	}
}

func (br *BotRunner) handleStart(chatID int64) {
	msg := br.showSubscriptionService.Handle(chatID)
	br.sendMessage(chatID, msg)

	br.mu.Lock()
	br.conversations[chatID] = &ConversationState{
		Step:         StepSearchStreet,
		LastActivity: br.clock(),
	}
	br.mu.Unlock()
}

func (br *BotRunner) handleStop(chatID int64) {
	// Clear any active conversation
	br.mu.Lock()
	delete(br.conversations, chatID)
	br.mu.Unlock()

	removed, err := br.userRepo.Remove(chatID)
	if err != nil {
		br.logger.Printf("error removing user %d: %v", chatID, err)
		br.sendMessage(chatID, "Сталася помилка. Спробуйте пізніше.")
		return
	}

	if removed {
		br.sendMessage(chatID, "Ви успішно відписалися від сповіщень про відключення електроенергії.")
	} else {
		br.sendMessage(chatID, "Ви не маєте активної підписки.")
	}
}

func (br *BotRunner) handleSubscription(chatID int64) {
	msg, err := br.showSubscriptionService.ShowCurrent(chatID)
	if err != nil {
		br.logger.Printf("error finding user %d: %v", chatID, err)
		br.sendMessage(chatID, "Сталася помилка. Спробуйте пізніше.")
		return
	}

	br.sendMessage(chatID, msg)
}

func (br *BotRunner) handleSearchStreet(chatID int64, text string) {
	br.mu.RLock()
	state := br.conversations[chatID]
	br.mu.RUnlock()

	if br.isExpired(state) {
		br.mu.Lock()
		delete(br.conversations, chatID)
		br.mu.Unlock()
		br.sendMessage(chatID, "Сесія застаріла. Будь ласка, почніть знову з /start")
		return
	}

	result := br.searchStreetService.Handle(text)

	if result.HasMultipleOptions() {
		rows := make([][]tgbotapi.KeyboardButton, len(result.StreetOptions))
		for i, street := range result.StreetOptions {
			rows[i] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(street.Name))
		}
		markup := tgbotapi.NewReplyKeyboard(rows...)
		markup.ResizeKeyboard = true
		markup.OneTimeKeyboard = true
		br.sendMessageWithReplyMarkup(chatID, result.Message, markup)
	} else {
		br.sendMessage(chatID, result.Message)
	}

	if result.HasExactMatch() {
		br.mu.Lock()
		br.conversations[chatID] = &ConversationState{
			Step:               StepSaveSubscription,
			SelectedStreetID:   *result.SelectedStreetID,
			SelectedStreetName: *result.SelectedStreetName,
			LastActivity:       br.clock(),
		}
		br.mu.Unlock()
	} else {
		br.mu.Lock()
		state.LastActivity = br.clock()
		br.mu.Unlock()
	}
}

func (br *BotRunner) handleSaveSubscription(chatID int64, text string) {
	br.mu.RLock()
	state := br.conversations[chatID]
	br.mu.RUnlock()

	if br.isExpired(state) {
		br.mu.Lock()
		delete(br.conversations, chatID)
		br.mu.Unlock()
		br.sendMessage(chatID, "Сесія застаріла. Будь ласка, почніть знову з /start")
		return
	}

	if state.SelectedStreetID == 0 || state.SelectedStreetName == "" {
		br.mu.Lock()
		delete(br.conversations, chatID)
		br.mu.Unlock()
		br.sendMessage(chatID, "Сесія застаріла. Будь ласка, почніть знову з /start")
		return
	}

	result := br.saveSubscriptionService.Handle(chatID, state.SelectedStreetID, state.SelectedStreetName, text)
	br.sendMessage(chatID, result.Message)

	if result.Success {
		br.mu.Lock()
		delete(br.conversations, chatID)
		br.mu.Unlock()
	} else {
		br.mu.Lock()
		state.LastActivity = br.clock()
		br.mu.Unlock()
	}
}

func (br *BotRunner) isExpired(state *ConversationState) bool {
	if state == nil {
		return true
	}
	return br.clock().Sub(state.LastActivity) > br.ttl
}

func (br *BotRunner) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := br.bot.Send(msg); err != nil {
		br.logger.Printf("failed to send message to %d: %v", chatID, err)
	}
}

func (br *BotRunner) sendMessageWithReplyMarkup(chatID int64, text string, markup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = markup
	if _, err := br.bot.Send(msg); err != nil {
		br.logger.Printf("failed to send message to %d: %v", chatID, err)
	}
}

func (br *BotRunner) runCleanup(ctx context.Context, tick <-chan time.Time) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			br.cleanupExpired()
		}
	}
}

func (br *BotRunner) cleanupExpired() {
	br.mu.Lock()
	defer br.mu.Unlock()
	now := br.clock()
	for chatID, state := range br.conversations {
		if now.Sub(state.LastActivity) > br.ttl {
			delete(br.conversations, chatID)
		}
	}
}

// GetState returns the conversation state for testing.
func (br *BotRunner) GetState(chatID int64) *ConversationState {
	br.mu.RLock()
	defer br.mu.RUnlock()
	return br.conversations[chatID]
}
