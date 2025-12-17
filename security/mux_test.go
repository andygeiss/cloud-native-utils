package security_test

import (
	"context"
	"embed"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

//go:embed assets
var efs embed.FS

func Test_NewServeMux_With_CanceledContext_Should_ReturnServiceUnavailable(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	mux, _ := security.NewServeMux(ctx, efs)
	req := httptest.NewRequest("GET", "/readiness", nil)
	w := httptest.NewRecorder()
	// Act
	mux.ServeHTTP(w, req)
	// Assert
	assert.That(t, "status code must be 503", w.Code, 503)
}

func Test_NewServeMux_With_LivenessEndpoint_Should_ReturnOK(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mux, _ := security.NewServeMux(ctx, efs)
	req := httptest.NewRequest("GET", "/liveness", nil)
	w := httptest.NewRecorder()
	// Act
	mux.ServeHTTP(w, req)
	// Assert
	assert.That(t, "status code must be 200", w.Code, 200)
	assert.That(t, "body must be OK", w.Body.String(), "OK")
}

func Test_NewServeMux_With_ReadinessEndpoint_Should_ReturnOK(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mux, _ := security.NewServeMux(ctx, efs)
	req := httptest.NewRequest("GET", "/readiness", nil)
	w := httptest.NewRecorder()
	// Act
	mux.ServeHTTP(w, req)
	// Assert
	assert.That(t, "status code must be 200", w.Code, 200)
}

func Test_NewServeMux_With_StaticAssets_Should_ServeFiles(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mux, _ := security.NewServeMux(ctx, efs)
	req := httptest.NewRequest("GET", "/static/keepalive.txt", nil)
	w := httptest.NewRecorder()
	// Act
	mux.ServeHTTP(w, req)
	// Assert
	assert.That(t, "status code must be 200", w.Code, 200)
	assert.That(t, "body must be correct", w.Body.String(), "localhost")
}

func Test_NewServeMux_With_UnknownRoute_Should_Return404(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mux, _ := security.NewServeMux(ctx, efs)
	req := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	// Act
	mux.ServeHTTP(w, req)
	// Assert
	assert.That(t, "status code must be 404", w.Code, 404)
}

func Test_NewServeMux_With_ValidContext_Should_ReturnNonNilMux(t *testing.T) {
	// Arrange
	ctx := context.Background()
	// Act
	mux, _ := security.NewServeMux(ctx, efs)
	// Assert
	assert.That(t, "mux must not be nil", mux != nil, true)
}
