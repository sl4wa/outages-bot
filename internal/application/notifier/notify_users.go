package notifier

import (
	"context"
	"errors"
	"fmt"
	"log"
	"outages-bot/internal/application/service"
	"outages-bot/internal/domain"
	"time"
)

// ErrRecipientUnavailable indicates the recipient can no longer receive messages
// (e.g. blocked the bot, deactivated account).
var ErrRecipientUnavailable = errors.New("recipient unavailable")

// NotificationSender sends notifications to users.
type NotificationSender interface {
	Send(userID int64, content NotificationContent) error
}

// NotificationContent carries the structured data needed to render a notification.
type NotificationContent struct {
	City       string
	StreetName string
	Buildings  []string
	Start      time.Time
	End        time.Time
	Comment    string
}

// OutageRepository is a dedup store for outage data fetched from the API.
type OutageRepository interface {
	Load() ([]*domain.Outage, error)
	Save(outages []*domain.Outage) error
}

// NotifyUsers fetches outages and sends notifications to all affected users.
type NotifyUsers struct {
	fetchService *service.FetchOutages
	sender       NotificationSender
	userRepo     domain.UserRepository
	outageRepo   OutageRepository
	logger       *log.Logger
}

// NewNotifyUsers creates a new NotifyUsers.
func NewNotifyUsers(fetchService *service.FetchOutages, sender NotificationSender, userRepo domain.UserRepository, outageRepo OutageRepository, logger *log.Logger) *NotifyUsers {
	if logger == nil {
		logger = log.Default()
	}
	return &NotifyUsers{
		fetchService: fetchService,
		sender:       sender,
		userRepo:     userRepo,
		outageRepo:   outageRepo,
		logger:       logger,
	}
}

// Handle fetches outages and sends notifications to all affected users.
func (n *NotifyUsers) Handle(ctx context.Context) error {
	outages, err := n.fetchService.Handle(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch outages: %w", err)
	}

	prev, err := n.outageRepo.Load()
	if err != nil {
		return fmt.Errorf("failed to load outage data: %w", err)
	}

	if prev != nil && domain.OutagesEqual(prev, outages) {
		n.logger.Printf("Outage data unchanged; checker/notifier logic skipped.")
		return nil
	}

	if prev == nil {
		n.logger.Printf("No prior outage data found; saving and continuing.")
	} else {
		n.logger.Printf("Outage data changed; saving and continuing.")
	}

	if err := n.outageRepo.Save(outages); err != nil {
		return fmt.Errorf("failed to save outage data: %w", err)
	}

	users := n.userRepo.FindAll()

	for _, user := range users {
		outage := user.FindOutageForNotification(outages)
		if outage == nil {
			continue
		}

		content := NotificationContent{
			City:       outage.Address.City,
			StreetName: outage.Address.StreetName,
			Buildings:  outage.Address.Buildings,
			Start:      outage.Period.StartDate,
			End:        outage.Period.EndDate,
			Comment:    outage.Description.Value,
		}

		if err := n.sender.Send(user.ID, content); err != nil {
			if errors.Is(err, ErrRecipientUnavailable) {
				if _, rmErr := n.userRepo.Remove(user.ID); rmErr != nil {
					n.logger.Printf("failed to remove blocked user %d: %v", user.ID, rmErr)
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
