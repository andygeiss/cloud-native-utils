package logging

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// PathSanitizer configures which path segments should be masked during logging.
// Example usage:
//
//	sanitizer := logging.NewPathSanitizer(map[string]string{
//		"session_id": "{session_id}",
//		"user_id": "{user_id}",
//		"order_id": "{order_id}",
//	})
//	handler := logging.WithLoggingCustom(logger, sanitizer, nextHandler)
type PathSanitizer struct {
	// PathValues maps path parameter names to their placeholder representations.
	PathValues map[string]string
}

// NewPathSanitizer creates a custom path sanitizer with specified path values.
func NewPathSanitizer(pathValues map[string]string) *PathSanitizer {
	return &PathSanitizer{
		PathValues: pathValues,
	}
}

// Sanitize applies masking rules to the HTTP request path.
func (ps *PathSanitizer) Sanitize(r *http.Request) string {
	// Never log query strings (can contain secrets like OIDC codes or payment tokens).
	path := r.URL.Path

	// Mask configured path identifiers.
	for paramName, placeholder := range ps.PathValues {
		value := r.PathValue(paramName)
		path = maskPathValue(path, value, placeholder)
	}

	return path
}

// maskPathValue replaces sensitive path segments with placeholders.
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

// WithLogging logs the request with method, path and duration without path masking.
// To mask sensitive path segments, use WithLoggingCustom with a configured PathSanitizer.
func WithLogging(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return WithLoggingCustom(logger, NewPathSanitizer(make(map[string]string)), next)
}

// WithLoggingCustom logs the request with method, path and duration using a custom path sanitizer.
func WithLoggingCustom(logger *slog.Logger, sanitizer *PathSanitizer, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		nowUTC := time.Now().UTC()

		// Log the request with method, path and duration.
		logger.Info(
			"http request handled",
			"method", r.Method,
			"path", sanitizer.Sanitize(r),
			"duration", time.Since(start),
			"now_utc", nowUTC,
		)
	}
}
