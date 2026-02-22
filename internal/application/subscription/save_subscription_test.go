package subscription

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveSubscription_ValidInput(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewSaveSubscriptionService(repo)
	result := svc.Handle(12345, 1, "Стрийська", "10")
	assert.True(t, result.Success)
	assert.Contains(t, result.Message, "Ви підписалися")
	assert.Contains(t, result.Message, "Стрийська")
	assert.Contains(t, result.Message, "10")
	assert.Len(t, repo.saved, 1)
}

func TestSaveSubscription_InvalidBuilding(t *testing.T) {
	tests := []struct {
		name     string
		building string
	}{
		{"lowercase suffix", "10-a"},
		{"slash format", "10/1"},
		{"letters only", "abc"},
		{"empty", ""},
		{"with space", "10 A"},
		{"special chars", "10!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepo()
			svc := NewSaveSubscriptionService(repo)
			result := svc.Handle(12345, 1, "Стрийська", tt.building)
			assert.False(t, result.Success)
			assert.NotEmpty(t, result.Message)
			// Verify Save was NOT called for invalid input
			assert.Len(t, repo.saved, 0)
		})
	}
}

func TestSaveSubscription_InvalidStreetID(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewSaveSubscriptionService(repo)
	result := svc.Handle(12345, 0, "Стрийська", "10")
	assert.False(t, result.Success)
	assert.Len(t, repo.saved, 0)
}

func TestSaveSubscription_WithCyrillicSuffix(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewSaveSubscriptionService(repo)
	result := svc.Handle(12345, 1, "Стрийська", "10-А")
	assert.True(t, result.Success)
	assert.Contains(t, result.Message, "10-А")
}

func TestSaveSubscription_MessageFormat(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewSaveSubscriptionService(repo)
	result := svc.Handle(12345, 1, "Наукова", "25")
	assert.True(t, result.Success)
	assert.Equal(t,
		"Ви підписалися на сповіщення про відключення електроенергії для вулиці Наукова, будинок 25.",
		result.Message,
	)
}
