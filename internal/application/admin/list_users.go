package admin

import (
	"outages-bot/internal/domain"
	"sort"
)

// ListUsers returns all users sorted by outage start date (descending), users without outage at the end.
func ListUsers(userRepo domain.UserRepository) ([]*domain.User, error) {
	users, err := userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	sort.Slice(users, func(i, j int) bool {
		a, b := users[i], users[j]
		if a.OutageInfo == nil && b.OutageInfo == nil {
			return false
		}
		if a.OutageInfo == nil {
			return false
		}
		if b.OutageInfo == nil {
			return true
		}
		return b.OutageInfo.Period.StartDate.Before(a.OutageInfo.Period.StartDate)
	})

	return users, nil
}
