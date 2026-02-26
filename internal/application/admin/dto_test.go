package admin

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPeriodFormatter_SameDay(t *testing.T) {
	start := time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC)
	end := time.Date(2024, 3, 15, 16, 30, 0, 0, time.UTC)

	result := PeriodFormatter(start, end)
	assert.Equal(t, "15.03.2024 08:00 - 16:30", result)
}

func TestPeriodFormatter_MultiDay(t *testing.T) {
	start := time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC)
	end := time.Date(2024, 3, 16, 16, 0, 0, 0, time.UTC)

	result := PeriodFormatter(start, end)
	assert.Equal(t, "15.03.2024 08:00 - 16.03.2024 16:00", result)
}

func TestPeriodFormatter_MidnightBoundary(t *testing.T) {
	start := time.Date(2024, 3, 15, 23, 0, 0, 0, time.UTC)
	end := time.Date(2024, 3, 16, 1, 0, 0, 0, time.UTC)

	result := PeriodFormatter(start, end)
	assert.Equal(t, "15.03.2024 23:00 - 16.03.2024 01:00", result)
}

func TestPeriodFormatter_SameStartAndEnd(t *testing.T) {
	ts := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)

	result := PeriodFormatter(ts, ts)
	assert.Equal(t, "15.03.2024 12:00 - 12:00", result)
}
