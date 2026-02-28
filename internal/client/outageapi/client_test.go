package outageapi

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestApiProvider_Non200_ReturnsEmptyAndLogsWarning(t *testing.T) {
	server := makeServer(t, 500, "error")
	defer server.Close()

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	provider := NewProvider(server.URL, fixedClock(), logger)

	result, err := provider.FetchOutages(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Contains(t, buf.String(), "500")
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
