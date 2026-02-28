package subscription

import "outages-bot/internal/domain"

// UnsubscribeService handles user unsubscription.
type UnsubscribeService struct {
	userRepo domain.UserRepository
}

// NewUnsubscribeService creates a new UnsubscribeService.
func NewUnsubscribeService(userRepo domain.UserRepository) *UnsubscribeService {
	return &UnsubscribeService{userRepo: userRepo}
}

// UnsubscribeResult holds the result of an unsubscribe operation.
type UnsubscribeResult struct {
	Message string
	Err     error
}

// Handle removes the user's subscription and returns an appropriate message.
func (s *UnsubscribeService) Handle(chatID int64) UnsubscribeResult {
	removed, err := s.userRepo.Remove(chatID)
	if err != nil {
		return UnsubscribeResult{
			Message: "Сталася помилка. Спробуйте пізніше.",
			Err:     err,
		}
	}

	if removed {
		return UnsubscribeResult{
			Message: "Ви успішно відписалися від сповіщень про відключення електроенергії.",
		}
	}

	return UnsubscribeResult{
		Message: "Ви не маєте активної підписки.",
	}
}
