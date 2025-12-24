package utils

import (
	"fmt"
	"time"
)

// GetAge returns a human-readable age string from a time
func GetAge(t time.Time) string {
	duration := time.Since(t)
	
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	} else if duration < time.Hour {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%dh", int(duration.Hours()))
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	}
}

// FormatDuration formats a duration as a human-readable string
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.2fm", d.Minutes())
	} else {
		return fmt.Sprintf("%.2fh", d.Hours())
	}
}

