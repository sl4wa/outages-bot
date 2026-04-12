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
	return notifyUsers.Handle(ctx)
}
