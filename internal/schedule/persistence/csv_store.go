package persistence

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/sl4wa/outages-bot/internal/schedule/schedule"
)

const StateFileName = "schedule.csv"

type CSVStateStore struct {
	Path string
}

func NewCSVStateStore(path string) CSVStateStore {
	return CSVStateStore{Path: path}
}

func (s CSVStateStore) Load() (map[time.Time]string, error) {
	result := make(map[time.Time]string)
	file, err := os.Open(s.Path)
	if errors.Is(err, os.ErrNotExist) {
		return result, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("WARNING: skipping malformed row in %s: %v", s.Path, err)
			continue
		}
		if len(record) < 2 {
			log.Printf("WARNING: skipping short row in %s: %v", s.Path, record)
			continue
		}
		date, err := schedule.ParseStateDate(record[0])
		if err != nil {
			log.Printf("WARNING: skipping row in %s with unparseable date %q: %v", s.Path, record[0], err)
			continue
		}
		result[date] = record[1]
	}
	return result, nil
}

func (s CSVStateStore) Save(state map[time.Time]string) error {
	if dir := filepath.Dir(s.Path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	tmpPath := s.Path + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)

	dates := make([]time.Time, 0, len(state))
	for date := range state {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })
	for _, date := range dates {
		if err := writer.Write([]string{schedule.FormatStateDate(date), state[date]}); err != nil {
			file.Close()
			os.Remove(tmpPath)
			return err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		file.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := file.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, s.Path); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return nil
}
