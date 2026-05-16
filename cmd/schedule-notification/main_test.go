package main

import (
	"os"
	"path/filepath"
	"testing"

	outageloe "github.com/sl4wa/outages-bot/internal/outage/loe"
	outagepersistence "github.com/sl4wa/outages-bot/internal/outage/persistence"
	"github.com/sl4wa/outages-bot/internal/schedule/loe"
	"github.com/sl4wa/outages-bot/internal/schedule/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigFromEnvUsesDataDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("SCHEDULE_API_URL", "https://example.test/api")
	t.Setenv("DATA_DIR", dir)

	config, err := configFromEnv()

	require.NoError(t, err)
	expectedDataDir, err := filepath.Abs(dir)
	require.NoError(t, err)
	expectedDataDir = filepath.Clean(expectedDataDir)
	assert.Equal(t, filepath.Join(expectedDataDir, persistence.StateFileName), config.StatePath)
	assert.Equal(t, filepath.Join(expectedDataDir, "users"), config.TelegramUsersDir)
	assert.Equal(t, filepath.Join(expectedDataDir, loe.DefaultCacheFileName), config.HTTPCachePath)
	assert.Equal(t, "https://example.test/api", config.APIURL)
}

func TestConfigFromEnvRequiresDataDir(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("SCHEDULE_API_URL", "https://example.test/api")
	t.Setenv("DATA_DIR", "")

	_, err := configFromEnv()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "DATA_DIR")
}

func TestConfigFromEnvRequiresToken(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "")
	t.Setenv("DATA_DIR", t.TempDir())

	_, err := configFromEnv()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "TELEGRAM_BOT_TOKEN")
}

func TestConfigFromEnvRequiresAPIURL(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("SCHEDULE_API_URL", "")
	t.Setenv("DATA_DIR", t.TempDir())

	_, err := configFromEnv()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "SCHEDULE_API_URL")
}

func TestConfigFromEnvUsesCustomDataDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("SCHEDULE_API_URL", "https://example.test/api")
	t.Setenv("DATA_DIR", dir)

	config, err := configFromEnv()

	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, persistence.StateFileName), config.StatePath)
	assert.Equal(t, filepath.Join(dir, "users"), config.TelegramUsersDir)
	assert.Equal(t, filepath.Join(dir, loe.DefaultCacheFileName), config.HTTPCachePath)
}

func TestConfigFromEnvDoesNotCollideWithOutageFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("SCHEDULE_API_URL", "https://example.test/api")
	t.Setenv("DATA_DIR", dir)

	config, err := configFromEnv()

	require.NoError(t, err)
	assert.NotEqual(t, filepath.Join(dir, outagepersistence.OutageSnapshotFileName), config.StatePath)
	assert.NotEqual(t, filepath.Join(dir, outageloe.DefaultCacheFileName), config.HTTPCachePath)
}

func TestAbsPathKeepsCleanAbsolutePath(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(wd, "x"), absPath("./x"))
}
