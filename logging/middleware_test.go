package logging_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/logging"
)

func Test_WithLogging_With_SuccessfulRequest_Should_LogInfo(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	handler := logging.WithLogging(logger, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	assert.That(t, "status must be 200", rec.Code, http.StatusOK)
	assert.That(t, "log must contain info level", strings.Contains(buf.String(), "INFO"), true)
}

func Test_WithLogging_With_ValidRequest_Should_LogMethod(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	handler := logging.WithLogging(logger, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	logOutput := buf.String()
	assert.That(t, "log must contain method", strings.Contains(logOutput, "GET"), true)
}

func Test_WithLogging_With_ValidRequest_Should_LogPath(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	handler := logging.WithLogging(logger, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	logOutput := buf.String()
	assert.That(t, "log must contain path", strings.Contains(logOutput, "/test-path"), true)
}

func Test_WithLogging_Without_Custom_Config_Should_Not_Mask_Paths(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	handler := logging.WithLogging(logger, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/ui/secret-session-123/checkout/return?token=paypal-token", nil)
	req.SetPathValue("session_id", "secret-session-123")

	// Act
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	// Assert
	assert.That(t, "status code should be 200", res.Code, http.StatusOK)
	logOutput := buf.String()
	assert.That(t, "log must not contain query token", strings.Contains(logOutput, "paypal-token"), false)
	assert.That(t, "log must contain raw session id (not masked by default)", strings.Contains(logOutput, "secret-session-123"), true)
}

func Test_WithLogging_With_ValidRequest_Should_LogRequestMessage(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	handler := logging.WithLogging(logger, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	logOutput := buf.String()
	assert.That(t, "log must contain request message", strings.Contains(logOutput, "http request handled"), true)
}

func Test_WithLoggingCustom_With_CustomSanitizer_Should_Mask_CustomPaths(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	// Create custom sanitizer with only custom path values
	customSanitizer := logging.NewPathSanitizer(map[string]string{
		"user_id": "{user_id}",
		"team_id": "{team_id}",
	})

	handler := logging.WithLoggingCustom(logger, customSanitizer, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/teams/team-123/users/user-456/profile", nil)
	req.SetPathValue("team_id", "team-123")
	req.SetPathValue("user_id", "user-456")
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	logOutput := buf.String()
	assert.That(t, "log must contain masked team_id", strings.Contains(logOutput, "{team_id}"), true)
	assert.That(t, "log must contain masked user_id", strings.Contains(logOutput, "{user_id}"), true)
	assert.That(t, "log must not contain raw team_id", strings.Contains(logOutput, "team-123"), false)
	assert.That(t, "log must not contain raw user_id", strings.Contains(logOutput, "user-456"), false)
	assert.That(t, "log must contain profile path", strings.Contains(logOutput, "/profile"), true)
}

func Test_WithLoggingCustom_With_EmptySanitizer_Should_Not_Mask_Anything(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	// Create empty sanitizer - masks nothing
	emptySanitizer := logging.NewPathSanitizer(map[string]string{})

	handler := logging.WithLoggingCustom(logger, emptySanitizer, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/users/secret-id-123/data", nil)
	req.SetPathValue("user_id", "secret-id-123")
	rec := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(rec, req)

	// Assert
	logOutput := buf.String()
	assert.That(t, "log must contain raw user_id when no sanitization configured", strings.Contains(logOutput, "secret-id-123"), true)
}
