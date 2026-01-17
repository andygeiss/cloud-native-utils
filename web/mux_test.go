package web_test

import (
	"context"
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/web"
)

//go:embed assets
var efs embed.FS

func Test_NewServeMux_With_CanceledContext_Should_ReturnServiceUnavailable(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	mux, _ := web.NewServeMux(ctx, efs)
	req := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	w := httptest.NewRecorder()
	// Act
	mux.ServeHTTP(w, req)
	// Assert
	assert.That(t, "status code must be 503", w.Code, 503)
}

func Test_NewServeMux_With_LivenessEndpoint_Should_ReturnOK(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mux, _ := web.NewServeMux(ctx, efs)
	req := httptest.NewRequest(http.MethodGet, "/liveness", nil)
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
	mux, _ := web.NewServeMux(ctx, efs)
	req := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	w := httptest.NewRecorder()
	// Act
	mux.ServeHTTP(w, req)
	// Assert
	assert.That(t, "status code must be 200", w.Code, 200)
}

func Test_NewServeMux_With_StaticAssets_Should_ServeFiles(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mux, _ := web.NewServeMux(ctx, efs)
	req := httptest.NewRequest(http.MethodGet, "/static/keepalive.txt", nil)
	w := httptest.NewRecorder()
	// Act
	mux.ServeHTTP(w, req)
	// Assert
	assert.That(t, "status code must be 200", w.Code, 200)
	assert.That(t, "body must be correct", w.Body.String(), "localhost\n")
}

func Test_NewServeMux_With_UnknownRoute_Should_Return404(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mux, _ := web.NewServeMux(ctx, efs)
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
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
	mux, _ := web.NewServeMux(ctx, efs)
	// Assert
	assert.That(t, "mux must not be nil", mux != nil, true)
}
