package application

import (
	"strings"
	"time"
)

// OutageDTO is a data transfer object for outages from the provider.
type OutageDTO struct {
	ID         int
	Start      time.Time
	End        time.Time
	City       string
	StreetID   int
	StreetName string
	Buildings  []string
	Comment    string
}

// NotificationSenderDTO is a data transfer object for sending notifications.
type NotificationSenderDTO struct {
	UserID     int64
	City       string
	StreetName string
	Buildings  []string
	Start      time.Time
	End        time.Time
	Comment    string
}

// NotificationSendError is a custom error type for notification send failures.
type NotificationSendError struct {
	UserID  int64
	Code    int
	Message string
}

func (e *NotificationSendError) Error() string {
	return e.Message
}

// IsBlocked returns true if the error indicates the user has blocked the bot.
func (e *NotificationSendError) IsBlocked() bool {
	return e.Code == 403 || strings.Contains(strings.ToLower(e.Message), "forbidden")
}

// UserInfoDTO is a data transfer object for Telegram user info.
type UserInfoDTO struct {
	ChatID    int64
	Username  string
	FirstName string
	LastName  string
}
