package db

import "strings"

func NormalizeJobStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "complete":
		return "completed"
	default:
		return strings.ToLower(strings.TrimSpace(status))
	}
}
