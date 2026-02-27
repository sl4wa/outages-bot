package repository

import (
	"os"
	"outages-bot/internal/domain"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserRepo(t *testing.T) *FileUserRepository {
	t.Helper()
	dir := t.TempDir()
	repo, err := NewFileUserRepository(dir)
	require.NoError(t, err)
	return repo
}

func makeTestUser(t *testing.T, id int64) *domain.User {
	t.Helper()
	addr, err := domain.NewUserAddress(1, "Стрийська", "10")
	require.NoError(t, err)
	return &domain.User{ID: id, Address: addr}
}

func TestFileUserRepository_SaveAndFind(t *testing.T) {
	repo := setupUserRepo(t)
	user := makeTestUser(t, 12345)

	err := repo.Save(user)
	require.NoError(t, err)

	found, err := repo.Find(12345)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, int64(12345), found.ID)
	assert.Equal(t, "Стрийська", found.Address.StreetName)
	assert.Equal(t, "10", found.Address.Building)
	assert.Equal(t, 1, found.Address.StreetID)
	assert.Nil(t, found.OutageInfo)
}

func TestFileUserRepository_SaveWithOutageInfo(t *testing.T) {
	repo := setupUserRepo(t)
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	start := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC)
	period, _ := domain.NewOutagePeriod(start, end)
	desc := domain.NewOutageDescription("Планове відключення")
	info := domain.NewOutageInfo(period, desc)
	user := &domain.User{ID: 12345, Address: addr, OutageInfo: &info}

	err := repo.Save(user)
	require.NoError(t, err)

	found, err := repo.Find(12345)
	require.NoError(t, err)
	require.NotNil(t, found)
	require.NotNil(t, found.OutageInfo)
	assert.Equal(t, start.Unix(), found.OutageInfo.Period.StartDate.Unix())
	assert.Equal(t, end.Unix(), found.OutageInfo.Period.EndDate.Unix())
	assert.Equal(t, "Планове відключення", found.OutageInfo.Description.Value)
}

func TestFileUserRepository_FindNotFound(t *testing.T) {
	repo := setupUserRepo(t)
	found, err := repo.Find(99999)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestFileUserRepository_Remove(t *testing.T) {
	repo := setupUserRepo(t)
	user := makeTestUser(t, 12345)
	repo.Save(user)

	removed, err := repo.Remove(12345)
	require.NoError(t, err)
	assert.True(t, removed)

	found, err := repo.Find(12345)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestFileUserRepository_RemoveNotFound(t *testing.T) {
	repo := setupUserRepo(t)
	removed, err := repo.Remove(99999)
	assert.NoError(t, err)
	assert.False(t, removed)
}

func TestFileUserRepository_FindAll(t *testing.T) {
	repo := setupUserRepo(t)
	repo.Save(makeTestUser(t, 111))
	repo.Save(makeTestUser(t, 222))
	repo.Save(makeTestUser(t, 333))

	users := repo.FindAll()
	assert.Len(t, users, 3)
}

func TestFileUserRepository_AtomicWrite(t *testing.T) {
	repo := setupUserRepo(t)
	user := makeTestUser(t, 12345)
	err := repo.Save(user)
	require.NoError(t, err)

	// Verify no temp file remains
	tmpPath := filepath.Join(repo.dataDir, "12345.yml.tmp")
	_, err = os.Stat(tmpPath)
	assert.True(t, os.IsNotExist(err))

	// Verify the actual file exists
	actualPath := filepath.Join(repo.dataDir, "12345.yml")
	_, err = os.Stat(actualPath)
	assert.NoError(t, err)
}

func TestFileUserRepository_RaceCondition(t *testing.T) {
	repo := setupUserRepo(t)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			user := makeTestUser(t, id)
			err := repo.Save(user)
			assert.NoError(t, err)
		}(int64(i))
	}

	wg.Wait()

	users := repo.FindAll()
	assert.Len(t, users, 10)
}

func TestFileUserRepository_MigrateLegacyTxtToYml(t *testing.T) {
	dir := t.TempDir()

	// Write old-format .txt files
	simple := "street_id: 1\nstreet_name: Стрийська\nbuilding: 10\n"
	withColon := "street_id: 2\nstreet_name: Личаківська\nbuilding: 5\nstart_date: 2024-01-01T08:00:00Z\nend_date: 2024-01-01T16:00:00Z\ncomment: Причина: ремонт\n"

	require.NoError(t, os.WriteFile(filepath.Join(dir, "111.txt"), []byte(simple), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "222.txt"), []byte(withColon), 0o644))

	// Creating repo triggers migration
	repo, err := NewFileUserRepository(dir)
	require.NoError(t, err)

	// .txt files should be gone
	txtFiles, _ := filepath.Glob(filepath.Join(dir, "*.txt"))
	assert.Empty(t, txtFiles)

	// .yml files should exist
	ymlFiles, _ := filepath.Glob(filepath.Join(dir, "*.yml"))
	assert.Len(t, ymlFiles, 2)

	// Find should work
	user1, err := repo.Find(111)
	require.NoError(t, err)
	require.NotNil(t, user1)
	assert.Equal(t, "Стрийська", user1.Address.StreetName)
	assert.Equal(t, "10", user1.Address.Building)

	user2, err := repo.Find(222)
	require.NoError(t, err)
	require.NotNil(t, user2)
	assert.Equal(t, "Личаківська", user2.Address.StreetName)
	assert.Equal(t, "5", user2.Address.Building)
	require.NotNil(t, user2.OutageInfo)
	assert.Equal(t, "Причина: ремонт", user2.OutageInfo.Description.Value)

	// FindAll should return both
	users := repo.FindAll()
	assert.Len(t, users, 2)
}

func TestFileUserRepository_MigrateSkipsWhenYmlExists(t *testing.T) {
	dir := t.TempDir()

	// Write a stale .txt and a newer .yml with different data
	txt := "street_id: 1\nstreet_name: Стрийська\nbuilding: 10\n"
	yml := "street_id: 2\nstreet_name: Личаківська\nbuilding: 5\n"

	require.NoError(t, os.WriteFile(filepath.Join(dir, "111.txt"), []byte(txt), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "111.yml"), []byte(yml), 0o644))

	repo, err := NewFileUserRepository(dir)
	require.NoError(t, err)

	// .yml should retain its original content (not overwritten by .txt)
	user, err := repo.Find(111)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Личаківська", user.Address.StreetName)
	assert.Equal(t, "5", user.Address.Building)

	// .txt should still be present (skipped, not deleted)
	_, err = os.Stat(filepath.Join(dir, "111.txt"))
	assert.NoError(t, err)
}

func TestFileUserRepository_FindAllSkipsMalformedFiles(t *testing.T) {
	dir := t.TempDir()
	repo, err := NewFileUserRepository(dir)
	require.NoError(t, err)

	// Save a valid user
	repo.Save(makeTestUser(t, 111))

	// Write a malformed YAML file
	malformedPath := filepath.Join(dir, "222.yml")
	require.NoError(t, os.WriteFile(malformedPath, []byte("not: valid: yaml: [[["), 0o644))

	users := repo.FindAll()
	assert.Len(t, users, 1)
	assert.Equal(t, int64(111), users[0].ID)
}

func TestFileUserRepository_LoadFromFile_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	repo, err := NewFileUserRepository(dir)
	require.NoError(t, err)

	badPath := filepath.Join(dir, "999.yml")
	require.NoError(t, os.WriteFile(badPath, []byte(": : : invalid"), 0o644))

	user, err := repo.Find(999)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestFileUserRepository_LoadFromFile_InvalidAddress(t *testing.T) {
	dir := t.TempDir()
	repo, err := NewFileUserRepository(dir)
	require.NoError(t, err)

	// street_id of 0 is invalid
	badYAML := "street_id: 0\nstreet_name: Test\nbuilding: 10\n"
	badPath := filepath.Join(dir, "999.yml")
	require.NoError(t, os.WriteFile(badPath, []byte(badYAML), 0o644))

	user, err := repo.Find(999)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "invalid user address")
}

func TestFileUserRepository_LoadFromFile_InvalidOutageDates(t *testing.T) {
	dir := t.TempDir()
	repo, err := NewFileUserRepository(dir)
	require.NoError(t, err)

	badYAML := "street_id: 1\nstreet_name: Test\nbuilding: 10\nstart_date: not-a-date\nend_date: also-not\n"
	badPath := filepath.Join(dir, "999.yml")
	require.NoError(t, os.WriteFile(badPath, []byte(badYAML), 0o644))

	user, err := repo.Find(999)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "invalid start_date")
}

func TestFileUserRepository_SaveOverwriteExisting(t *testing.T) {
	repo := setupUserRepo(t)
	user := makeTestUser(t, 12345)
	require.NoError(t, repo.Save(user))

	// Overwrite with different address
	addr, err := domain.NewUserAddress(2, "Молдавська", "5")
	require.NoError(t, err)
	updatedUser := &domain.User{ID: 12345, Address: addr}
	require.NoError(t, repo.Save(updatedUser))

	found, err := repo.Find(12345)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "Молдавська", found.Address.StreetName)
	assert.Equal(t, "5", found.Address.Building)
	assert.Equal(t, 2, found.Address.StreetID)
}

func TestFileUserRepository_MigrateSkipsMissingStreetID(t *testing.T) {
	dir := t.TempDir()

	// .txt without street_id
	txt := "street_name: Стрийська\nbuilding: 10\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "111.txt"), []byte(txt), 0o644))

	_, err := NewFileUserRepository(dir)
	require.NoError(t, err)

	// .txt should still be present (not deleted)
	_, err = os.Stat(filepath.Join(dir, "111.txt"))
	assert.NoError(t, err)

	// .yml should NOT have been created
	_, err = os.Stat(filepath.Join(dir, "111.yml"))
	assert.True(t, os.IsNotExist(err))
}
