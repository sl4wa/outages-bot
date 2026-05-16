package persistence

import (
	"fmt"
	"log"
	"os"
	"outages-bot/internal/outage"
	"outages-bot/internal/users"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type userFile struct {
	StreetID   int    `toml:"street_id"`
	StreetName string `toml:"street_name"`
	Building   string `toml:"building"`
	StartDate  string `toml:"start_date,omitempty"`
	EndDate    string `toml:"end_date,omitempty"`
	Comment    string `toml:"comment,omitempty"`
}

// FileUserRepository persists users as individual TOML files.
type FileUserRepository struct {
	dataDir string
}

// NewFileUserRepository creates a new FileUserRepository with the given data directory.
func NewFileUserRepository(dataDir string) (*FileUserRepository, error) {
	if err := os.MkdirAll(dataDir, 0o770); err != nil {
		return nil, fmt.Errorf("failed to create user data directory: %w", err)
	}

	return &FileUserRepository{dataDir: dataDir}, nil
}

// FindAll returns all users from disk.
func (r *FileUserRepository) FindAll() []*users.User {
	entries, _ := filepath.Glob(filepath.Join(r.dataDir, "*.toml"))
	usersList := make([]*users.User, 0, len(entries))
	for _, path := range entries {
		user, err := r.loadFromFile(path)
		if err != nil {
			log.Printf("WARNING: skipping malformed user file %s: %v", filepath.Base(path), err)
			continue
		}
		usersList = append(usersList, user)
	}
	return usersList
}

// Find returns a user by chat ID, or nil if not found.
func (r *FileUserRepository) Find(chatID int64) (*users.User, error) {
	tomlPath := r.filePath(chatID)
	if _, err := os.Stat(tomlPath); err == nil {
		return r.loadFromFile(tomlPath)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat user file: %w", err)
	}
	return nil, nil
}

// Save persists a user to disk as TOML using atomic write (temp file + rename).
func (r *FileUserRepository) Save(user *users.User) error {
	uf := userFile{
		StreetID:   user.Address.StreetID,
		StreetName: user.Address.StreetName,
		Building:   user.Address.Building,
	}

	if user.OutageInfo != nil {
		uf.StartDate = user.OutageInfo.Period.StartDate.Format(time.RFC3339)
		uf.EndDate = user.OutageInfo.Period.EndDate.Format(time.RFC3339)
		uf.Comment = user.OutageInfo.Description.Value
	}

	return r.writeTomlFile(user.ID, uf)
}

// Remove deletes the user's .toml file. Returns (false, nil) if not found.
func (r *FileUserRepository) Remove(chatID int64) (bool, error) {
	if err := os.Remove(r.filePath(chatID)); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to remove user file: %w", err)
	}
}

func (r *FileUserRepository) filePath(chatID int64) string {
	return filepath.Join(r.dataDir, fmt.Sprintf("%d.toml", chatID))
}

func (r *FileUserRepository) writeTomlFile(chatID int64, uf userFile) error {
	content, err := toml.Marshal(&uf)
	if err != nil {
		return fmt.Errorf("failed to marshal user file: %w", err)
	}

	tomlPath := r.filePath(chatID)
	tmpPath := tomlPath + ".tmp"
	if err := os.WriteFile(tmpPath, content, 0o644); err != nil {
		return fmt.Errorf("failed to write temp user file: %w", err)
	}
	if err := os.Rename(tmpPath, tomlPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename user file: %w", err)
	}
	return nil
}

func (r *FileUserRepository) loadFromFile(path string) (*users.User, error) {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	idStr := strings.TrimSuffix(base, ext)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user file name %s: %w", base, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read user file: %w", err)
	}

	var uf userFile
	if err := toml.Unmarshal(data, &uf); err != nil {
		return nil, fmt.Errorf("failed to parse user file %s: %w", base, err)
	}

	return decodeUserFile(uf, id)
}

func decodeUserFile(uf userFile, id int64) (*users.User, error) {
	addr, err := users.NewAddress(uf.StreetID, uf.StreetName, uf.Building)
	if err != nil {
		return nil, fmt.Errorf("invalid user address in %d: %w", id, err)
	}

	var outageInfo *users.OutageInfo
	if uf.StartDate != "" && uf.EndDate != "" {
		startDate, err := time.Parse(time.RFC3339, uf.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date in %d: %w", id, err)
		}
		endDate, err := time.Parse(time.RFC3339, uf.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date in %d: %w", id, err)
		}
		period, err := outage.NewPeriod(startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("invalid outage period in %d: %w", id, err)
		}
		desc := outage.NewDescription(uf.Comment)
		info := users.NewOutageInfo(period, desc)
		outageInfo = &info
	}

	return &users.User{
		ID:         id,
		Address:    addr,
		OutageInfo: outageInfo,
	}, nil
}
