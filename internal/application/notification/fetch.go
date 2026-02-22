package notification

import (
	"context"
	"outages-bot/internal/application"
	"outages-bot/internal/domain"
)

// OutageFetchService converts provider DTOs into domain Outage entities.
type OutageFetchService struct {
	provider application.OutageProvider
}

// NewOutageFetchService creates a new OutageFetchService.
func NewOutageFetchService(provider application.OutageProvider) *OutageFetchService {
	return &OutageFetchService{provider: provider}
}

// Handle fetches outages from the provider and converts them to domain entities.
func (s *OutageFetchService) Handle(ctx context.Context) ([]*domain.Outage, error) {
	dtos, err := s.provider.FetchOutages(ctx)
	if err != nil {
		return nil, err
	}

	outages := make([]*domain.Outage, 0, len(dtos))
	for _, dto := range dtos {
		period, err := domain.NewOutagePeriod(dto.Start, dto.End)
		if err != nil {
			continue
		}
		addr, err := domain.NewOutageAddress(dto.StreetID, dto.StreetName, dto.Buildings, dto.City)
		if err != nil {
			continue
		}
		outages = append(outages, &domain.Outage{
			ID:          dto.ID,
			Period:      period,
			Address:     addr,
			Description: domain.NewOutageDescription(dto.Comment),
		})
	}

	return outages, nil
}
