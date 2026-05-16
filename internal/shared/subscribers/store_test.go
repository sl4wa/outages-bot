package subscribers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilePath(t *testing.T) {
	s := NewFileStore("/data/users")
	assert.Equal(t, filepath.Join("/data/users", "12345.toml"), s.FilePath(12345))
}

func TestChatIDsSortedAndFiltered(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "222.toml"), nil, 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "111.toml"), nil, 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "-100123.toml"), nil, 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "abc.toml"), nil, 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "333.yml"), nil, 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "444.toml"), 0o755))

	ids, err := NewFileStore(dir).ChatIDs()

	require.NoError(t, err)
	assert.Equal(t, []int64{-100123, 111, 222}, ids)
}

func TestChatIDsRejectsNonCanonicalNumericNames(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "001.toml"), nil, 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "+123.toml"), nil, 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "42.toml"), nil, 0o644))

	ids, err := NewFileStore(dir).ChatIDs()

	require.NoError(t, err)
	assert.Equal(t, []int64{42}, ids)
}

func TestChatIDsMissingDir(t *testing.T) {
	ids, err := NewFileStore(filepath.Join(t.TempDir(), "missing")).ChatIDs()

	require.NoError(t, err)
	assert.Nil(t, ids)
}

func TestChatIDsUnreadableDirErrors(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "users")
	require.NoError(t, os.Mkdir(dir, 0o000))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })

	ids, err := NewFileStore(dir).ChatIDs()

	require.Error(t, err)
	assert.Nil(t, ids)
}

func TestWriteCreatesDirAndFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "users")
	s := NewFileStore(dir)

	require.NoError(t, s.Write(42, []byte("hello")))

	got, err := os.ReadFile(filepath.Join(dir, "42.toml"))
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), got)

	_, err = os.Stat(filepath.Join(dir, "42.toml.tmp"))
	assert.True(t, os.IsNotExist(err))
}

func TestWriteOverwrites(t *testing.T) {
	dir := t.TempDir()
	s := NewFileStore(dir)

	require.NoError(t, s.Write(42, []byte("first")))
	require.NoError(t, s.Write(42, []byte("second")))

	got, err := s.Read(42)
	require.NoError(t, err)
	assert.Equal(t, []byte("second"), got)
}

func TestReadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewFileStore(dir)
	require.NoError(t, s.Write(7, []byte("payload")))

	got, err := s.Read(7)

	require.NoError(t, err)
	assert.Equal(t, []byte("payload"), got)
}

func TestReadMissingIsNotExist(t *testing.T) {
	_, err := NewFileStore(t.TempDir()).Read(99)

	require.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrNotExist))
}

func TestRemoveExisting(t *testing.T) {
	dir := t.TempDir()
	s := NewFileStore(dir)
	require.NoError(t, s.Write(7, []byte("x")))

	removed, err := s.Remove(7)

	require.NoError(t, err)
	assert.True(t, removed)

	_, err = s.Read(7)
	assert.True(t, errors.Is(err, os.ErrNotExist))
}

func TestRemoveMissing(t *testing.T) {
	removed, err := NewFileStore(t.TempDir()).Remove(99)

	require.NoError(t, err)
	assert.False(t, removed)
}
