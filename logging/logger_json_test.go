package logging_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/logging"
)

func Test_NewJsonLogger_With_DebugLevel_Should_EnableDebugLogging(t *testing.T) {
	// Arrange
	t.Setenv("LOGGING_LEVEL", "DEBUG")
	ctx := context.Background()

	// Act
	logger := logging.NewJsonLogger()

	// Assert
	assert.That(t, "logging level must be debug", logger.Handler().Enabled(ctx, slog.LevelDebug), true)
}

func Test_NewJsonLogger_With_DefaultLevel_Should_EnableInfoLogging(t *testing.T) {
	// Arrange
	_ = os.Unsetenv("LOGGING_LEVEL")
	ctx := context.Background()

	// Act
	logger := logging.NewJsonLogger()

	// Assert
	assert.That(t, "logging level must be info", logger.Handler().Enabled(ctx, slog.LevelInfo), true)
	assert.That(t, "debug should be disabled", logger.Handler().Enabled(ctx, slog.LevelDebug), false)
}

func Test_NewJsonLogger_With_ErrorLevel_Should_EnableErrorLogging(t *testing.T) {
	// Arrange
	t.Setenv("LOGGING_LEVEL", "ERror")
	ctx := context.Background()

	// Act
	logger := logging.NewJsonLogger()

	// Assert
	assert.That(t, "logging level must be error", logger.Handler().Enabled(ctx, slog.LevelError), true)
}

func Test_NewJsonLogger_With_InfoLevel_Should_EnableInfoLogging(t *testing.T) {
	// Arrange
	t.Setenv("LOGGING_LEVEL", "info")
	ctx := context.Background()

	// Act
	logger := logging.NewJsonLogger()

	// Assert
	assert.That(t, "logging level must be info", logger.Handler().Enabled(ctx, slog.LevelInfo), true)
}

func Test_NewJsonLogger_With_InvalidLevel_Should_DefaultToInfoLogging(t *testing.T) {
	// Arrange
	t.Setenv("LOGGING_LEVEL", "INVALID")
	ctx := context.Background()

	// Act
	logger := logging.NewJsonLogger()

	// Assert
	assert.That(t, "logging level must default to info", logger.Handler().Enabled(ctx, slog.LevelInfo), true)
	assert.That(t, "debug should be disabled", logger.Handler().Enabled(ctx, slog.LevelDebug), false)
}

func Test_NewJsonLogger_With_WarnLevel_Should_EnableWarnLogging(t *testing.T) {
	// Arrange
	t.Setenv("LOGGING_LEVEL", "Warn")
	ctx := context.Background()

	// Act
	logger := logging.NewJsonLogger()

	// Assert
	assert.That(t, "logging level must be warn", logger.Handler().Enabled(ctx, slog.LevelWarn), true)
}
