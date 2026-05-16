package loe

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/sl4wa/outages-bot/internal/shared/httpcache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPCache200RequiresCommitToSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "outages.http-cache")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"abc"`)
		_, _ = w.Write([]byte("body"))
	}))
	defer server.Close()
	cache := NewHTTPCache(path)
	cache.Client = server.Client()

	body, err := cache.Fetch(context.Background(), server.URL)

	require.NoError(t, err)
	assert.Equal(t, "body", body)
	assert.NoFileExists(t, path)
	require.NoError(t, cache.Commit())
	assert.FileExists(t, path)
}

func TestHTTPCache304UsesSavedBody(t *testing.T) {
	path := filepath.Join(t.TempDir(), "outages.http-cache")
	require.NoError(t, httpcache.Save(path, `"etag1"`, []byte("cached body")))
	var capturedETag string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedETag = r.Header.Get("If-None-Match")
		w.WriteHeader(http.StatusNotModified)
	}))
	defer server.Close()
	cache := NewHTTPCache(path)
	cache.Client = server.Client()

	body, err := cache.Fetch(context.Background(), server.URL)

	require.NoError(t, err)
	assert.Equal(t, "cached body", body)
	assert.Equal(t, `"etag1"`, capturedETag)
}

func TestHTTPCache304WithoutCacheThrows(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotModified)
	}))
	defer server.Close()
	cache := NewHTTPCache(filepath.Join(t.TempDir(), "missing"))
	cache.Client = server.Client()

	_, err := cache.Fetch(context.Background(), server.URL)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "304")
}

func TestHTTPCacheFetcherError(t *testing.T) {
	cache := NewHTTPCache("")
	cache.Client = &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("network failure")
	})}

	_, err := cache.Fetch(context.Background(), "http://example.test")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "network failure")
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
