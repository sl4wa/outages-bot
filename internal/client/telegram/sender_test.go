package telegram

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"outages-bot/internal/application"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTelegramServer(t *testing.T, sendHandler func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *tgbotapi.BotAPI) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bottest-token/getMe" {
			resp := tgbotapi.APIResponse{Ok: true, Result: json.RawMessage(`{"id":123,"is_bot":true,"first_name":"Test"}`)}
			json.NewEncoder(w).Encode(resp)
			return
		}
		if r.URL.Path == "/bottest-token/sendMessage" {
			sendHandler(w, r)
			return
		}
		w.WriteHeader(404)
	}))
	t.Cleanup(server.Close)

	api, err := tgbotapi.NewBotAPIWithAPIEndpoint("test-token", server.URL+"/bot%s/%s")
	require.NoError(t, err)
	return server, api
}

func testDTO() application.NotificationSenderDTO {
	return application.NotificationSenderDTO{
		UserID: 100,
		Text:   "test notification",
	}
}

func TestSender_SuccessfulSend(t *testing.T) {
	_, api := makeTelegramServer(t, func(w http.ResponseWriter, r *http.Request) {
		resp := tgbotapi.APIResponse{Ok: true, Result: json.RawMessage(`{"message_id":1,"chat":{"id":100},"text":"test"}`)}
		json.NewEncoder(w).Encode(resp)
	})

	sender := NewNotificationSender(api)
	err := sender.Send(testDTO())
	assert.NoError(t, err)
}

func TestSender_Forbidden403(t *testing.T) {
	_, api := makeTelegramServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		resp := tgbotapi.APIResponse{Ok: false, ErrorCode: 403, Description: "Forbidden: bot was blocked by the user"}
		json.NewEncoder(w).Encode(resp)
	})

	sender := NewNotificationSender(api)
	err := sender.Send(testDTO())
	require.Error(t, err)

	sendErr, ok := err.(*application.NotificationSendError)
	require.True(t, ok)
	assert.Equal(t, 403, sendErr.Code)
	assert.True(t, sendErr.IsBlocked())
}

func TestSender_BadRequest400(t *testing.T) {
	_, api := makeTelegramServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		resp := tgbotapi.APIResponse{Ok: false, ErrorCode: 400, Description: "Bad Request: message is too long"}
		json.NewEncoder(w).Encode(resp)
	})

	sender := NewNotificationSender(api)
	err := sender.Send(testDTO())
	require.Error(t, err)

	sendErr, ok := err.(*application.NotificationSendError)
	require.True(t, ok)
	assert.Equal(t, 400, sendErr.Code)
	assert.False(t, sendErr.IsBlocked())
}

func TestSender_TooManyRequests429(t *testing.T) {
	_, api := makeTelegramServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		resp := tgbotapi.APIResponse{Ok: false, ErrorCode: 429, Description: "Too Many Requests: retry after 30"}
		json.NewEncoder(w).Encode(resp)
	})

	sender := NewNotificationSender(api)
	err := sender.Send(testDTO())
	require.Error(t, err)

	sendErr, ok := err.(*application.NotificationSendError)
	require.True(t, ok)
	assert.Equal(t, 429, sendErr.Code)
	assert.False(t, sendErr.IsBlocked())
}

func TestSender_NetworkError_Code0(t *testing.T) {
	// Create a server that's immediately closed to simulate network error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bottest-token/getMe" {
			resp := tgbotapi.APIResponse{Ok: true, Result: json.RawMessage(`{"id":123,"is_bot":true,"first_name":"Test"}`)}
			json.NewEncoder(w).Encode(resp)
			return
		}
	}))
	api, err := tgbotapi.NewBotAPIWithAPIEndpoint("test-token", server.URL+"/bot%s/%s")
	require.NoError(t, err)
	server.Close() // Close to cause network error

	sender := NewNotificationSender(api)
	err = sender.Send(testDTO())
	require.Error(t, err)

	sendErr, ok := err.(*application.NotificationSendError)
	require.True(t, ok)
	assert.Equal(t, 0, sendErr.Code) // fallback code
}

func TestSender_MalformedJSON(t *testing.T) {
	_, api := makeTelegramServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("not json at all"))
	})

	sender := NewNotificationSender(api)
	err := sender.Send(testDTO())
	require.Error(t, err)

	sendErr, ok := err.(*application.NotificationSendError)
	require.True(t, ok)
	assert.Equal(t, 0, sendErr.Code) // non-API error falls back to 0
}

func TestSender_HTMLParseMode(t *testing.T) {
	var capturedParseMode string
	_, api := makeTelegramServer(t, func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		capturedParseMode = r.FormValue("parse_mode")
		resp := tgbotapi.APIResponse{Ok: true, Result: json.RawMessage(`{"message_id":1,"chat":{"id":100},"text":"test"}`)}
		json.NewEncoder(w).Encode(resp)
	})

	sender := NewNotificationSender(api)
	err := sender.Send(testDTO())
	require.NoError(t, err)
	assert.Equal(t, "HTML", capturedParseMode)
}

func TestSender_MessageText(t *testing.T) {
	var capturedText string
	_, api := makeTelegramServer(t, func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		capturedText = r.FormValue("text")
		resp := tgbotapi.APIResponse{Ok: true, Result: json.RawMessage(`{"message_id":1,"chat":{"id":100},"text":"test"}`)}
		json.NewEncoder(w).Encode(resp)
	})

	sender := NewNotificationSender(api)
	dto := testDTO()
	err := sender.Send(dto)
	require.NoError(t, err)

	assert.Equal(t, dto.Text, capturedText)
}

func TestSender_ForbiddenInMessage_IsBlocked(t *testing.T) {
	_, api := makeTelegramServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		resp := tgbotapi.APIResponse{Ok: false, ErrorCode: 403, Description: "Forbidden: user is deactivated"}
		json.NewEncoder(w).Encode(resp)
	})

	sender := NewNotificationSender(api)
	err := sender.Send(testDTO())
	require.Error(t, err)

	sendErr, ok := err.(*application.NotificationSendError)
	require.True(t, ok)
	assert.True(t, sendErr.IsBlocked())
}
