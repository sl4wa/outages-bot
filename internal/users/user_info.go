package users

// Info is a data transfer object for Telegram user info.
type Info struct {
	ChatID    int64
	Username  string
	FirstName string
	LastName  string
}

// InfoProvider retrieves user info.
type InfoProvider interface {
	GetUserInfo(chatID int64) (Info, error)
}
