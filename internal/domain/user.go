package domain

// User represents a subscribed user.
type User struct {
	ID         int64
	Address    UserAddress
	OutageInfo *OutageInfo
}

// WithNotifiedOutage returns a new User with the outage info set from the given outage.
func (u *User) WithNotifiedOutage(outage *Outage) *User {
	info := NewOutageInfo(outage.Period, outage.Description)
	return &User{
		ID:         u.ID,
		Address:    u.Address,
		OutageInfo: &info,
	}
}

// IsAlreadyNotifiedAbout checks if the user has already been notified about the given outage info.
func (u *User) IsAlreadyNotifiedAbout(info OutageInfo) bool {
	if u.OutageInfo == nil {
		return false
	}
	return u.OutageInfo.Equals(info)
}

// UserRepository defines the interface for user data access.
type UserRepository interface {
	FindAll() ([]*User, error)
	Find(chatID int64) (*User, error)
	Save(user *User) error
	Remove(chatID int64) (bool, error)
}
