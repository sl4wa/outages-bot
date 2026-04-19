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

	"gopkg.in/yaml.v3"
)

type userFile struct {
	StreetID   int    `yaml:"street_id"`
	StreetName string `yaml:"street_name"`
	Building   string `yaml:"building"`
	StartDate  string `yaml:"start_date,omitempty"`
	EndDate    string `yaml:"end_date,omitempty"`
	Comment    string `yaml:"comment,omitempty"`
}

// FileUserRepository persists users as individual YAML files.
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
	entries, _ := filepath.Glob(filepath.Join(r.dataDir, "*.yml"))

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

// Find returns a user by chat ID, or (nil, nil) if not found.
func (r *FileUserRepository) Find(chatID int64) (*users.User, error) {
	path := r.filePath(chatID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}
	return r.loadFromFile(path)
}

// Save persists a user to disk using atomic write (temp file + rename).
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

	content, err := yaml.Marshal(&uf)
	if err != nil {
		return fmt.Errorf("failed to marshal user file: %w", err)
	}

	// Atomic write: temp file + rename
	tmpPath := r.filePath(user.ID) + ".tmp"
	if err := os.WriteFile(tmpPath, content, 0o644); err != nil {
		return fmt.Errorf("failed to write temp user file: %w", err)
	}
	if err := os.Rename(tmpPath, r.filePath(user.ID)); err != nil {
		os.Remove(tmpPath) // cleanup on rename failure
		return fmt.Errorf("failed to rename user file: %w", err)
	}
	return nil
}

// Remove deletes a user file. Returns (false, nil) if not found.
func (r *FileUserRepository) Remove(chatID int64) (bool, error) {
	path := r.filePath(chatID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	if err := os.Remove(path); err != nil {
		return false, fmt.Errorf("failed to remove user file: %w", err)
	}
	return true, nil
}

func (r *FileUserRepository) filePath(chatID int64) string {
	return filepath.Join(r.dataDir, fmt.Sprintf("%d.yml", chatID))
}

func (r *FileUserRepository) loadFromFile(path string) (*users.User, error) {
	base := filepath.Base(path)
	idStr := strings.TrimSuffix(base, ".yml")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user file name %s: %w", base, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read user file: %w", err)
	}

	var uf userFile
	if err := yaml.Unmarshal(data, &uf); err != nil {
		return nil, fmt.Errorf("failed to parse user file %s: %w", base, err)
	}

	addr, err := users.NewAddress(uf.StreetID, uf.StreetName, uf.Building)
	if err != nil {
		return nil, fmt.Errorf("invalid user address in %s: %w", base, err)
	}

	var outageInfo *users.OutageInfo
	if uf.StartDate != "" && uf.EndDate != "" {
		startDate, err := time.Parse(time.RFC3339, uf.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date in %s: %w", base, err)
		}
		endDate, err := time.Parse(time.RFC3339, uf.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date in %s: %w", base, err)
		}
		period, err := outage.NewPeriod(startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("invalid outage period in %s: %w", base, err)
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
