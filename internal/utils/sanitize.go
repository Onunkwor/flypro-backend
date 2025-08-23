package utils

import (
	"html"
	"strings"
)

func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	s = html.EscapeString(s)
	return s
}

func NormalizeCurrency(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

func NormalizeCategory(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
