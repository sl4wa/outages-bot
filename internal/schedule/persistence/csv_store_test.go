package persistence

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSVStateStoreLoadSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "outages.csv")
	store := NewCSVStateStore(path)
	date := time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)
	text := "line 1\nline \"2\""

	require.NoError(t, store.Save(map[time.Time]string{date: text}))
	loaded, err := store.Load()

	require.NoError(t, err)
	assert.Equal(t, map[time.Time]string{date: text}, loaded)
}

func TestCSVStateStoreMissingFile(t *testing.T) {
	loaded, err := NewCSVStateStore(filepath.Join(t.TempDir(), "missing.csv")).Load()

	require.NoError(t, err)
	assert.Empty(t, loaded)
}

func TestCSVStateStoreLoadSkipsMalformedRows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "outages.csv")
	content := "2026-02-13,good row\n" +
		"short-row\n" +
		"bad-date,ignored\n" +
		"2026-02-14,another good row\n"
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	loaded, err := NewCSVStateStore(path).Load()

	require.NoError(t, err)
	assert.Equal(t, map[time.Time]string{
		time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC): "good row",
		time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC): "another good row",
	}, loaded)
}
