package outageapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"outages-bot/internal/application"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const defaultAPIURL = "https://power-api.loe.lviv.ua/api/pw_accidents?pagination=false&otg.id=28&city.id=693"

var newlineRegex = regexp.MustCompile(`[\r\n]+`)

// Provider fetches outages from the Lviv power outage API.
type Provider struct {
	baseURL string
	client  *http.Client
	clock   func() time.Time
	logger  *log.Logger
}

// NewProvider creates a new Provider.
func NewProvider(baseURL string, clock func() time.Time, logger *log.Logger) *Provider {
	if baseURL == "" {
		baseURL = defaultAPIURL
	}
	if clock == nil {
		clock = time.Now
	}
	if logger == nil {
		logger = log.Default()
	}
	return &Provider{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
		clock:   clock,
		logger:  logger,
	}
}

type apiResponse struct {
	HydraMember []json.RawMessage `json:"hydra:member"`
}

type apiRow struct {
	ID            interface{}     `json:"id"`
	DateEvent     string          `json:"dateEvent"`
	DatePlanIn    string          `json:"datePlanIn"`
	Koment        string          `json:"koment"`
	BuildingNames json.RawMessage `json:"buildingNames"`
	City          json.RawMessage `json:"city"`
	Street        json.RawMessage `json:"street"`
}

type cityObj struct {
	Name string `json:"name"`
}

type streetObj struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

// FetchOutages fetches outages from the API.
func (p *Provider) FetchOutages(ctx context.Context) ([]application.OutageDTO, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch outages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p.logger.Printf("WARNING: outage API returned status %d", resp.StatusCode)
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	var outages []application.OutageDTO
	seen := make(map[string]int) // key → index in outages

	for _, raw := range apiResp.HydraMember {
		var row apiRow
		if err := json.Unmarshal(raw, &row); err != nil {
			continue
		}

		id := toInt(row.ID)

		comment := newlineRegex.ReplaceAllString(row.Koment, " ")
		comment = strings.TrimSpace(comment)

		buildings := p.parseBuildings(row.BuildingNames)

		var city cityObj
		if row.City != nil {
			json.Unmarshal(row.City, &city)
		}

		var street streetObj
		if row.Street != nil {
			json.Unmarshal(row.Street, &street)
		}

		streetID := toInt(street.ID)

		start := p.parseDate(row.DateEvent)
		end := p.parseDate(row.DatePlanIn)

		key := fmt.Sprintf("%d|%s|%d|%d", streetID, strings.Join(buildings, ","), start.Unix(), end.Unix())

		dto := application.OutageDTO{
			ID:         id,
			Start:      start,
			End:        end,
			City:       city.Name,
			StreetID:   streetID,
			StreetName: street.Name,
			Buildings:  buildings,
			Comment:    comment,
		}

		if idx, ok := seen[key]; ok {
			outages[idx] = dto
		} else {
			seen[key] = len(outages)
			outages = append(outages, dto)
		}
	}

	return outages, nil
}

func (p *Provider) parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return p.clock()
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Try other common formats
		t, err = time.Parse("2006-01-02T15:04:05", dateStr)
		if err != nil {
			return p.clock()
		}
	}
	return t
}

func (p *Provider) parseBuildings(raw json.RawMessage) []string {
	if raw == nil {
		return nil
	}

	// Try as array first
	var arr []interface{}
	if err := json.Unmarshal(raw, &arr); err == nil {
		var buildings []string
		for _, item := range arr {
			s := fmt.Sprintf("%v", item)
			s = strings.TrimSpace(s)
			if s != "" {
				buildings = append(buildings, s)
			}
		}
		return buildings
	}

	// Try as string
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		var buildings []string
		for _, part := range strings.Split(s, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				buildings = append(buildings, part)
			}
		}
		return buildings
	}

	return nil
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case json.Number:
		n, err := val.Int64()
		if err != nil {
			f, ferr := val.Float64()
			if ferr != nil || math.IsNaN(f) || math.IsInf(f, 0) {
				return 0
			}
			return int(f)
		}
		return int(n)
	case string:
		val = strings.TrimSpace(val)
		n, err := strconv.Atoi(val)
		if err != nil {
			f, ferr := strconv.ParseFloat(val, 64)
			if ferr != nil || math.IsNaN(f) || math.IsInf(f, 0) {
				return 0
			}
			return int(f)
		}
		return n
	default:
		return 0
	}
}
