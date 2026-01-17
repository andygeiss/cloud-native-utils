package logging

import (
	"log/slog"
	"os"
	"strings"
)

// NewJsonLogger creates a new structured logger in JSON format.
func NewJsonLogger() *slog.Logger {
	// Configure the level by using the environment.
	var level slog.Leveler

	lvl := os.Getenv("LOGGING_LEVEL")
	lvl = strings.ToUpper(lvl)

	switch lvl {
	case "DEBUG":
		level = slog.LevelDebug
	case "ERROR":
		level = slog.LevelError
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
	}

	// Create a new handler for structured logs.
	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewJSONHandler(os.Stdout, opts)

	// Create and return a new logger.
	return slog.New(handler)
}
