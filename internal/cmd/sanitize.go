package cli

import (
	"regexp"
	"strings"
)

// sanitizeRegex matches Unicode format characters (\p{Cf}) and Hangul filler (U+3164).
var sanitizeRegex = regexp.MustCompile(`[\p{Cf}\x{3164}]`)

// Sanitize removes invisible/control Unicode characters from a string.
// Returns "-" for empty or whitespace-only results.
func Sanitize(value string) string {
	if value == "" {
		return "-"
	}

	cleaned := sanitizeRegex.ReplaceAllString(value, "")
	trimmed := strings.TrimSpace(cleaned)

	if trimmed == "" {
		return "-"
	}

	return trimmed
}
