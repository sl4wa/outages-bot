package persistence

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sl4wa/outages-bot/internal/outage/outage"
	"github.com/sl4wa/outages-bot/internal/outage/users"
	"github.com/sl4wa/outages-bot/internal/shared/subscribers"

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
	store subscribers.FileStore
}

// NewFileUserRepository creates a new FileUserRepository with the given data directory.
func NewFileUserRepository(dataDir string) (*FileUserRepository, error) {
	if err := os.MkdirAll(dataDir, 0o770); err != nil {
		return nil, fmt.Errorf("failed to create user data directory: %w", err)
	}

	return &FileUserRepository{store: subscribers.NewFileStore(dataDir)}, nil
}

// FindAll returns all users from disk.
func (r *FileUserRepository) FindAll() []*users.User {
	ids, err := r.store.ChatIDs()
	if err != nil {
		log.Printf("WARNING: failed to list user files: %v", err)
		return nil
	}
	usersList := make([]*users.User, 0, len(ids))
	for _, id := range ids {
		user, err := r.load(id)
		if err != nil {
			log.Printf("WARNING: skipping malformed user file %d.toml: %v", id, err)
			continue
		}
		usersList = append(usersList, user)
	}
	return usersList
}

// Find returns a user by chat ID, or nil if not found.
func (r *FileUserRepository) Find(chatID int64) (*users.User, error) {
	user, err := r.load(chatID)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
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

	content, err := toml.Marshal(&uf)
	if err != nil {
		return fmt.Errorf("failed to marshal user file: %w", err)
	}
	if err := r.store.Write(user.ID, content); err != nil {
		return fmt.Errorf("failed to save user file: %w", err)
	}
	return nil
}

// Remove deletes the user's .toml file. Returns (false, nil) if not found.
func (r *FileUserRepository) Remove(chatID int64) (bool, error) {
	return r.store.Remove(chatID)
}

func (r *FileUserRepository) load(chatID int64) (*users.User, error) {
	data, err := r.store.Read(chatID)
	if err != nil {
		return nil, err
	}

	var uf userFile
	if err := toml.Unmarshal(data, &uf); err != nil {
		return nil, fmt.Errorf("failed to parse user file %d.toml: %w", chatID, err)
	}

	return decodeUserFile(uf, chatID)
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
