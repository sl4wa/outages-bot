package httpcache

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache")

	require.NoError(t, Save(path, `"abc"`, []byte("body")))
	entry := Load(path)

	require.NotNil(t, entry)
	assert.Equal(t, `"abc"`, entry.ETag)
	assert.Equal(t, []byte("body"), entry.Body)
}

func TestLoadMalformedReturnsNil(t *testing.T) {
	assert.Nil(t, Load(filepath.Join(t.TempDir(), "missing")))
}

func TestGet200ReturnsFreshBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"v1"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("fresh"))
	}))
	defer server.Close()

	result, err := Get(context.Background(), server.Client(), server.URL, "")

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, []byte("fresh"), result.Body)
	assert.Equal(t, `"v1"`, result.ETag)
	assert.False(t, result.FromCache)
}

func TestGet304UsesCacheAndSendsETag(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache")
	require.NoError(t, Save(path, `"v1"`, []byte("cached")))
	var gotETag string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotETag = r.Header.Get("If-None-Match")
		w.WriteHeader(http.StatusNotModified)
	}))
	defer server.Close()

	result, err := Get(context.Background(), server.Client(), server.URL, path)

	require.NoError(t, err)
	assert.Equal(t, `"v1"`, gotETag)
	assert.Equal(t, []byte("cached"), result.Body)
	assert.True(t, result.FromCache)
}

func TestGet304WithoutCacheFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotModified)
	}))
	defer server.Close()

	_, err := Get(context.Background(), server.Client(), server.URL, filepath.Join(t.TempDir(), "missing"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "304")
}

func TestGetNetworkError(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("network")
	})}

	_, err := Get(context.Background(), client, "http://example.test", "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "network")
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
