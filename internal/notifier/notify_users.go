package notifier

import (
	"context"
	"errors"
	"fmt"
	"log"
	"outages-bot/internal/outage"
	"outages-bot/internal/users"
)

// NotifyUsers fetches outages and sends notifications to all affected users.
type NotifyUsers struct {
	fetchService *outage.FetchOutages
	sender       Sender
	userRepo     users.Repository
	outageRepo   outage.SnapshotStore
	logger       *log.Logger
}

// NewNotifyUsers creates a new NotifyUsers.
func NewNotifyUsers(fetchService *outage.FetchOutages, sender Sender, userRepo users.Repository, outageRepo outage.SnapshotStore, logger *log.Logger) *NotifyUsers {
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

	if prev != nil && outage.OutagesEqual(prev, outages) {
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

		content := Content{
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
