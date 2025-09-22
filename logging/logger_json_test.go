package logging_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/logging"
)

func TestNewJsonLogger_With_Debug_Level(t *testing.T) {
	os.Setenv("LOGGING_LEVEL", "DEBUG")
	ctx := context.Background()
	logger := logging.NewJsonLogger()
	assert.That(t, "logging level must be debug", logger.Handler().Enabled(ctx, slog.LevelDebug), true)
}

func TestNewJsonLogger_With_Error_Level(t *testing.T) {
	os.Setenv("LOGGING_LEVEL", "ERror")
	ctx := context.Background()
	logger := logging.NewJsonLogger()
	assert.That(t, "logging level must be error", logger.Handler().Enabled(ctx, slog.LevelError), true)
}

func TestNewJsonLogger_With_Info_Level(t *testing.T) {
	os.Setenv("LOGGING_LEVEL", "info")
	ctx := context.Background()
	logger := logging.NewJsonLogger()
	assert.That(t, "logging level must be info", logger.Handler().Enabled(ctx, slog.LevelInfo), true)
}

func TestNewJsonLogger_With_Warn_Level(t *testing.T) {
	os.Setenv("LOGGING_LEVEL", "Warn")
	ctx := context.Background()
	logger := logging.NewJsonLogger()
	assert.That(t, "logging level must be warn", logger.Handler().Enabled(ctx, slog.LevelWarn), true)
}
