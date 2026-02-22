package repository

import (
	"encoding/csv"
	"fmt"
	"os"
	"outages-bot/internal/domain"
	"strconv"
)

// FileStreetRepository reads streets from a CSV file.
type FileStreetRepository struct {
	streets []domain.Street
}

// NewFileStreetRepository creates a FileStreetRepository by loading streets from the given file path.
func NewFileStreetRepository(filePath string) (*FileStreetRepository, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open streets file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse streets file: %w", err)
	}

	if len(records) < 1 {
		return &FileStreetRepository{}, nil
	}

	streets := make([]domain.Street, 0, len(records)-1)
	for _, record := range records[1:] {
		if len(record) < 2 {
			continue
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("invalid street id %q: %w", record[0], err)
		}
		streets = append(streets, domain.Street{ID: id, Name: record[1]})
	}

	return &FileStreetRepository{streets: streets}, nil
}

// GetAllStreets returns all loaded streets.
func (r *FileStreetRepository) GetAllStreets() ([]domain.Street, error) {
	return r.streets, nil
}
