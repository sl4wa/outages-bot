package loe

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderParsesHydraMemberPayload(t *testing.T) {
	payload := `{"hydra:member":[{"menuItems":[{"children":[{"id":1023,"rawHtml":"<div><p>Графік погодинних відключень на 16.03.2026</p><p>Інформація станом на 18:25 15.03.2026</p><p>Група 5.2. Електроенергії немає з 17:00 до 19:30.</p></div>"}]}]}]}`

	result, err := Provider{LoadPayload: func(context.Context) (string, error) { return payload, nil }}.GetSchedules(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotNil(t, result[0].SourceID)
	assert.Equal(t, 1023, *result[0].SourceID)
	assert.Equal(t, time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC), result[0].ScheduleDate)
	require.NotNil(t, result[0].UpdatedAt)
	assert.Equal(t, time.Date(2026, 3, 15, 18, 25, 0, 0, time.UTC), *result[0].UpdatedAt)
}

func TestProviderParsesJSONArrayFormat(t *testing.T) {
	payload := `[{"menuItems":[{"rawHtml":"<div><p>Графік погодинних відключень на 16.03.2026</p><p>Інформація станом на 18:25 15.03.2026</p><p>Група 5.2. Електроенергії немає з 17:00 до 19:30.</p></div>","children":[]}]}]`

	result, err := Provider{LoadPayload: func(context.Context) (string, error) { return payload, nil }}.GetSchedules(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC), result[0].ScheduleDate)
	require.NotNil(t, result[0].UpdatedAt)
	assert.Equal(t, time.Date(2026, 3, 15, 18, 25, 0, 0, time.UTC), *result[0].UpdatedAt)
}

func TestProviderParsesDescriptionWhenRawHTMLIsBlank(t *testing.T) {
	payload := `{"hydra:member":[{"menuItems":[{"id":2048,"rawHtml":"","description":"Отримано вказівку НЕК «Укренерго» про відміну застосування ГПВ на 18.05.2026 до окремого розпорядження НЕК «Укренерго».\nУ випадку нового розпорядження НЕК «Укренерго» про застосування погодинних відключень, графік буде опубліковано додатково.\nІнформація станом на 19:45 17.05.2026"}]}]}`

	result, err := Provider{LoadPayload: func(context.Context) (string, error) { return payload, nil }}.GetSchedules(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotNil(t, result[0].SourceID)
	assert.Equal(t, 2048, *result[0].SourceID)
	assert.Equal(t, time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC), result[0].ScheduleDate)
	require.NotNil(t, result[0].UpdatedAt)
	assert.Equal(t, time.Date(2026, 5, 17, 19, 45, 0, 0, time.UTC), *result[0].UpdatedAt)
	assert.Contains(t, result[0].Text, "відміну застосування ГПВ на 18.05.2026")
}

func TestProviderPrefersRawHTMLOverDescription(t *testing.T) {
	payload := `{"hydra:member":[{"menuItems":[{"id":2049,"rawHtml":"<div><p>Графік погодинних відключень на 19.05.2026</p><p>Інформація станом на 08:15 18.05.2026</p><p>Група 1.1. Електроенергія є.</p></div>","description":"Отримано вказівку НЕК «Укренерго» про відміну застосування ГПВ на 20.05.2026.\nІнформація станом на 20:30 19.05.2026"}]}]}`

	result, err := Provider{LoadPayload: func(context.Context) (string, error) { return payload, nil }}.GetSchedules(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC), result[0].ScheduleDate)
	require.NotNil(t, result[0].UpdatedAt)
	assert.Equal(t, time.Date(2026, 5, 18, 8, 15, 0, 0, time.UTC), *result[0].UpdatedAt)
	assert.Contains(t, result[0].Text, "Група 1.1. Електроенергія є.")
	assert.NotContains(t, result[0].Text, "20.05.2026")
}

func TestProviderInvalidPayloads(t *testing.T) {
	_, err := Provider{LoadPayload: func(context.Context) (string, error) { return "", nil }}.GetSchedules(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "blank")

	_, err = Provider{LoadPayload: func(context.Context) (string, error) { return "{bad", nil }}.GetSchedules(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not valid JSON")
}

func TestProviderCallsAcceptedWhenParseYieldsEmpty(t *testing.T) {
	payload := `{"hydra:member":[]}`
	called := false

	result, err := Provider{
		LoadPayload: func(context.Context) (string, error) { return payload, nil },
		OnPayloadAccepted: func() error {
			called = true
			return nil
		},
	}.GetSchedules(context.Background())

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.True(t, called)
}

func TestProviderCallsAcceptedAfterValidParse(t *testing.T) {
	payload := `[{"menuItems":[{"rawHtml":"<div><p>Графік погодинних відключень на 16.03.2026</p><p>Інформація станом на 18:25 15.03.2026</p><p>Група 1.1. Електроенергія є.</p></div>"}]}]`
	called := false

	_, err := Provider{
		LoadPayload: func(context.Context) (string, error) { return payload, nil },
		OnPayloadAccepted: func() error {
			called = true
			return nil
		},
	}.GetSchedules(context.Background())

	require.NoError(t, err)
	assert.True(t, called)
}
