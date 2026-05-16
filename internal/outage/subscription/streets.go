package subscription

import (
	"strings"

	"github.com/sl4wa/outages-bot/internal/outage/users"
)

type streetSearchResult struct {
	street  *users.Street
	options []users.Street
}

func (w *Workflow) searchStreet(query string) (streetSearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return streetSearchResult{}, ErrEmptyStreetQuery
	}

	q := strings.ToLower(query)
	var matches []users.Street
	for _, street := range w.streetRepo.GetAllStreets() {
		if street.NameEquals(q) {
			match := street
			return streetSearchResult{street: &match}, nil
		}
		if street.NameContains(q) {
			matches = append(matches, street)
		}
	}

	switch len(matches) {
	case 0:
		return streetSearchResult{}, ErrStreetNotFound
	case 1:
		match := matches[0]
		return streetSearchResult{street: &match}, nil
	default:
		return streetSearchResult{options: matches}, nil
	}
}
