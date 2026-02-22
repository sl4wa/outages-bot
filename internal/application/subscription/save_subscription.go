package subscription

import (
	"fmt"
	"outages-bot/internal/domain"
)

// SaveSubscriptionResult holds the result of saving a subscription.
type SaveSubscriptionResult struct {
	Message string
	Success bool
}

// SaveSubscriptionService handles saving user subscriptions.
type SaveSubscriptionService struct {
	userRepo domain.UserRepository
}

// NewSaveSubscriptionService creates a new SaveSubscriptionService.
func NewSaveSubscriptionService(userRepo domain.UserRepository) *SaveSubscriptionService {
	return &SaveSubscriptionService{userRepo: userRepo}
}

// Handle saves a user subscription. Returns validation errors as unsuccessful results.
func (s *SaveSubscriptionService) Handle(chatID int64, streetID int, streetName, building string) *SaveSubscriptionResult {
	addr, err := domain.NewUserAddress(streetID, streetName, building)
	if err != nil {
		return &SaveSubscriptionResult{Message: err.Error(), Success: false}
	}

	user := &domain.User{
		ID:      chatID,
		Address: addr,
	}

	if err := s.userRepo.Save(user); err != nil {
		return &SaveSubscriptionResult{Message: "Сталася помилка. Спробуйте пізніше.", Success: false}
	}

	return &SaveSubscriptionResult{
		Message: fmt.Sprintf(
			"Ви підписалися на сповіщення про відключення електроенергії для вулиці %s, будинок %s.",
			streetName, building,
		),
		Success: true,
	}
}
