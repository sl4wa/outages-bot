package repository

import (
	"os"
	"path/filepath"
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

func TestFileStreetRepository_FileNotFound(t *testing.T) {
	_, err := NewFileStreetRepository("nonexistent/streets.csv")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open streets file")
}

func TestFileStreetRepository_EmptyCSV(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.csv")
	require.NoError(t, os.WriteFile(path, []byte("id,name\n"), 0o644))

	repo, err := NewFileStreetRepository(path)
	require.NoError(t, err)

	streets, err := repo.GetAllStreets()
	require.NoError(t, err)
	assert.Empty(t, streets)
}

func TestFileStreetRepository_InvalidStreetID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.csv")
	require.NoError(t, os.WriteFile(path, []byte("id,name\nabc,Стрийська\n"), 0o644))

	_, err := NewFileStreetRepository(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid street id")
}

func TestFileStreetRepository_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.csv")
	require.NoError(t, os.WriteFile(path, []byte(""), 0o644))

	repo, err := NewFileStreetRepository(path)
	require.NoError(t, err)

	streets, err := repo.GetAllStreets()
	require.NoError(t, err)
	assert.Empty(t, streets)
}
