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

// FindOutageForNotification finds the first matching outage for a user that they haven't been notified about.
func (u *User) FindOutageForNotification(allOutages []*Outage) *Outage {
	for _, outage := range allOutages {
		if !outage.Address.CoversUserAddress(u.Address) {
			continue
		}

		outageInfo := NewOutageInfo(outage.Period, outage.Description)

		if u.OutageInfo != nil && u.OutageInfo.Equals(outageInfo) {
			return nil
		}

		return outage
	}

	return nil
}

// UserRepository defines the interface for user data access.
type UserRepository interface {
	FindAll() []*User
	Find(chatID int64) (*User, error)
	Save(user *User) error
	Remove(chatID int64) (bool, error)
}
