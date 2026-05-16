package loe

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sl4wa/outages-bot/internal/shared/httpcache"
)

const DefaultCacheFileName = "schedule.http-cache"

type HTTPCache struct {
	Path    string
	Client  *http.Client
	pending *httpcache.FetchResult
}

func NewHTTPCache(path string) *HTTPCache {
	return &HTTPCache{
		Path:   path,
		Client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *HTTPCache) Fetch(ctx context.Context, url string) (string, error) {
	c.pending = nil
	result, err := httpcache.Get(ctx, c.Client, url, c.Path)
	if err != nil {
		return "", fmt.Errorf("LOE API request failed: %w", err)
	}
	switch {
	case result.StatusCode == http.StatusNotModified:
		return string(result.Body), nil
	case result.StatusCode >= 200 && result.StatusCode <= 299:
		c.pending = &result
		return string(result.Body), nil
	default:
		return "", fmt.Errorf("LOE API request failed with status %d", result.StatusCode)
	}
}

func (c *HTTPCache) Commit() error {
	if c.pending == nil || c.Path == "" {
		c.pending = nil
		return nil
	}
	err := httpcache.Save(c.Path, c.pending.ETag, c.pending.Body)
	c.pending = nil
	return err
}
