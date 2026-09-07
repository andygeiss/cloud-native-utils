package logging_test

import (
	"log/slog"

	"github.com/andygeiss/cloud-native-utils/logging"
)

// This example has no Output comment on purpose: the logger writes JSON with a
// timestamp, so the bytes differ on every run.
func ExampleNewJsonLogger() {
	// LOGGING_LEVEL picks the level; it defaults to INFO.
	logger := logging.NewJsonLogger()

	logger.Info("server started", slog.Int("port", 8080))
	logger.Error("upstream unreachable", slog.String("host", "db.internal"))
}
