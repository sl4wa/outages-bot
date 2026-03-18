package cli

import (
	"context"
	"log"
	"outages-bot/internal/application/notifier"
)

// RunNotifierCommand fetches outages and sends notifications.
func RunNotifierCommand(
	ctx context.Context,
	notifyUsers *notifier.NotifyUsers,
	logger *log.Logger,
) error {
	if err := notifyUsers.Handle(ctx); err != nil {
		return err
	}

	logger.Printf("Successfully dispatched notifications.")
	return nil
}
