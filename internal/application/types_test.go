package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotificationSendError_IsBlocked_Code403(t *testing.T) {
	e := &NotificationSendError{UserID: 1, Code: 403, Message: "Forbidden: bot was blocked by the user"}
	assert.True(t, e.IsBlocked())
}

func TestNotificationSendError_IsBlocked_Code200(t *testing.T) {
	e := &NotificationSendError{UserID: 1, Code: 200, Message: "OK"}
	assert.False(t, e.IsBlocked())
}

func TestNotificationSendError_IsBlocked_MessageForbidden(t *testing.T) {
	e := &NotificationSendError{UserID: 1, Code: 0, Message: "Forbidden: bot was blocked"}
	assert.True(t, e.IsBlocked())
}

func TestNotificationSendError_IsBlocked_MessageForbiddenMixedCase(t *testing.T) {
	e := &NotificationSendError{UserID: 1, Code: 0, Message: "FORBIDDEN by user"}
	assert.True(t, e.IsBlocked())
}

func TestNotificationSendError_IsBlocked_NoMatch(t *testing.T) {
	e := &NotificationSendError{UserID: 1, Code: 500, Message: "Internal server error"}
	assert.False(t, e.IsBlocked())
}
