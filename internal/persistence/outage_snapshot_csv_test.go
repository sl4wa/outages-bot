package persistence

import (
	"path/filepath"
	"testing"
	"time"

	"outages-bot/internal/outage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	t0 = time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	t1 = time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC)
)

func makeOutage(streetID int, streetName string, buildings []string, start, end time.Time, comment string) *outage.Outage {
	period, _ := outage.NewOutagePeriod(start, end)
	addr, _ := outage.NewOutageAddress(streetID, streetName, buildings, "Львів")
	return &outage.Outage{
		Period:      period,
		Address:     addr,
		Description: outage.NewOutageDescription(comment),
	}
}

func TestFileOutageRepository_Load_MissingFile_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	repo := NewFileOutageRepository(filepath.Join(dir, "snap.csv"))

	outages, err := repo.Load()
	require.NoError(t, err)
	assert.Nil(t, outages)
}

func TestFileOutageRepository_SaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	repo := NewFileOutageRepository(filepath.Join(dir, "snap.csv"))

	want := []*outage.Outage{makeOutage(1, "Стрийська", []string{"10", "12"}, t0, t1, "test")}

	require.NoError(t, repo.Save(want))

	got, err := repo.Load()
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, want[0].Address.StreetID, got[0].Address.StreetID)
	assert.Equal(t, want[0].Address.Buildings, got[0].Address.Buildings)
	assert.Equal(t, want[0].Description.Value, got[0].Description.Value)
	assert.Equal(t, want[0].Period.StartDate.Unix(), got[0].Period.StartDate.Unix())
	assert.Equal(t, want[0].Period.EndDate.Unix(), got[0].Period.EndDate.Unix())
}

func TestFileOutageRepository_Save_IsAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.csv")
	repo := NewFileOutageRepository(path)

	outages := []*outage.Outage{makeOutage(1, "Стрийська", []string{"1"}, t0, t1, "x")}
	require.NoError(t, repo.Save(outages))

	// temp file should not remain after save
	assert.NoFileExists(t, path+".tmp")
	assert.FileExists(t, path)
}

func TestFileOutageRepository_UTCNormalization(t *testing.T) {
	dir := t.TempDir()
	repo := NewFileOutageRepository(filepath.Join(dir, "snap.csv"))

	kyiv := time.FixedZone("UTC+3", 3*60*60)
	outages := []*outage.Outage{makeOutage(1, "Стрийська", []string{"10"}, t0.In(kyiv), t1.In(kyiv), "c")}
	require.NoError(t, repo.Save(outages))

	got, err := repo.Load()
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, t0.Unix(), got[0].Period.StartDate.Unix())
	assert.Equal(t, t1.Unix(), got[0].Period.EndDate.Unix())
}
