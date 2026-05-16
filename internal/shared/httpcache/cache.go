package httpcache

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
	ETag string
	Body []byte
}

type FetchResult struct {
	StatusCode int
	Body       []byte
	ETag       string
	FromCache  bool
}

func Load(path string) *Entry {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	parts := strings.SplitN(string(data), "\n", 2)
	if len(parts) < 2 {
		return nil
	}
	body := []byte(parts[1])
	if len(body) == 0 {
		return nil
	}
	return &Entry{ETag: parts[0], Body: body}
}

func Save(path string, etag string, body []byte) error {
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, append([]byte(etag+"\n"), body...), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return nil
}

func Get(ctx context.Context, client *http.Client, url string, cachePath string) (FetchResult, error) {
	var cache *Entry
	if cachePath != "" {
		cache = Load(cachePath)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return FetchResult{}, fmt.Errorf("create request: %w", err)
	}
	if cache != nil && cache.ETag != "" {
		req.Header.Set("If-None-Match", cache.ETag)
	}

	resp, err := client.Do(req)
	if err != nil {
		return FetchResult{}, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		if cache == nil {
			return FetchResult{}, fmt.Errorf("got 304 but no cached body")
		}
		return FetchResult{StatusCode: resp.StatusCode, Body: cache.Body, ETag: cache.ETag, FromCache: true}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FetchResult{}, fmt.Errorf("read response body: %w", err)
	}

	return FetchResult{
		StatusCode: resp.StatusCode,
		Body:       body,
		ETag:       resp.Header.Get("ETag"),
	}, nil
}
