package bot

import "outages-bot/internal/domain"

// Unsubscribe handles user unsubscription.
type Unsubscribe struct {
	userRepo domain.UserRepository
}

// NewUnsubscribe creates a new Unsubscribe.
func NewUnsubscribe(userRepo domain.UserRepository) *Unsubscribe {
	return &Unsubscribe{userRepo: userRepo}
}

// UnsubscribeResult holds the result of an unsubscribe operation.
type UnsubscribeResult struct {
	Message string
	Err     error
}

// Handle removes the user's subscription and returns an appropriate message.
func (s *Unsubscribe) Handle(chatID int64) UnsubscribeResult {
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
