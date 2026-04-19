package notifier

import "errors"

// ErrRecipientUnavailable indicates the recipient can no longer receive messages.
var ErrRecipientUnavailable = errors.New("recipient unavailable")
