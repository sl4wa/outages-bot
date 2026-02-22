package subscription

import (
	"fmt"
	"outages-bot/internal/domain"
	"strings"
)

// SearchStreetResult holds the result of a street search.
type SearchStreetResult struct {
	Message            string
	StreetOptions      []domain.Street
	SelectedStreetID   *int
	SelectedStreetName *string
}

// HasMultipleOptions returns true if multiple street options are available.
func (r *SearchStreetResult) HasMultipleOptions() bool {
	return len(r.StreetOptions) > 0
}

// HasExactMatch returns true if exactly one street was matched.
func (r *SearchStreetResult) HasExactMatch() bool {
	return r.SelectedStreetID != nil && r.SelectedStreetName != nil
}

// SearchStreetService handles street search logic.
type SearchStreetService struct {
	streetRepo domain.StreetRepository
}

// NewSearchStreetService creates a new SearchStreetService.
func NewSearchStreetService(streetRepo domain.StreetRepository) *SearchStreetService {
	return &SearchStreetService{streetRepo: streetRepo}
}

// Handle searches for streets matching the query.
func (s *SearchStreetService) Handle(query string) (*SearchStreetResult, error) {
	query = strings.TrimSpace(query)

	if query == "" {
		return &SearchStreetResult{Message: "Введіть назву вулиці."}, nil
	}

	q := strings.ToLower(query)
	streets, err := s.streetRepo.GetAllStreets()
	if err != nil {
		return nil, err
	}

	var results []domain.Street
	for _, street := range streets {
		if street.NameEquals(q) {
			id := street.ID
			name := street.Name
			return &SearchStreetResult{
				Message:            fmt.Sprintf("Ви обрали вулицю: %s\nБудь ласка, введіть номер будинку:", street.Name),
				SelectedStreetID:   &id,
				SelectedStreetName: &name,
			}, nil
		}
		if street.NameContains(q) {
			results = append(results, street)
		}
	}

	if len(results) == 0 {
		return &SearchStreetResult{Message: "Вулицю не знайдено. Спробуйте ще раз."}, nil
	}

	if len(results) == 1 {
		id := results[0].ID
		name := results[0].Name
		return &SearchStreetResult{
			Message:            fmt.Sprintf("Ви обрали вулицю: %s\nБудь ласка, введіть номер будинку:", results[0].Name),
			SelectedStreetID:   &id,
			SelectedStreetName: &name,
		}, nil
	}

	return &SearchStreetResult{
		Message:       "Будь ласка, оберіть вулицю:",
		StreetOptions: results,
	}, nil
}
