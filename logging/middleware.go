package logging

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func sanitizePath(r *http.Request) string {
	// Never log query strings (can contain secrets like OIDC codes or payment tokens).
	path := r.URL.Path

	// Mask common path identifiers that are considered secrets or sensitive.
	path = maskPathValue(path, r.PathValue("session_id"), "{session_id}")
	path = maskPathValue(path, r.PathValue("external_order_id"), "{external_order_id}")

	// These are not necessarily secrets, but masking reduces correlation/leakage.
	path = maskPathValue(path, r.PathValue("booking_id"), "{booking_id}")
	path = maskPathValue(path, r.PathValue("hold_id"), "{hold_id}")
	path = maskPathValue(path, r.PathValue("invoice_id"), "{invoice_id}")

	return path
}

func maskPathValue(path, value, placeholder string) string {
	if value == "" {
		return path
	}

	// Replace full segment occurrences only.
	segment := "/" + value

	// Middle segments.
	path = strings.ReplaceAll(path, segment+"/", "/"+placeholder+"/")
	// Trailing segment.
	if trimmed, ok := strings.CutSuffix(path, segment); ok {
		path = trimmed + "/" + placeholder
	}

	return path
}

// WithLogging logs the request with method, path and duration.
func WithLogging(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		nowUTC := time.Now().UTC()

		// Log the request with method, path and duration.
		logger.Info(
			"http request handled",
			"method", r.Method,
			"path", sanitizePath(r),
			"duration", time.Since(start),
			"now_utc", nowUTC,
		)
	}
}
