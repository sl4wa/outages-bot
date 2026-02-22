package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStreetRepository_LoadsStreets(t *testing.T) {
	repo, err := NewFileStreetRepository("testdata/streets.csv")
	require.NoError(t, err)

	streets, err := repo.GetAllStreets()
	require.NoError(t, err)
	assert.Len(t, streets, 2)
}

func TestFileStreetRepository_StreetsContent(t *testing.T) {
	repo, err := NewFileStreetRepository("testdata/streets.csv")
	require.NoError(t, err)

	streets, _ := repo.GetAllStreets()
	assert.Equal(t, 12444, streets[0].ID)
	assert.Equal(t, "Молдавська", streets[0].Name)
	assert.Equal(t, 12445, streets[1].ID)
	assert.Equal(t, "Стрийська", streets[1].Name)
}
