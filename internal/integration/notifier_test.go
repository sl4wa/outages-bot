package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"outages-bot/internal/application"
	"strings"
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
	apiProvider := outageapi.NewProvider(s.server.URL, clock, nil)
	fetchService := notification.NewOutageFetchService(apiProvider)
	notifService := notification.NewService(s.sender, s.userRepo, nil)

	outages, err := fetchService.Handle(context.Background())
	require.NoError(s.T(), err)
	notifService.Handle(outages)
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

func (s *NotifierSuite) makeAPIBodyMultiple(outages []struct {
	ID        int
	StreetID  int
	Buildings []string
	Comment   string
	Start     string
	End       string
}) {
	members := make([]string, 0, len(outages))
	for _, o := range outages {
		buildingsJSON, _ := json.Marshal(o.Buildings)
		member := fmt.Sprintf(
			`{"id":%d,"dateEvent":"%s","datePlanIn":"%s","koment":"%s","buildingNames":%s,"city":{"name":"Львів"},"street":{"id":%d,"name":"Стрийська"}}`,
			o.ID, o.Start, o.End, o.Comment, string(buildingsJSON), o.StreetID,
		)
		members = append(members, member)
	}
	s.apiBody = `{"hydra:member":[` + strings.Join(members, ",") + `]}`
}

func (s *NotifierSuite) TestMultipleOutages_UserMatchesFirst() {
	s.makeAPIBodyMultiple([]struct {
		ID        int
		StreetID  int
		Buildings []string
		Comment   string
		Start     string
		End       string
	}{
		{ID: 1, StreetID: 1, Buildings: []string{"10"}, Comment: "first", Start: "2024-01-01T08:00:00+00:00", End: "2024-01-01T16:00:00+00:00"},
		{ID: 2, StreetID: 2, Buildings: []string{"20"}, Comment: "second", Start: "2024-01-01T08:00:00+00:00", End: "2024-01-01T16:00:00+00:00"},
	})
	s.saveUser(100, 1, "10") // matches outage 1 only

	s.runPipeline()

	assert.Len(s.T(), s.sender.sent, 1)
	assert.Equal(s.T(), int64(100), s.sender.sent[0].UserID)
	assert.Equal(s.T(), "first", s.sender.sent[0].Comment)
}

func (s *NotifierSuite) TestDedup_DifferentOutage_SecondRunSendsNew() {
	s.makeAPIBody(1, []string{"10"}, "first outage")
	s.saveUser(100, 1, "10")

	// First run: sends notification for "first outage"
	s.runPipeline()
	assert.Len(s.T(), s.sender.sent, 1)

	// Second run with a different outage (different comment = different description)
	s.sender.sent = nil
	s.makeAPIBodyMultiple([]struct {
		ID        int
		StreetID  int
		Buildings []string
		Comment   string
		Start     string
		End       string
	}{
		{ID: 2, StreetID: 1, Buildings: []string{"10"}, Comment: "new outage", Start: "2024-02-01T08:00:00+00:00", End: "2024-02-01T16:00:00+00:00"},
	})
	s.runPipeline()
	assert.Len(s.T(), s.sender.sent, 1)
	assert.Equal(s.T(), "new outage", s.sender.sent[0].Comment)
}

func TestNotifierSuite(t *testing.T) {
	suite.Run(t, new(NotifierSuite))
}
