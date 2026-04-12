package notifier

import (
	"errors"
	"outages-bot/internal/domain"
)

type mockOutageRepo struct {
	outages []*domain.Outage
	saved   []*domain.Outage
	loadErr error
	saveErr error
}

func (m *mockOutageRepo) Load() ([]*domain.Outage, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	return m.outages, nil
}

func (m *mockOutageRepo) Save(outages []*domain.Outage) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.saved = outages
	m.outages = outages
	return nil
}

var errSaveFailed = errors.New("save failed")
