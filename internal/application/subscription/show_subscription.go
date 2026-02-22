package subscription

import (
	"fmt"
	"outages-bot/internal/domain"
)

// ShowSubscriptionService handles showing the current subscription or prompting for a new one.
type ShowSubscriptionService struct {
	userRepo domain.UserRepository
}

// NewShowSubscriptionService creates a new ShowSubscriptionService.
func NewShowSubscriptionService(userRepo domain.UserRepository) *ShowSubscriptionService {
	return &ShowSubscriptionService{userRepo: userRepo}
}

// ShowCurrent returns the current subscription status without the update prompt.
// Unlike Handle, it returns errors to the caller instead of swallowing them.
func (s *ShowSubscriptionService) ShowCurrent(chatID int64) (string, error) {
	user, err := s.userRepo.Find(chatID)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "Ви не маєте активної підписки.", nil
	}

	return fmt.Sprintf(
		"Ваша поточна підписка:\nВулиця: %s\nБудинок: %s",
		user.Address.StreetName,
		user.Address.Building,
	), nil
}

// Handle returns a message showing the current subscription or prompting for a new one.
// If the repository returns an error (e.g., corrupted data), it falls back to the new-user prompt.
func (s *ShowSubscriptionService) Handle(chatID int64) string {
	user, err := s.userRepo.Find(chatID)
	if err != nil {
		// Matches PHP behavior: catch Throwable, treat as null
		return "Будь ласка, введіть назву вулиці:"
	}

	if user == nil {
		return "Будь ласка, введіть назву вулиці:"
	}

	return fmt.Sprintf(
		"Ваша поточна підписка:\nВулиця: %s\nБудинок: %s\n\nБудь ласка, введіть нову назву вулиці для оновлення підписки:",
		user.Address.StreetName,
		user.Address.Building,
	)
}
