package telegram

import (
	"errors"
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeErrorNil(t *testing.T) {
	require.NoError(t, NormalizeError(nil))
}

func TestNormalizeErrorForbiddenIsRecipientUnavailable(t *testing.T) {
	apiErr := &tgbotapi.Error{Code: 403, Message: "Forbidden: bot was blocked by the user"}

	err := NormalizeError(apiErr)

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrRecipientUnavailable), "expected ErrRecipientUnavailable in chain")

	var unwrapped *tgbotapi.Error
	require.True(t, errors.As(err, &unwrapped), "expected *tgbotapi.Error in chain")
	assert.Equal(t, 403, unwrapped.Code)
}

func TestNormalizeErrorOtherAPIErrorIsWrapped(t *testing.T) {
	apiErr := &tgbotapi.Error{Code: 400, Message: "bad request"}

	err := NormalizeError(apiErr)

	require.Error(t, err)
	assert.False(t, errors.Is(err, ErrRecipientUnavailable))

	var unwrapped *tgbotapi.Error
	require.True(t, errors.As(err, &unwrapped))
	assert.Equal(t, 400, unwrapped.Code)
}

func TestNormalizeErrorGenericErrorIsWrapped(t *testing.T) {
	err := NormalizeError(errors.New("network down"))

	require.Error(t, err)
	assert.False(t, errors.Is(err, ErrRecipientUnavailable))
	assert.True(t, strings.HasPrefix(err.Error(), "telegram send:"), "got: %q", err.Error())
}
