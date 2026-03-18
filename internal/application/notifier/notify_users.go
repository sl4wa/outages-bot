package notifier

import (
	"context"
	"errors"
	"fmt"
	"log"
	"outages-bot/internal/application/service"
	"outages-bot/internal/domain"
	"strings"
)

// NotificationSender sends notifications to users.
type NotificationSender interface {
	Send(dto NotificationSenderDTO) error
}

// NotificationSenderDTO is a data transfer object for sending notifications.
type NotificationSenderDTO struct {
	UserID int64
	Text   string
}

// NotificationSendError is a custom error type for notification send failures.
type NotificationSendError struct {
	UserID  int64
	Code    int
	Message string
}

func (e *NotificationSendError) Error() string {
	return e.Message
}

// IsBlocked returns true if the error indicates the user has blocked the bot.
func (e *NotificationSendError) IsBlocked() bool {
	return e.Code == 403 || strings.Contains(strings.ToLower(e.Message), "forbidden")
}

// NotifyUsers fetches outages and sends notifications to all affected users.
type NotifyUsers struct {
	fetchService *service.FetchOutages
	sender       NotificationSender
	userRepo     domain.UserRepository
	logger       *log.Logger
}

// NewNotifyUsers creates a new NotifyUsers.
func NewNotifyUsers(fetchService *service.FetchOutages, sender NotificationSender, userRepo domain.UserRepository, logger *log.Logger) *NotifyUsers {
	if logger == nil {
		logger = log.Default()
	}
	return &NotifyUsers{
		fetchService: fetchService,
		sender:       sender,
		userRepo:     userRepo,
		logger:       logger,
	}
}

// Handle fetches outages and sends notifications to all affected users.
func (n *NotifyUsers) Handle(ctx context.Context) error {
	outages, err := n.fetchService.Handle(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch outages: %w", err)
	}

	users := n.userRepo.FindAll()

	for _, user := range users {
		outage := user.FindOutageForNotification(outages)
		if outage == nil {
			continue
		}

		dto := NotificationSenderDTO{
			UserID: user.ID,
			Text:   formatOutageNotification(outage),
		}

		if err := n.sender.Send(dto); err != nil {
			var sendErr *NotificationSendError
			if errors.As(err, &sendErr) && sendErr.IsBlocked() {
				if _, rmErr := n.userRepo.Remove(sendErr.UserID); rmErr != nil {
					n.logger.Printf("failed to remove blocked user %d: %v", sendErr.UserID, rmErr)
				}
			}
			// Non-blocking errors: user NOT removed, NOT saved, continue
			continue
		}

		updatedUser := user.WithNotifiedOutage(outage)
		if err := n.userRepo.Save(updatedUser); err != nil {
			n.logger.Printf("failed to save user %d: %v", user.ID, err)
		}
	}

	return nil
}
