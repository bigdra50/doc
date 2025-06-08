package utils

import (
	"fmt"
	"time"
)

// FormatDuration formats a duration for display
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%.1fm", d.Minutes())
}

// MaskAPIKey masks an API key for safe logging
func MaskAPIKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) < 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
