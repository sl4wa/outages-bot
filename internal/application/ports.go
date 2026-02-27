package application

import "context"

// OutageProvider fetches outage data from an external source.
type OutageProvider interface {
	FetchOutages(ctx context.Context) ([]OutageDTO, error)
}

// NotificationSender sends notifications to users.
type NotificationSender interface {
	Send(dto NotificationSenderDTO) error
}

// UserInfoProvider retrieves user info.
type UserInfoProvider interface {
	GetUserInfo(chatID int64) (UserInfoDTO, error)
}
