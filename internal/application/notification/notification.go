package notification

import (
	"errors"
	"log"
	"outages-bot/internal/application"
	"outages-bot/internal/domain"
)

// NotificationService handles sending outage notifications to affected users.
type NotificationService struct {
	sender   application.NotificationSender
	userRepo domain.UserRepository
	logger   *log.Logger
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(sender application.NotificationSender, userRepo domain.UserRepository, logger *log.Logger) *NotificationService {
	if logger == nil {
		logger = log.Default()
	}
	return &NotificationService{
		sender:   sender,
		userRepo: userRepo,
		logger:   logger,
	}
}

// Handle sends notifications for the given outages to all affected users.
func (s *NotificationService) Handle(outages []*domain.Outage) {
	users := s.userRepo.FindAll()

	for _, user := range users {
		outage := domain.FindOutageForNotification(user, outages)
		if outage == nil {
			continue
		}

		dto := application.NotificationSenderDTO{
			UserID:     user.ID,
			City:       outage.Address.City,
			StreetName: outage.Address.StreetName,
			Buildings:  outage.Address.Buildings,
			Start:      outage.Period.StartDate,
			End:        outage.Period.EndDate,
			Comment:    outage.Description.Value,
		}

		if err := s.sender.Send(dto); err != nil {
			var sendErr *application.NotificationSendError
			if errors.As(err, &sendErr) && sendErr.IsBlocked() {
				if _, rmErr := s.userRepo.Remove(sendErr.UserID); rmErr != nil {
					s.logger.Printf("failed to remove blocked user %d: %v", sendErr.UserID, rmErr)
				}
			}
			// Non-blocking errors: user NOT removed, NOT saved, continue
			continue
		}

		updatedUser := user.WithNotifiedOutage(outage)
		if err := s.userRepo.Save(updatedUser); err != nil {
			s.logger.Printf("failed to save user %d: %v", user.ID, err)
		}
	}
}
