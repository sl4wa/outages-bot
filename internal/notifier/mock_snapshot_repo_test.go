package notifier

import (
	"errors"
	"outages-bot/internal/outage"
)

type mockOutageRepo struct {
	outages []*outage.Outage
	saved   []*outage.Outage
	loadErr error
	saveErr error
}

func (m *mockOutageRepo) Load() ([]*outage.Outage, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	return m.outages, nil
}

func (m *mockOutageRepo) Save(outages []*outage.Outage) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.saved = outages
	m.outages = outages
	return nil
}

var errSaveFailed = errors.New("save failed")
