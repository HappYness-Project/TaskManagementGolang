package utils

import "strings"

type ContextKey string

const RequestIdentifier = "X-Request-ID"

func IsDevMode(s string) bool {
	return strings.Contains(s, "local") || strings.Contains(s, "dev")
}
