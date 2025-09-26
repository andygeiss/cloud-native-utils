package logging

import (
	"log/slog"
	"net/http"
	"time"
)

// WithLogging logs the request with method, path and duration.
func WithLogging(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)

		// Log the request with method, path and duration.
		logger.Info(
			"http request handled",
			"method", r.Method,
			"path", r.RequestURI,
			"duration", time.Since(start),
		)
	}
}
