package utils

import (
	"strings"
	"time"
)

type ContextKey string

const RequestIdentifier = "X-Request-ID"

func IsDevMode(s string) bool {
	return strings.Contains(s, "local") || strings.Contains(s, "dev")
}
func FormatTimeToISO(timeToFormat time.Time) string {
	return timeToFormat.Format(time.RFC3339)
}

func CurrentISOTime() string {
	return FormatTimeToISO(time.Now().UTC())
}
