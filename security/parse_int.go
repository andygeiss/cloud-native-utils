package security

import (
	"os"
	"strconv"
)

// ParseInt parses the value of the environment variable with the given key as an integer.
// If the value is not set or cannot be parsed, the default integer is returned.
func ParseInt(key string, def int) int {
	if value := os.Getenv(key); value != "" {
		if duration, err := strconv.Atoi(value); err == nil {
			return duration
		}
	}
	return def
}
