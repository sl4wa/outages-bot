package users

// Repository defines the interface for user data access.
type Repository interface {
	FindAll() []*User
	Find(chatID int64) (*User, error)
	Save(user *User) error
	Remove(chatID int64) (bool, error)
}
