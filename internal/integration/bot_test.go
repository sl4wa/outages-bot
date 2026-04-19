package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"outages-bot/internal/persistence"
	"outages-bot/internal/telegram"
	"outages-bot/internal/users"
	"path/filepath"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BotSuite struct {
	suite.Suite
	runner     *telegram.BotRunner
	userRepo   *persistence.FileUserRepository
	streetRepo *persistence.FileStreetRepository
	dataDir    string
}

func (s *BotSuite) SetupTest() {
	s.dataDir = s.T().TempDir()

	var err error
	s.userRepo, err = persistence.NewFileUserRepository(filepath.Join(s.dataDir, "users"))
	require.NoError(s.T(), err)

	s.streetRepo, err = persistence.NewFileStreetRepository("testdata/streets.csv")
	require.NoError(s.T(), err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := tgbotapi.APIResponse{Ok: true}
		if r.URL.Path == "/bottest-token/getMe" {
			resp.Result = json.RawMessage(`{"id":123,"is_bot":true,"first_name":"Test"}`)
		} else {
			resp.Result = json.RawMessage(`{"message_id":1,"chat":{"id":100},"text":""}`)
		}
		json.NewEncoder(w).Encode(resp)
	}))
	s.T().Cleanup(server.Close)

	api, err := tgbotapi.NewBotAPIWithAPIEndpoint("test-token", server.URL+"/bot%s/%s")
	require.NoError(s.T(), err)

	cleanupCh := make(chan time.Time)
	s.runner = telegram.NewBotRunner(telegram.BotRunnerConfig{
		Bot:              api,
		SearchStreet:     users.NewSearchStreet(s.streetRepo),
		ShowSubscription: users.NewShowSubscription(s.userRepo),
		SaveSubscription: users.NewSaveSubscription(s.userRepo),
		Unsubscribe:      users.NewUnsubscribe(s.userRepo),
		CleanupTicker:    cleanupCh,
		TTL:              30 * time.Minute,
	})
	s.T().Cleanup(s.runner.Close)
}

func makeIntMsg(chatID int64, text string) *tgbotapi.Message {
	return &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: chatID},
		Text: text,
	}
}

func makeIntCmd(chatID int64, command string) *tgbotapi.Message {
	return &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: chatID},
		Text: "/" + command,
		Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len("/" + command)},
		},
	}
}

func (s *BotSuite) TestSearchSaveVerifyFile() {
	// Start conversation
	s.runner.HandleMessage(makeIntCmd(100, "start"))

	// Search for a street from streets.csv fixture
	s.runner.HandleMessage(makeIntMsg(100, "Молдавська"))

	// Save subscription with building
	s.runner.HandleMessage(makeIntMsg(100, "25"))

	// Verify user file on disk
	user, err := s.userRepo.Find(100)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), user)
	assert.Equal(s.T(), int64(100), user.ID)
	assert.Equal(s.T(), "Молдавська", user.Address.StreetName)
	assert.Equal(s.T(), "25", user.Address.Building)
	assert.Equal(s.T(), 12444, user.Address.StreetID)
}

func (s *BotSuite) TestShowSubscription_ExistingUser() {
	addr, _ := users.NewAddress(12444, "Молдавська", "10")
	s.userRepo.Save(&users.User{ID: 100, Address: addr})

	s.runner.HandleMessage(makeIntCmd(100, "start"))

	// Should be in search step (got shown existing subscription)
	state := s.runner.GetState(100)
	require.NotNil(s.T(), state)
	assert.Equal(s.T(), telegram.StepSearchStreet, state.Step)
}

func (s *BotSuite) TestShowSubscription_NewUser() {
	s.runner.HandleMessage(makeIntCmd(100, "start"))
	state := s.runner.GetState(100)
	require.NotNil(s.T(), state)
	assert.Equal(s.T(), telegram.StepSearchStreet, state.Step)
}

func (s *BotSuite) TestRemoveUser_FileDeleted() {
	addr, _ := users.NewAddress(12444, "Молдавська", "10")
	s.userRepo.Save(&users.User{ID: 100, Address: addr})

	s.runner.HandleMessage(makeIntCmd(100, "stop"))

	user, err := s.userRepo.Find(100)
	require.NoError(s.T(), err)
	assert.Nil(s.T(), user)
}

func TestBotSuite(t *testing.T) {
	suite.Run(t, new(BotSuite))
}
