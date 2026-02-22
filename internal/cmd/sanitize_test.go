package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal ASCII text", "Hello World", "Hello World"},
		{"zero-width joiner", "Hello\u200DWorld", "HelloWorld"},
		{"zero-width non-joiner", "Hello\u200CWorld", "HelloWorld"},
		{"soft hyphen", "Hello\u00ADWorld", "HelloWorld"},
		{"hangul filler", "Hello\u3164World", "HelloWorld"},
		{"multiple invisible chars", "He\u200Dl\u200Cl\u00ADo", "Hello"},
		{"only invisible chars", "\u200D\u200C\u00AD", "-"},
		{"empty string", "", "-"},
		{"whitespace only after stripping", "\u200D \u200C", "-"},
		{"cyrillic text preserved", "Іван Петрович", "Іван Петрович"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Sanitize(tt.input))
		})
	}
}
