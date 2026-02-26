package cli

import (
	"bytes"
	"context"
	"errors"
	"outages-bot/internal/application"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockOutageProviderForOutages struct {
	outages []application.OutageDTO
	err     error
}

func (m *mockOutageProviderForOutages) FetchOutages(_ context.Context) ([]application.OutageDTO, error) {
	return m.outages, m.err
}

func TestRunOutagesCommand_PrintsTable(t *testing.T) {
	provider := &mockOutageProviderForOutages{
		outages: []application.OutageDTO{
			{
				ID:         1,
				StreetID:   100,
				StreetName: "Стрийська",
				Buildings:  []string{"10", "12"},
				Start:      time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC),
				End:        time.Date(2024, 3, 15, 16, 0, 0, 0, time.UTC),
				Comment:    "Ремонт",
			},
		},
	}

	var buf bytes.Buffer
	err := RunOutagesCommand(context.Background(), provider, &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "StreetID")
	assert.Contains(t, output, "Street")
	assert.Contains(t, output, "Buildings")
	assert.Contains(t, output, "Period")
	assert.Contains(t, output, "Comment")
	assert.Contains(t, output, "100")
	assert.Contains(t, output, "Стрийська")
	assert.Contains(t, output, "10, 12")
	assert.Contains(t, output, "15.03.2024 08:00 - 16:00")
	assert.Contains(t, output, "Ремонт")
}

func TestRunOutagesCommand_Empty(t *testing.T) {
	provider := &mockOutageProviderForOutages{outages: []application.OutageDTO{}}

	var buf bytes.Buffer
	err := RunOutagesCommand(context.Background(), provider, &buf)
	require.NoError(t, err)
	assert.Equal(t, "No outages found.\n", buf.String())
}

func TestRunOutagesCommand_ProviderError(t *testing.T) {
	provider := &mockOutageProviderForOutages{err: errors.New("timeout")}

	var buf bytes.Buffer
	err := RunOutagesCommand(context.Background(), provider, &buf)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch outages")
	assert.Contains(t, err.Error(), "timeout")
}
