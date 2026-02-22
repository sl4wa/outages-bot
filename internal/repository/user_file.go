package repository

import (
	"fmt"
	"os"
	"outages-bot/internal/domain"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FileUserRepository persists users as individual text files.
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
func (r *FileUserRepository) FindAll() ([]*domain.User, error) {
	entries, err := filepath.Glob(filepath.Join(r.dataDir, "*.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to list user files: %w", err)
	}

	users := make([]*domain.User, 0, len(entries))
	for _, path := range entries {
		user, err := r.loadFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load user from %s: %w", path, err)
		}
		users = append(users, user)
	}
	return users, nil
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
	fields := map[string]string{
		"street_id":   strconv.Itoa(user.Address.StreetID),
		"street_name": user.Address.StreetName,
		"building":    user.Address.Building,
		"start_date":  "",
		"end_date":    "",
		"comment":     "",
	}

	if user.OutageInfo != nil {
		fields["start_date"] = user.OutageInfo.Period.StartDate.Format(time.RFC3339)
		fields["end_date"] = user.OutageInfo.Period.EndDate.Format(time.RFC3339)
		fields["comment"] = user.OutageInfo.Description.Value
	}

	// Preserve field order
	order := []string{"street_id", "street_name", "building", "start_date", "end_date", "comment"}
	var lines []string
	for _, key := range order {
		lines = append(lines, key+": "+fields[key])
	}
	content := strings.Join(lines, "\n")

	// Atomic write: temp file + rename
	tmpPath := r.filePath(user.ID) + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(content), 0o644); err != nil {
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
	return filepath.Join(r.dataDir, fmt.Sprintf("%d.txt", chatID))
}

func (r *FileUserRepository) loadFromFile(path string) (*domain.User, error) {
	base := filepath.Base(path)
	idStr := strings.TrimSuffix(base, ".txt")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user file name %s: %w", base, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read user file: %w", err)
	}

	fields := map[string]string{
		"street_id":   "0",
		"street_name": "",
		"building":    "",
		"start_date":  "",
		"end_date":    "",
		"comment":     "",
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		if _, ok := fields[key]; ok {
			fields[key] = val
		}
	}

	streetID, _ := strconv.Atoi(fields["street_id"])
	addr, err := domain.NewUserAddress(streetID, fields["street_name"], fields["building"])
	if err != nil {
		return nil, fmt.Errorf("invalid user address in %s: %w", base, err)
	}

	var outageInfo *domain.OutageInfo
	if fields["start_date"] != "" && fields["end_date"] != "" {
		startDate, err := time.Parse(time.RFC3339, fields["start_date"])
		if err != nil {
			return nil, fmt.Errorf("invalid start_date in %s: %w", base, err)
		}
		endDate, err := time.Parse(time.RFC3339, fields["end_date"])
		if err != nil {
			return nil, fmt.Errorf("invalid end_date in %s: %w", base, err)
		}
		period, err := domain.NewOutagePeriod(startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("invalid outage period in %s: %w", base, err)
		}
		desc := domain.NewOutageDescription(fields["comment"])
		info := domain.NewOutageInfo(period, desc)
		outageInfo = &info
	}

	return &domain.User{
		ID:         id,
		Address:    addr,
		OutageInfo: outageInfo,
	}, nil
}
