package security

import (
	"os"
	"strconv"
	"time"
)

// ParseDurationOrDefault parses the value of the environment variable with the given key as a duration.
// If the value is not set or cannot be parsed, the default duration is returned.
func ParseDurationOrDefault(key string, def time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return def
}

// ParseIntOrDefault parses the value of the environment variable with the given key as an integer.
// If the value is not set or cannot be parsed, the default integer is returned.
func ParseIntOrDefault(key string, def int) int {
	if value := os.Getenv(key); value != "" {
		if duration, err := strconv.Atoi(value); err == nil {
			return duration
		}
	}
	return def
}
