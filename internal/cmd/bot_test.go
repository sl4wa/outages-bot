package cli

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"outages-bot/internal/application/subscription"
	"outages-bot/internal/domain"
	"sync"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testUserRepo struct {
	mu      sync.RWMutex
	users   map[int64]*domain.User
	findErr error
	rmErr   error
}

func newTestUserRepo() *testUserRepo {
	return &testUserRepo{users: make(map[int64]*domain.User)}
}

func (r *testUserRepo) FindAll() ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.User
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}

func (r *testUserRepo) Find(chatID int64) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.findErr != nil {
		return nil, r.findErr
	}
	return r.users[chatID], nil
}

func (r *testUserRepo) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *testUserRepo) Remove(chatID int64) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.rmErr != nil {
		return false, r.rmErr
	}
	if _, ok := r.users[chatID]; ok {
		delete(r.users, chatID)
		return true, nil
	}
	return false, nil
}

type testStreetRepo struct {
	streets []domain.Street
}

func (r *testStreetRepo) GetAllStreets() ([]domain.Street, error) {
	return r.streets, nil
}

// sentMessage captures sent messages
type sentMessage struct {
	ChatID      int64
	Text        string
	ReplyMarkup string
}

func setupBot(t *testing.T) (*BotRunner, *testUserRepo, *[]sentMessage) {
	t.Helper()

	var messages []sentMessage
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock sendMessage endpoint
		if r.URL.Path == "/bottest-token/sendMessage" {
			r.ParseForm()
			chatIDStr := r.FormValue("chat_id")
			text := r.FormValue("text")
			replyMarkup := r.FormValue("reply_markup")
			var chatID int64
			json.Unmarshal([]byte(chatIDStr), &chatID)
			mu.Lock()
			messages = append(messages, sentMessage{ChatID: chatID, Text: text, ReplyMarkup: replyMarkup})
			mu.Unlock()

			resp := tgbotapi.APIResponse{Ok: true, Result: json.RawMessage(`{"message_id":1,"chat":{"id":` + chatIDStr + `},"text":""}`)}
			json.NewEncoder(w).Encode(resp)
			return
		}
		// Default: getMe
		resp := tgbotapi.APIResponse{Ok: true, Result: json.RawMessage(`{"id":123,"is_bot":true,"first_name":"Test"}`)}
		json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(server.Close)

	api, err := tgbotapi.NewBotAPIWithAPIEndpoint("test-token", server.URL+"/bot%s/%s")
	require.NoError(t, err)

	userRepo := newTestUserRepo()
	streetRepo := &testStreetRepo{
		streets: []domain.Street{
			{ID: 1, Name: "Стрийська"},
			{ID: 2, Name: "Наукова"},
			{ID: 3, Name: "Стрілецька"},
		},
	}

	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }

	cleanupCh := make(chan time.Time)
	br := NewBotRunner(BotRunnerConfig{
		Bot:                     api,
		SearchStreetService:     subscription.NewSearchStreetService(streetRepo),
		ShowSubscriptionService: subscription.NewShowSubscriptionService(userRepo),
		SaveSubscriptionService: subscription.NewSaveSubscriptionService(userRepo),
		UserRepo:                userRepo,
		Logger:                  log.Default(),
		Clock:                   clock,
		CleanupTicker:           cleanupCh,
		TTL:                     30 * time.Minute,
	})
	t.Cleanup(br.Close)

	return br, userRepo, &messages
}

func makeMsg(chatID int64, text string) *tgbotapi.Message {
	return &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: chatID},
		Text: text,
	}
}

func makeCmd(chatID int64, command string) *tgbotapi.Message {
	return &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: chatID},
		Text: "/" + command,
		Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len("/" + command)},
		},
	}
}

func TestBot_StartSetsSearchStreetStep(t *testing.T) {
	br, _, _ := setupBot(t)
	br.HandleMessage(makeCmd(100, "start"))
	state := br.GetState(100)
	require.NotNil(t, state)
	assert.Equal(t, StepSearchStreet, state.Step)
}

func TestBot_SearchStreet_ExactMatch_AdvancesToSave(t *testing.T) {
	br, _, _ := setupBot(t)
	br.HandleMessage(makeCmd(100, "start"))
	br.HandleMessage(makeMsg(100, "Наукова"))
	state := br.GetState(100)
	require.NotNil(t, state)
	assert.Equal(t, StepSaveSubscription, state.Step)
	assert.Equal(t, 2, state.SelectedStreetID)
	assert.Equal(t, "Наукова", state.SelectedStreetName)
}

func TestBot_SearchStreet_MultipleMatch_StaysInSearch(t *testing.T) {
	br, _, msgs := setupBot(t)
	br.HandleMessage(makeCmd(100, "start"))
	br.HandleMessage(makeMsg(100, "Стр"))
	state := br.GetState(100)
	require.NotNil(t, state)
	assert.Equal(t, StepSearchStreet, state.Step)

	// Verify reply keyboard with street options was sent
	var replyMarkupJSON string
	for _, m := range *msgs {
		if m.ChatID == 100 && m.Text == "Будь ласка, оберіть вулицю:" {
			replyMarkupJSON = m.ReplyMarkup
		}
	}
	require.NotEmpty(t, replyMarkupJSON, "expected ReplyMarkup to be set")

	var keyboard tgbotapi.ReplyKeyboardMarkup
	err := json.Unmarshal([]byte(replyMarkupJSON), &keyboard)
	require.NoError(t, err)
	assert.True(t, keyboard.ResizeKeyboard)
	assert.True(t, keyboard.OneTimeKeyboard)
	require.Len(t, keyboard.Keyboard, 2)
	assert.Equal(t, "Стрийська", keyboard.Keyboard[0][0].Text)
	assert.Equal(t, "Стрілецька", keyboard.Keyboard[1][0].Text)
}

func TestBot_SaveSubscription_ValidInput_CompletesFlow(t *testing.T) {
	br, userRepo, _ := setupBot(t)
	br.HandleMessage(makeCmd(100, "start"))
	br.HandleMessage(makeMsg(100, "Наукова"))
	br.HandleMessage(makeMsg(100, "10"))

	// Conversation should be cleared
	state := br.GetState(100)
	assert.Nil(t, state)

	// User should be saved
	user, err := userRepo.Find(100)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Наукова", user.Address.StreetName)
	assert.Equal(t, "10", user.Address.Building)
}

func TestBot_SaveSubscription_InvalidInput_StaysInSave(t *testing.T) {
	br, _, _ := setupBot(t)
	br.HandleMessage(makeCmd(100, "start"))
	br.HandleMessage(makeMsg(100, "Наукова"))
	br.HandleMessage(makeMsg(100, "invalid!"))

	state := br.GetState(100)
	require.NotNil(t, state)
	assert.Equal(t, StepSaveSubscription, state.Step)
}

func TestBot_StopWithSubscription_RemovesUser(t *testing.T) {
	br, userRepo, msgs := setupBot(t)
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	userRepo.users[100] = &domain.User{ID: 100, Address: addr}

	br.HandleMessage(makeCmd(100, "stop"))

	_, err := userRepo.Find(100)
	require.NoError(t, err)

	// Check message
	found := false
	for _, m := range *msgs {
		if m.ChatID == 100 && m.Text == "Ви успішно відписалися від сповіщень про відключення електроенергії." {
			found = true
		}
	}
	assert.True(t, found)
}

func TestBot_StopWithNoSubscription(t *testing.T) {
	br, _, msgs := setupBot(t)
	br.HandleMessage(makeCmd(100, "stop"))

	found := false
	for _, m := range *msgs {
		if m.ChatID == 100 && m.Text == "Ви не маєте активної підписки." {
			found = true
		}
	}
	assert.True(t, found)
}

func TestBot_StopRemoveError(t *testing.T) {
	br, userRepo, msgs := setupBot(t)
	userRepo.rmErr = errors.New("disk error")

	br.HandleMessage(makeCmd(100, "stop"))

	found := false
	for _, m := range *msgs {
		if m.ChatID == 100 && m.Text == "Сталася помилка. Спробуйте пізніше." {
			found = true
		}
	}
	assert.True(t, found)
}

func TestBot_SubscriptionFindError(t *testing.T) {
	br, userRepo, msgs := setupBot(t)
	userRepo.findErr = errors.New("disk error")

	br.HandleMessage(makeCmd(100, "subscription"))

	found := false
	for _, m := range *msgs {
		if m.ChatID == 100 && m.Text == "Сталася помилка. Спробуйте пізніше." {
			found = true
		}
	}
	assert.True(t, found)
}

func TestBot_StartDuringConversation_Resets(t *testing.T) {
	br, _, _ := setupBot(t)
	br.HandleMessage(makeCmd(100, "start"))
	br.HandleMessage(makeMsg(100, "Наукова"))
	// Now in SaveSubscription step
	state := br.GetState(100)
	assert.Equal(t, StepSaveSubscription, state.Step)

	// /start resets
	br.HandleMessage(makeCmd(100, "start"))
	state = br.GetState(100)
	assert.Equal(t, StepSearchStreet, state.Step)
}

func TestBot_StopDuringConversation_ClearsState(t *testing.T) {
	br, _, _ := setupBot(t)
	br.HandleMessage(makeCmd(100, "start"))
	br.HandleMessage(makeCmd(100, "stop"))
	state := br.GetState(100)
	assert.Nil(t, state)
}

func TestBot_ExpiredSession(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	currentTime := now
	clock := func() time.Time { return currentTime }

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := tgbotapi.APIResponse{Ok: true, Result: json.RawMessage(`{"id":123,"is_bot":true,"first_name":"Test"}`)}
		if r.URL.Path != "/bottest-token/getMe" {
			resp.Result = json.RawMessage(`{"message_id":1,"chat":{"id":100},"text":""}`)
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	api, _ := tgbotapi.NewBotAPIWithAPIEndpoint("test-token", server.URL+"/bot%s/%s")
	userRepo := newTestUserRepo()
	streetRepo := &testStreetRepo{streets: []domain.Street{{ID: 1, Name: "Наукова"}}}
	cleanupCh := make(chan time.Time)

	br := NewBotRunner(BotRunnerConfig{
		Bot:                     api,
		SearchStreetService:     subscription.NewSearchStreetService(streetRepo),
		ShowSubscriptionService: subscription.NewShowSubscriptionService(userRepo),
		SaveSubscriptionService: subscription.NewSaveSubscriptionService(userRepo),
		UserRepo:                userRepo,
		Clock:                   clock,
		CleanupTicker:           cleanupCh,
		TTL:                     30 * time.Minute,
	})
	defer br.Close()

	br.HandleMessage(makeCmd(100, "start"))
	assert.NotNil(t, br.GetState(100))

	// Advance time past TTL
	currentTime = now.Add(31 * time.Minute)

	// Trigger cleanup
	cleanupCh <- currentTime

	// Give cleanup a moment to run
	time.Sleep(10 * time.Millisecond)

	assert.Nil(t, br.GetState(100))
}

func TestBot_SubscriptionShowsExistingUser(t *testing.T) {
	br, userRepo, msgs := setupBot(t)
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	userRepo.users[100] = &domain.User{ID: 100, Address: addr}

	br.HandleMessage(makeCmd(100, "subscription"))

	found := false
	for _, m := range *msgs {
		if m.ChatID == 100 {
			if m.Text == "Ваша поточна підписка:\nВулиця: Стрийська\nБудинок: 10" {
				found = true
			}
		}
	}
	assert.True(t, found)
}
