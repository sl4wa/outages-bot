package cli

import (
	"context"
	"fmt"
	"log"
	"outages-bot/internal/application/notification"
)

// RunNotifierCommand fetches outages and sends notifications.
func RunNotifierCommand(
	ctx context.Context,
	fetchService *notification.OutageFetchService,
	notificationService *notification.NotificationService,
	logger *log.Logger,
) error {
	outages, err := fetchService.Handle(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch outages: %w", err)
	}

	count, err := notificationService.Handle(outages)
	if err != nil {
		return fmt.Errorf("failed to send notifications: %w", err)
	}

	logger.Printf("Successfully dispatched %d outages.", count)
	return nil
}
