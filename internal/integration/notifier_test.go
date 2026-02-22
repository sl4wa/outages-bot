package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"outages-bot/internal/application"
	"outages-bot/internal/application/notification"
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

type mockNotifSender struct {
	sent []application.NotificationSenderDTO
	errs map[int64]error
}

func (m *mockNotifSender) Send(dto application.NotificationSenderDTO) error {
	m.sent = append(m.sent, dto)
	if err, ok := m.errs[dto.UserID]; ok {
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
}

func (s *NotifierSuite) makeServer() {
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(s.apiBody))
	}))
	s.T().Cleanup(s.server.Close)
}

func (s *NotifierSuite) makeAPIBody(streetID int, buildings []string, comment string) {
	start := "2024-01-01T08:00:00+00:00"
	end := "2024-01-01T16:00:00+00:00"
	buildingsJSON, _ := json.Marshal(buildings)
	s.apiBody = `{"hydra:member":[{"id":1,"dateEvent":"` + start + `","datePlanIn":"` + end + `","koment":"` + comment + `","buildingNames":` + string(buildingsJSON) + `,"city":{"name":"Львів"},"street":{"id":` + json.Number(string(rune('0'+streetID))).String() + `,"name":"Стрийська"}}]}`
}

func (s *NotifierSuite) saveUser(chatID int64, streetID int, building string) {
	addr, err := domain.NewUserAddress(streetID, "Стрийська", building)
	require.NoError(s.T(), err)
	user := &domain.User{ID: chatID, Address: addr}
	require.NoError(s.T(), s.userRepo.Save(user))
}

func (s *NotifierSuite) runPipeline() {
	s.makeServer()
	clock := func() time.Time { return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC) }
	apiProvider := outageapi.NewApiOutageProvider(s.server.URL, clock, nil)
	fetchService := notification.NewOutageFetchService(apiProvider)
	notifService := notification.NewNotificationService(s.sender, s.userRepo, nil)

	outages, err := fetchService.Handle(context.Background())
	require.NoError(s.T(), err)
	_, err = notifService.Handle(outages)
	require.NoError(s.T(), err)
}

func (s *NotifierSuite) TestMatchingUser_NotificationSent() {
	s.makeAPIBody(1, []string{"10", "12"}, "test")
	s.saveUser(100, 1, "10")

	s.runPipeline()

	assert.Len(s.T(), s.sender.sent, 1)
	assert.Equal(s.T(), int64(100), s.sender.sent[0].UserID)

	// User file should be updated with outage info
	user, err := s.userRepo.Find(100)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), user.OutageInfo)
}

func (s *NotifierSuite) TestNonMatchingUser_NoNotification() {
	s.makeAPIBody(1, []string{"10"}, "test")
	s.saveUser(100, 2, "10") // different street

	s.runPipeline()

	assert.Empty(s.T(), s.sender.sent)

	// User file unchanged (no outage info)
	user, err := s.userRepo.Find(100)
	require.NoError(s.T(), err)
	assert.Nil(s.T(), user.OutageInfo)
}

func (s *NotifierSuite) TestBlockedUser_FileDeleted() {
	s.makeAPIBody(1, []string{"10"}, "test")
	s.saveUser(100, 1, "10")
	s.sender.errs[100] = &application.NotificationSendError{UserID: 100, Code: 403, Message: "Forbidden"}

	s.runPipeline()

	user, err := s.userRepo.Find(100)
	require.NoError(s.T(), err)
	assert.Nil(s.T(), user) // deleted
}

func (s *NotifierSuite) TestNonBlockedError_UserNotRemoved() {
	s.makeAPIBody(1, []string{"10"}, "test")
	s.saveUser(100, 1, "10")
	s.sender.errs[100] = &application.NotificationSendError{UserID: 100, Code: 500, Message: "Server Error"}

	s.runPipeline()

	user, err := s.userRepo.Find(100)
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Nil(s.T(), user.OutageInfo) // not saved
}

func (s *NotifierSuite) TestDedup_SecondRunSendsNothing() {
	s.makeAPIBody(1, []string{"10"}, "test")
	s.saveUser(100, 1, "10")

	s.runPipeline()
	assert.Len(s.T(), s.sender.sent, 1)

	// Second run
	s.sender.sent = nil
	s.runPipeline()
	assert.Empty(s.T(), s.sender.sent)
}

func (s *NotifierSuite) TestMultipleUsers_CorrectSubsetNotified() {
	s.makeAPIBody(1, []string{"10", "12"}, "test")
	s.saveUser(100, 1, "10") // matches
	s.saveUser(200, 1, "14") // no match
	s.saveUser(300, 1, "12") // matches

	s.runPipeline()

	assert.Len(s.T(), s.sender.sent, 2)
	sentIDs := make([]int64, len(s.sender.sent))
	for i, msg := range s.sender.sent {
		sentIDs[i] = msg.UserID
	}
	assert.Contains(s.T(), sentIDs, int64(100))
	assert.Contains(s.T(), sentIDs, int64(300))
}

func TestNotifierSuite(t *testing.T) {
	suite.Run(t, new(NotifierSuite))
}
