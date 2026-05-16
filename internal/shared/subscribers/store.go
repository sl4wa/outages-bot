package subscribers

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const fileExt = ".toml"

type FileStore struct {
	Dir string
}

func NewFileStore(dir string) FileStore {
	return FileStore{Dir: dir}
}

func (s FileStore) FilePath(chatID int64) string {
	return filepath.Join(s.Dir, strconv.FormatInt(chatID, 10)+fileExt)
}

func (s FileStore) ChatIDs() ([]int64, error) {
	entries, err := os.ReadDir(s.Dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read subscribers dir: %w", err)
	}
	var result []int64
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != fileExt {
			continue
		}
		stem := strings.TrimSuffix(name, fileExt)
		id, err := strconv.ParseInt(stem, 10, 64)
		if err != nil || strconv.FormatInt(id, 10) != stem {
			continue
		}
		result = append(result, id)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result, nil
}

func (s FileStore) Read(chatID int64) ([]byte, error) {
	return os.ReadFile(s.FilePath(chatID))
}

func (s FileStore) Write(chatID int64, data []byte) error {
	if err := os.MkdirAll(s.Dir, 0o770); err != nil {
		return fmt.Errorf("create subscribers dir: %w", err)
	}
	path := s.FilePath(chatID)
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return fmt.Errorf("write temp subscriber file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename subscriber file: %w", err)
	}
	return nil
}

func (s FileStore) Remove(chatID int64) (bool, error) {
	if err := os.Remove(s.FilePath(chatID)); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("remove subscriber file: %w", err)
	}
}
