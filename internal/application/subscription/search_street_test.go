package subscription

import (
	"outages-bot/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockStreetRepo struct {
	streets []domain.Street
}

func (m *mockStreetRepo) GetAllStreets() []domain.Street {
	return m.streets
}

func newSearchService() *SearchStreetService {
	repo := &mockStreetRepo{
		streets: []domain.Street{
			{ID: 1, Name: "Стрийська"},
			{ID: 2, Name: "Наукова"},
			{ID: 3, Name: "Стрілецька"},
		},
	}
	return NewSearchStreetService(repo)
}

func TestSearchStreet_EmptyQuery(t *testing.T) {
	svc := newSearchService()
	result := svc.Handle("")
	assert.Equal(t, "Введіть назву вулиці.", result.Message)
	assert.False(t, result.HasExactMatch())
	assert.False(t, result.HasMultipleOptions())
}

func TestSearchStreet_WhitespaceQuery(t *testing.T) {
	svc := newSearchService()
	result := svc.Handle("   ")
	assert.Equal(t, "Введіть назву вулиці.", result.Message)
}

func TestSearchStreet_NotFound(t *testing.T) {
	svc := newSearchService()
	result := svc.Handle("Невідома")
	assert.Equal(t, "Вулицю не знайдено. Спробуйте ще раз.", result.Message)
	assert.False(t, result.HasExactMatch())
}

func TestSearchStreet_ExactMatch(t *testing.T) {
	svc := newSearchService()
	result := svc.Handle("Наукова")
	assert.True(t, result.HasExactMatch())
	assert.Equal(t, 2, *result.SelectedStreetID)
	assert.Equal(t, "Наукова", *result.SelectedStreetName)
	assert.Contains(t, result.Message, "Ви обрали вулицю: Наукова")
	assert.Contains(t, result.Message, "Будь ласка, введіть номер будинку:")
}

func TestSearchStreet_ExactMatchCaseInsensitive(t *testing.T) {
	svc := newSearchService()
	result := svc.Handle("наукова")
	assert.True(t, result.HasExactMatch())
	assert.Equal(t, 2, *result.SelectedStreetID)
}

func TestSearchStreet_SinglePartialMatch(t *testing.T) {
	svc := newSearchService()
	result := svc.Handle("науков")
	assert.True(t, result.HasExactMatch())
	assert.Equal(t, 2, *result.SelectedStreetID)
}

func TestSearchStreet_MultipleMatches(t *testing.T) {
	svc := newSearchService()
	result := svc.Handle("стр")
	assert.True(t, result.HasMultipleOptions())
	assert.False(t, result.HasExactMatch())
	assert.Equal(t, "Будь ласка, оберіть вулицю:", result.Message)
	assert.Len(t, result.StreetOptions, 2)
}

func TestSearchStreet_MultipleMatchesContainsCorrectStreets(t *testing.T) {
	svc := newSearchService()
	result := svc.Handle("стр")
	names := make([]string, len(result.StreetOptions))
	for i, s := range result.StreetOptions {
		names[i] = s.Name
	}
	assert.Contains(t, names, "Стрийська")
	assert.Contains(t, names, "Стрілецька")
}

func TestSearchStreet_ExactMatchTakesPrecedence(t *testing.T) {
	svc := newSearchService()
	result := svc.Handle("Стрийська")
	assert.True(t, result.HasExactMatch())
	assert.Equal(t, 1, *result.SelectedStreetID)
}
