package loe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixedClock() func() time.Time {
	t := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	return func() time.Time { return t }
}

func makeServer(t *testing.T, statusCode int, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}))
}

const validBody = `{"hydra:member":[{"id":1,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"test","buildingNames":"10","city":{"name":"Львів"},"street":{"id":1,"name":"Стрийська"}}]}`

func TestApiProvider_Non200_ReturnsError(t *testing.T) {
	server := makeServer(t, 500, "error")
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)

	result, err := provider.FetchOutages(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
	assert.Nil(t, result)
}

func TestApiProvider_MalformedJSON_ReturnsError(t *testing.T) {
	server := makeServer(t, 200, "not json")
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	_, err := provider.FetchOutages(context.Background())
	assert.Error(t, err)
}

func TestApiProvider_MissingHydraMember_ReturnsEmpty(t *testing.T) {
	server := makeServer(t, 200, `{}`)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestApiProvider_BuildingNamesAsString(t *testing.T) {
	body := `{"hydra:member":[{"id":1,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"test","buildingNames":"10, 12, 14","city":{"name":"Львів"},"street":{"id":1,"name":"Стрийська"}}]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, []string{"10", "12", "14"}, result[0].Buildings)
}

func TestApiProvider_BuildingNamesAsArray(t *testing.T) {
	body := `{"hydra:member":[{"id":1,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"test","buildingNames":["10","12"],"city":{"name":"Львів"},"street":{"id":1,"name":"Стрийська"}}]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, []string{"10", "12"}, result[0].Buildings)
}

func TestApiProvider_Dedup(t *testing.T) {
	body := `{"hydra:member":[
		{"id":1,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"a","buildingNames":"10","city":{"name":"Львів"},"street":{"id":1,"name":"S"}},
		{"id":2,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"b","buildingNames":"10","city":{"name":"Львів"},"street":{"id":1,"name":"S"}}
	]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestApiProvider_MissingDates_UsesInjectedClock(t *testing.T) {
	body := `{"hydra:member":[{"id":1,"dateEvent":"","datePlanIn":"","koment":"test","buildingNames":"10","city":{"name":"Львів"},"street":{"id":1,"name":"S"}}]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	clock := fixedClock()
	provider := NewProvider(server.URL, clock, nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, clock().Unix(), result[0].Start.Unix())
	assert.Equal(t, clock().Unix(), result[0].End.Unix())
}

func TestApiProvider_EmptyCityAndStreet(t *testing.T) {
	body := `{"hydra:member":[{"id":1,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"test","buildingNames":"10"}]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "", result[0].City)
	assert.Equal(t, "", result[0].StreetName)
	assert.Equal(t, 0, result[0].StreetID)
}

func TestApiProvider_CommentWithCRLF(t *testing.T) {
	body := `{"hydra:member":[{"id":1,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"line1\r\nline2","buildingNames":"10","street":{"id":1,"name":"S"}}]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "line1 line2", result[0].Comment)
}

func TestApiProvider_CommentWithMultipleNewlines(t *testing.T) {
	body := `{"hydra:member":[{"id":1,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"line1\n\n\nline2","buildingNames":"10","street":{"id":1,"name":"S"}}]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "line1 line2", result[0].Comment)
}

func TestApiProvider_CommentWithLeadingTrailingWhitespace(t *testing.T) {
	body := `{"hydra:member":[{"id":1,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"  hello  ","buildingNames":"10","street":{"id":1,"name":"S"}}]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "hello", result[0].Comment)
}

func TestApiProvider_NormalComment(t *testing.T) {
	body := `{"hydra:member":[{"id":1,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"Normal comment","buildingNames":"10","street":{"id":1,"name":"S"}}]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Normal comment", result[0].Comment)
}

func TestApiProvider_StringIDs(t *testing.T) {
	body := `{"hydra:member":[{"id":"123","dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"test","buildingNames":"10","city":{"name":"Львів"},"street":{"id":"456","name":"Стрийська"}}]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 123, result[0].ID)
	assert.Equal(t, 456, result[0].StreetID)
}

func TestApiProvider_DecimalStringIDs(t *testing.T) {
	body := `{"hydra:member":[{"id":"123.0","dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"test","buildingNames":"10","city":{"name":"Львів"},"street":{"id":"456.0","name":"Стрийська"}}]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 123, result[0].ID)
	assert.Equal(t, 456, result[0].StreetID)
}

func TestApiProvider_DedupPreservesOrder(t *testing.T) {
	// A and C share the same dedup key (same street, buildings, dates).
	// B has a different street. Result should be [C, B] — C overwrites A at index 0.
	body := `{"hydra:member":[
		{"id":1,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"a","buildingNames":"10","city":{"name":"Львів"},"street":{"id":1,"name":"S1"}},
		{"id":2,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"b","buildingNames":"10","city":{"name":"Львів"},"street":{"id":2,"name":"S2"}},
		{"id":3,"dateEvent":"2024-01-01T08:00:00+00:00","datePlanIn":"2024-01-01T16:00:00+00:00","koment":"c","buildingNames":"10","city":{"name":"Львів"},"street":{"id":1,"name":"S1"}}
	]}`
	server := makeServer(t, 200, body)
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 2)
	// Index 0: C overwrote A (last-wins at original position)
	assert.Equal(t, 3, result[0].ID)
	assert.Equal(t, "c", result[0].Comment)
	assert.Equal(t, "S1", result[0].StreetName)
	// Index 1: B unchanged
	assert.Equal(t, 2, result[1].ID)
	assert.Equal(t, "b", result[1].Comment)
	assert.Equal(t, "S2", result[1].StreetName)
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{"json.Number int", json.Number("123"), 123},
		{"json.Number float", json.Number("123.0"), 123},
		{"json.Number zero", json.Number("0"), 0},
		{"json.Number invalid", json.Number("abc"), 0},
		{"float64", float64(42.9), 42},
		{"int", int(7), 7},
		{"string int", "99", 99},
		{"string float with spaces", " 55.5 ", 55},
		{"string invalid", "notanumber", 0},
		{"nil", nil, 0},
		{"bool", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, toInt(tt.input))
		})
	}
}

func TestProvider_200_WritesCacheAfterSuccessfulParse(t *testing.T) {
	dir := t.TempDir()
	cacheFile := filepath.Join(dir, "outages.http-cache")
	etag := `"abc123"`

	callCount := 0
	var secondCallHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Header().Set("ETag", etag)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(validBody))
		} else {
			secondCallHeader = r.Header.Get("If-None-Match")
			w.WriteHeader(http.StatusNotModified)
		}
	}))
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil).WithCacheFile(cacheFile)

	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.FileExists(t, cacheFile)

	result2, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	require.Len(t, result2, 1)
	assert.Equal(t, etag, secondCallHeader)
}

func TestProvider_304_WithCache_ReturnsCachedBody(t *testing.T) {
	dir := t.TempDir()
	cacheFile := filepath.Join(dir, "outages.http-cache")
	etag := `"v1"`
	require.NoError(t, saveCache(cacheFile, etag, []byte(validBody)))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotModified)
	}))
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil).WithCacheFile(cacheFile)
	result, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Стрийська", result[0].StreetName)
}

func TestProvider_304_WithoutCache_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	cacheFile := filepath.Join(dir, "outages.http-cache")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotModified)
	}))
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil).WithCacheFile(cacheFile)
	_, err := provider.FetchOutages(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "304")
}

func TestProvider_MalformedCacheFile_TreatedAsMiss(t *testing.T) {
	dir := t.TempDir()
	cacheFile := filepath.Join(dir, "outages.http-cache")
	require.NoError(t, os.WriteFile(cacheFile, []byte("no-newline-content"), 0o644))

	var capturedHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Get("If-None-Match")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(validBody))
	}))
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil).WithCacheFile(cacheFile)
	_, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	assert.Empty(t, capturedHeader)
}

func TestProvider_EmptyBodyCache_TreatedAsMiss(t *testing.T) {
	dir := t.TempDir()
	cacheFile := filepath.Join(dir, "outages.http-cache")
	// ETag\n with empty body
	require.NoError(t, os.WriteFile(cacheFile, []byte("\"v1\"\n"), 0o644))

	var capturedHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Get("If-None-Match")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(validBody))
	}))
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil).WithCacheFile(cacheFile)
	_, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	assert.Empty(t, capturedHeader)
}

func TestProvider_NoCacheFile_NoEtagHeader(t *testing.T) {
	var capturedHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Get("If-None-Match")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(validBody))
	}))
	defer server.Close()

	provider := NewProvider(server.URL, fixedClock(), nil)
	_, err := provider.FetchOutages(context.Background())
	require.NoError(t, err)
	assert.Empty(t, capturedHeader)
}

func TestProvider_200_MalformedBody_PreservesExistingCache(t *testing.T) {
	dir := t.TempDir()
	cacheFile := filepath.Join(dir, "outages.http-cache")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"new"`)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	require.NoError(t, saveCache(cacheFile, `"old"`, []byte(validBody)))
	originalContent, err := os.ReadFile(cacheFile)
	require.NoError(t, err)

	provider := NewProvider(server.URL, fixedClock(), nil).WithCacheFile(cacheFile)
	_, err = provider.FetchOutages(context.Background())
	require.Error(t, err)

	afterContent, err := os.ReadFile(cacheFile)
	require.NoError(t, err)
	assert.Equal(t, originalContent, afterContent)
}
