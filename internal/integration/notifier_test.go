package integration

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"outages-bot/internal/application/notifier"
	"outages-bot/internal/application/service"
	"outages-bot/internal/client/outageapi"
	"outages-bot/internal/domain"
	"outages-bot/internal/repository"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type sentNotification struct {
	UserID  int64
	Content notifier.NotificationContent
}

type mockNotifSender struct {
	sent []sentNotification
	errs map[int64]error
}

func (m *mockNotifSender) Send(userID int64, content notifier.NotificationContent) error {
	m.sent = append(m.sent, sentNotification{UserID: userID, Content: content})
	if err, ok := m.errs[userID]; ok {
		return err
	}
	return nil
}

type NotifierSuite struct {
	suite.Suite
	server   *httptest.Server
	userRepo *repository.FileUserRepository
	dataDir  string
	apiBody  string
	sender   *mockNotifSender
}

func (s *NotifierSuite) SetupTest() {
	s.dataDir = s.T().TempDir()
	var err error
	s.userRepo, err = repository.NewFileUserRepository(filepath.Join(s.dataDir, "users"))
	require.NoError(s.T(), err)
	s.sender = &mockNotifSender{errs: make(map[int64]error)}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(s.apiBody))
	}))
	s.T().Cleanup(s.server.Close)
}

func (s *NotifierSuite) loadFixture() {
	data, err := os.ReadFile("testdata/loe_data.json")
	require.NoError(s.T(), err)
	s.apiBody = string(data)
}

func (s *NotifierSuite) filterFixtureByComment(comment string) string {
	var resp struct {
		Context string            `json:"@context"`
		ID      string            `json:"@id"`
		Type    string            `json:"@type"`
		Members []json.RawMessage `json:"hydra:member"`
	}
	err := json.Unmarshal([]byte(s.apiBody), &resp)
	require.NoError(s.T(), err)

	var filtered []json.RawMessage
	for _, raw := range resp.Members {
		var row struct {
			Koment string `json:"koment"`
		}
		if err := json.Unmarshal(raw, &row); err != nil {
			continue
		}
		if row.Koment == comment {
			filtered = append(filtered, raw)
		}
	}

	result, err := json.Marshal(struct {
		Context string            `json:"@context"`
		ID      string            `json:"@id"`
		Type    string            `json:"@type"`
		Members []json.RawMessage `json:"hydra:member"`
	}{
		Context: resp.Context,
		ID:      resp.ID,
		Type:    resp.Type,
		Members: filtered,
	})
	require.NoError(s.T(), err)
	return string(result)
}

func (s *NotifierSuite) saveUser(chatID int64, streetID int, streetName, building string) {
	addr, err := domain.NewUserAddress(streetID, streetName, building)
	require.NoError(s.T(), err)
	user := &domain.User{ID: chatID, Address: addr}
	require.NoError(s.T(), s.userRepo.Save(user))
}

func (s *NotifierSuite) runPipeline() {
	clock := func() time.Time { return time.Date(2024, 11, 28, 9, 0, 0, 0, time.UTC) }
	apiProvider := outageapi.NewProvider(s.server.URL, clock, nil)
	fetchService := service.NewFetchOutages(apiProvider)
	notifyUsers := notifier.NewNotifyUsers(fetchService, s.sender, s.userRepo, nil)

	err := notifyUsers.Handle(context.Background())
	require.NoError(s.T(), err)
}

func (s *NotifierSuite) TestMatchingUser_NotificationSent() {
	s.loadFixture()
	s.saveUser(100, 12445, "Стрийська", "45")

	s.runPipeline()

	assert.Len(s.T(), s.sender.sent, 1)
	assert.Equal(s.T(), int64(100), s.sender.sent[0].UserID)

	user, err := s.userRepo.Find(100)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), user.OutageInfo)
}

func (s *NotifierSuite) TestNonMatchingUser_NoNotification() {
	s.loadFixture()
	s.saveUser(100, 12445, "Стрийська", "999")

	s.runPipeline()

	assert.Empty(s.T(), s.sender.sent)

	user, err := s.userRepo.Find(100)
	require.NoError(s.T(), err)
	assert.Nil(s.T(), user.OutageInfo)
}

func (s *NotifierSuite) TestBlockedUser_FileDeleted() {
	s.loadFixture()
	s.saveUser(100, 12445, "Стрийська", "45")
	s.sender.errs[100] = notifier.ErrRecipientUnavailable

	s.runPipeline()

	user, err := s.userRepo.Find(100)
	require.NoError(s.T(), err)
	assert.Nil(s.T(), user)
}

func (s *NotifierSuite) TestNonBlockedError_UserNotRemoved() {
	s.loadFixture()
	s.saveUser(100, 12445, "Стрийська", "45")
	s.sender.errs[100] = errors.New("server error")

	s.runPipeline()

	user, err := s.userRepo.Find(100)
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Nil(s.T(), user.OutageInfo)
}

func (s *NotifierSuite) TestDedup_SecondRunSendsNothing() {
	s.loadFixture()
	s.saveUser(100, 12445, "Стрийська", "45")

	s.runPipeline()
	assert.Len(s.T(), s.sender.sent, 1)

	s.sender.sent = nil
	s.runPipeline()
	assert.Empty(s.T(), s.sender.sent)
}

func (s *NotifierSuite) TestMultipleUsers_CorrectSubsetNotified() {
	s.loadFixture()
	s.saveUser(100, 12445, "Стрийська", "45")  // matches
	s.saveUser(200, 12445, "Стрийська", "14")  // no match
	s.saveUser(300, 12445, "Стрийська", "108") // matches

	s.runPipeline()

	assert.Len(s.T(), s.sender.sent, 2)
	sentIDs := make([]int64, len(s.sender.sent))
	for i, msg := range s.sender.sent {
		sentIDs[i] = msg.UserID
	}
	assert.Contains(s.T(), sentIDs, int64(100))
	assert.Contains(s.T(), sentIDs, int64(300))
}

func (s *NotifierSuite) TestFirstMatchingOutageSelected() {
	s.loadFixture()
	s.saveUser(100, 12445, "Стрийська", "45")

	s.runPipeline()

	require.Len(s.T(), s.sender.sent, 1)
	content := s.sender.sent[0].Content
	// The first entry covering building 45 in fixture iteration order is the
	// standalone "45" row: 2024-11-28T06:33:00+00:00 → 2024-11-28T10:57:00+00:00
	assert.Equal(s.T(), time.Date(2024, 11, 28, 6, 33, 0, 0, time.UTC).Unix(), content.Start.Unix())
	assert.Equal(s.T(), time.Date(2024, 11, 28, 10, 57, 0, 0, time.UTC).Unix(), content.End.Unix())
	assert.Equal(s.T(), "Застосування ГАВ", content.Comment)
}

func (s *NotifierSuite) TestNewOutageAfterPriorNotification() {
	s.loadFixture()
	s.apiBody = s.filterFixtureByComment("Застосування ГАВ")
	s.saveUser(100, 12445, "Стрийська", "108")

	s.runPipeline()

	require.Len(s.T(), s.sender.sent, 1)
	content := s.sender.sent[0].Content
	assert.Equal(s.T(), time.Date(2024, 11, 28, 6, 47, 0, 0, time.UTC).Unix(), content.Start.Unix())
	assert.Equal(s.T(), time.Date(2024, 11, 28, 10, 0, 0, 0, time.UTC).Unix(), content.End.Unix())
	assert.Equal(s.T(), "Застосування ГАВ", content.Comment)

	// Second run with ГПВ entries only
	s.sender.sent = nil
	s.loadFixture()
	s.apiBody = s.filterFixtureByComment("Застосування ГПВ")
	s.runPipeline()

	require.Len(s.T(), s.sender.sent, 1)
	content = s.sender.sent[0].Content
	assert.Equal(s.T(), time.Date(2024, 11, 28, 8, 0, 0, 0, time.UTC).Unix(), content.Start.Unix())
	assert.Equal(s.T(), time.Date(2024, 11, 28, 10, 0, 0, 0, time.UTC).Unix(), content.End.Unix())
	assert.Equal(s.T(), "Застосування ГПВ", content.Comment)
}

func TestNotifierSuite(t *testing.T) {
	suite.Run(t, new(NotifierSuite))
}
