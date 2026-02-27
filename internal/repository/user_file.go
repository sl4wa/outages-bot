package repository

import (
	"fmt"
	"log"
	"os"
	"outages-bot/internal/domain"
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

// FileUserRepository persists users as individual text files.
type FileUserRepository struct {
	dataDir string
}

// NewFileUserRepository creates a new FileUserRepository with the given data directory.
func NewFileUserRepository(dataDir string) (*FileUserRepository, error) {
	if err := os.MkdirAll(dataDir, 0o770); err != nil {
		return nil, fmt.Errorf("failed to create user data directory: %w", err)
	}

	// TODO: Remove migration after all environments have been migrated
	txtFiles, _ := filepath.Glob(filepath.Join(dataDir, "*.txt"))
	for _, txtPath := range txtFiles {
		data, err := os.ReadFile(txtPath)
		if err != nil {
			log.Printf("WARNING: migration: failed to read %s: %v", filepath.Base(txtPath), err)
			continue
		}
		uf, err := parseLegacy(data)
		if err != nil {
			log.Printf("WARNING: migration: failed to parse %s: %v", filepath.Base(txtPath), err)
			continue
		}
		ymlData, err := yaml.Marshal(&uf)
		if err != nil {
			log.Printf("WARNING: migration: failed to marshal %s: %v", filepath.Base(txtPath), err)
			continue
		}
		ymlPath := strings.TrimSuffix(txtPath, ".txt") + ".yml"
		if _, err := os.Stat(ymlPath); err == nil {
			log.Printf("WARNING: migration: skipping %s because %s already exists", filepath.Base(txtPath), filepath.Base(ymlPath))
			continue
		}
		if err := os.WriteFile(ymlPath, ymlData, 0o644); err != nil {
			log.Printf("WARNING: migration: failed to write %s: %v", filepath.Base(ymlPath), err)
			continue
		}
		if err := os.Remove(txtPath); err != nil {
			log.Printf("WARNING: migration: failed to delete %s: %v", filepath.Base(txtPath), err)
			continue
		}
		log.Printf("Migrated %s -> %s", filepath.Base(txtPath), filepath.Base(ymlPath))
	}

	return &FileUserRepository{dataDir: dataDir}, nil
}

// FindAll returns all users from disk.
func (r *FileUserRepository) FindAll() []*domain.User {
	entries, _ := filepath.Glob(filepath.Join(r.dataDir, "*.yml"))

	users := make([]*domain.User, 0, len(entries))
	for _, path := range entries {
		user, err := r.loadFromFile(path)
		if err != nil {
			log.Printf("WARNING: skipping malformed user file %s: %v", filepath.Base(path), err)
			continue
		}
		users = append(users, user)
	}
	return users
}

// Find returns a user by chat ID, or (nil, nil) if not found.
func (r *FileUserRepository) Find(chatID int64) (*domain.User, error) {
	path := r.filePath(chatID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}
	return r.loadFromFile(path)
}

// Save persists a user to disk using atomic write (temp file + rename).
func (r *FileUserRepository) Save(user *domain.User) error {
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

func (r *FileUserRepository) loadFromFile(path string) (*domain.User, error) {
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

	addr, err := domain.NewUserAddress(uf.StreetID, uf.StreetName, uf.Building)
	if err != nil {
		return nil, fmt.Errorf("invalid user address in %s: %w", base, err)
	}

	var outageInfo *domain.OutageInfo
	if uf.StartDate != "" && uf.EndDate != "" {
		startDate, err := time.Parse(time.RFC3339, uf.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date in %s: %w", base, err)
		}
		endDate, err := time.Parse(time.RFC3339, uf.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date in %s: %w", base, err)
		}
		period, err := domain.NewOutagePeriod(startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("invalid outage period in %s: %w", base, err)
		}
		desc := domain.NewOutageDescription(uf.Comment)
		info := domain.NewOutageInfo(period, desc)
		outageInfo = &info
	}

	return &domain.User{
		ID:         id,
		Address:    addr,
		OutageInfo: outageInfo,
	}, nil
}

// parseLegacy parses the old line-based key: value format used in .txt user files.
func parseLegacy(data []byte) (userFile, error) {
	var uf userFile
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, ": ")
		if idx < 0 {
			continue
		}
		key := line[:idx]
		value := line[idx+2:]
		switch key {
		case "street_id":
			id, err := strconv.Atoi(value)
			if err != nil {
				return uf, fmt.Errorf("invalid street_id: %w", err)
			}
			uf.StreetID = id
		case "street_name":
			uf.StreetName = value
		case "building":
			uf.Building = value
		case "start_date":
			uf.StartDate = value
		case "end_date":
			uf.EndDate = value
		case "comment":
			uf.Comment = value
		}
	}
	if uf.StreetID <= 0 || uf.StreetName == "" || uf.Building == "" {
		return uf, fmt.Errorf("missing required fields (street_id, street_name, building)")
	}
	return uf, nil
}
