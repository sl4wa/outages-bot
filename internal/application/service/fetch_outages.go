package service

import (
	"context"
	"outages-bot/internal/domain"
	"time"
)

// OutageProvider fetches outage data from an external source.
type OutageProvider interface {
	FetchOutages(ctx context.Context) ([]OutageDTO, error)
}

// OutageDTO is a data transfer object for outages from the provider.
type OutageDTO struct {
	ID         int
	Start      time.Time
	End        time.Time
	City       string
	StreetID   int
	StreetName string
	Buildings  []string
	Comment    string
}

// FetchOutages converts provider DTOs into domain Outage entities.
type FetchOutages struct {
	provider OutageProvider
}

// NewFetchOutages creates a new FetchOutages.
func NewFetchOutages(provider OutageProvider) *FetchOutages {
	return &FetchOutages{provider: provider}
}

// Handle fetches outages from the provider and converts them to domain entities.
func (s *FetchOutages) Handle(ctx context.Context) ([]*domain.Outage, error) {
	dtos, err := s.provider.FetchOutages(ctx)
	if err != nil {
		return nil, err
	}

	return dtosToOutages(dtos), nil
}

func dtosToOutages(dtos []OutageDTO) []*domain.Outage {
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
	return outages
}
