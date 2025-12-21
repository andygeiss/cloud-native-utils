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

func Test_WithLogging_With_SensitiveRequest_Should_Not_Log_Query_Or_SessionID(t *testing.T) {
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
	assert.That(t, "log must contain masked path", strings.Contains(logOutput, "/ui/{session_id}/checkout/return"), true)
	assert.That(t, "log must not contain query token", strings.Contains(logOutput, "paypal-token"), false)
	assert.That(t, "log must not contain raw session id", strings.Contains(logOutput, "secret-session-123"), false)
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
