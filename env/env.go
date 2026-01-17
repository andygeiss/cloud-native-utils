package env

import (
	"os"
	"strconv"
	"time"
)

// Get retrieves an environment variable and parses it to the type of the default value.
// If the value is not set or cannot be parsed, the default value is returned.
// Supported types: bool, int, float64, string, time.Duration
func Get[T any](key string, def T) T {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	result, ok := parse[T](value, def)
	if !ok {
		return def
	}
	return result
}

func parse[T any](value string, def T) (T, bool) {
	var result any
	var err error

	switch any(def).(type) {
	case bool:
		result, err = strconv.ParseBool(value)
	case int:
		result, err = strconv.Atoi(value)
	case float64:
		result, err = strconv.ParseFloat(value, 64)
	case string:
		return any(value).(T), true
	case time.Duration:
		result, err = time.ParseDuration(value)
	default:
		return def, false
	}

	if err != nil {
		return def, false
	}
	return result.(T), true
}
