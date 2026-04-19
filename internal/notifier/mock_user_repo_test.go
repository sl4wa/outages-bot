package notifier

import "outages-bot/internal/users"

type mockUserRepo struct {
	users     map[int64]*users.User
	saved     []*users.User
	removed   []int64
	saveErr   error
	removeErr error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[int64]*users.User)}
}

func (m *mockUserRepo) FindAll() []*users.User {
	var result []*users.User
	for _, u := range m.users {
		result = append(result, u)
	}
	return result
}

func (m *mockUserRepo) Find(chatID int64) (*users.User, error) {
	u, ok := m.users[chatID]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepo) Save(user *users.User) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.saved = append(m.saved, user)
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Remove(chatID int64) (bool, error) {
	if m.removeErr != nil {
		return false, m.removeErr
	}
	m.removed = append(m.removed, chatID)
	if _, ok := m.users[chatID]; ok {
		delete(m.users, chatID)
		return true, nil
	}
	return false, nil
}
