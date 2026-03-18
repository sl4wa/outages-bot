package users

// UserInfoDTO is a data transfer object for Telegram user info.
type UserInfoDTO struct {
	ChatID    int64
	Username  string
	FirstName string
	LastName  string
}

// UserInfoProvider retrieves user info.
type UserInfoProvider interface {
	GetUserInfo(chatID int64) (UserInfoDTO, error)
}
