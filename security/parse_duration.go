package security

import (
	"os"
	"time"
)

// ParseDuration parses the value of the environment variable with the given key as a duration.
// If the value is not set or cannot be parsed, the default duration is returned.
func ParseDuration(key string, def time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return def
}
