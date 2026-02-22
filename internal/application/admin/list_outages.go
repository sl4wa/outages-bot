package admin

import (
	"context"
	"outages-bot/internal/application"
)

// ListOutages fetches all outages from the provider.
func ListOutages(ctx context.Context, provider application.OutageProvider) ([]application.OutageDTO, error) {
	return provider.FetchOutages(ctx)
}
