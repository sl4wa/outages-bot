package outageapi

import (
	"os"
	"strings"
)

type cacheEntry struct {
	etag string
	body []byte
}

// loadCache returns nil if the file does not exist, cannot be parsed, or has an empty body.
func loadCache(path string) *cacheEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	parts := strings.SplitN(string(data), "\n", 2)
	if len(parts) < 2 {
		return nil
	}
	etag := parts[0]
	body := []byte(parts[1])
	if len(body) == 0 {
		return nil
	}
	return &cacheEntry{etag: etag, body: body}
}

// saveCache writes etag and body to path atomically using a temp-file rename.
func saveCache(path string, etag string, body []byte) error {
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(etag+"\n"+string(body)), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return nil
}
