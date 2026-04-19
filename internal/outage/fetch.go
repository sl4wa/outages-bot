package outage

import (
	"context"
)

// FetchOutages converts provider DTOs into domain Outage entities.
type FetchOutages struct {
	provider RawProvider
}

// NewFetchOutages creates a new FetchOutages.
func NewFetchOutages(provider RawProvider) *FetchOutages {
	return &FetchOutages{provider: provider}
}

// Handle fetches outages from the provider and converts them to domain entities.
func (s *FetchOutages) Handle(ctx context.Context) ([]*Outage, error) {
	dtos, err := s.provider.FetchOutages(ctx)
	if err != nil {
		return nil, err
	}

	return dtosToOutages(dtos), nil
}

func dtosToOutages(dtos []RawOutage) []*Outage {
	outages := make([]*Outage, 0, len(dtos))
	for _, dto := range dtos {
		period, err := NewOutagePeriod(dto.Start, dto.End)
		if err != nil {
			continue
		}
		addr, err := NewOutageAddress(dto.StreetID, dto.StreetName, dto.Buildings, dto.City)
		if err != nil {
			continue
		}
		outages = append(outages, &Outage{
			ID:          dto.ID,
			Period:      period,
			Address:     addr,
			Description: NewOutageDescription(dto.Comment),
		})
	}
	return outages
}
