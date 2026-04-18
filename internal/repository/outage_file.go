package repository

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"outages-bot/internal/domain"
)

const OutageSnapshotFileName = "outages.csv"

// FileOutageRepository persists outage data for deduplication as a CSV file.
type FileOutageRepository struct {
	path string
}

// NewFileOutageRepository creates a FileOutageRepository that stores the
// outage data at path.
func NewFileOutageRepository(path string) *FileOutageRepository {
	return &FileOutageRepository{path: path}
}

// Load reads the last saved outage data. Returns (nil, nil) when no data exists yet.
func (r *FileOutageRepository) Load() ([]*domain.Outage, error) {
	data, err := os.ReadFile(r.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read outage data: %w", err)
	}

	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse outage data: %w", err)
	}

	// skip header row
	if len(records) == 0 {
		return nil, nil
	}
	records = records[1:]

	outages := make([]*domain.Outage, 0, len(records))
	for _, row := range records {
		if len(row) != 7 {
			return nil, fmt.Errorf("unexpected column count %d", len(row))
		}
		start, err := time.Parse(time.RFC3339, row[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse start time: %w", err)
		}
		end, err := time.Parse(time.RFC3339, row[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse end time: %w", err)
		}
		city := row[2]
		streetID, err := strconv.Atoi(row[3])
		if err != nil {
			return nil, fmt.Errorf("failed to parse street_id: %w", err)
		}
		streetName := row[4]
		var buildings []string
		if row[5] != "" {
			buildings = strings.Split(row[5], "|")
		}
		comment := row[6]

		period, err := domain.NewOutagePeriod(start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to parse outage period: %w", err)
		}
		addr, err := domain.NewOutageAddress(streetID, streetName, buildings, city)
		if err != nil {
			return nil, fmt.Errorf("failed to parse outage address: %w", err)
		}
		outages = append(outages, &domain.Outage{
			Period:      period,
			Address:     addr,
			Description: domain.NewOutageDescription(comment),
		})
	}
	return outages, nil
}

// Save writes outages to disk using atomic write (temp file + rename).
func (r *FileOutageRepository) Save(outages []*domain.Outage) error {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write([]string{"start", "end", "city", "street_id", "street_name", "buildings", "comment"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}
	for _, o := range outages {
		row := []string{
			o.Period.StartDate.UTC().Format(time.RFC3339),
			o.Period.EndDate.UTC().Format(time.RFC3339),
			o.Address.City,
			strconv.Itoa(o.Address.StreetID),
			o.Address.StreetName,
			strings.Join(o.Address.Buildings, "|"),
			o.Description.Value,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to marshal outage data: %w", err)
	}

	tmpPath := r.path + ".tmp"
	if err := os.WriteFile(tmpPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write temp outage data: %w", err)
	}
	if err := os.Rename(tmpPath, r.path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename outage data file: %w", err)
	}
	return nil
}
