package security_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestServeMux_Is_Not_Nil(t *testing.T) {
	ctx := context.Background()
	mux := security.Mux(ctx)
	assert.That(t, "mux must not be nil", mux != nil, true)
}

func TestServeMux_Has_Health_Check(t *testing.T) {
	ctx := context.Background()
	mux := security.Mux(ctx)
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.That(t, "status code must be 200", w.Code, 200)
	assert.That(t, "body must be OK", w.Body.String(), "OK")
}

func TestServeMux_Has_Readiness_Check_When_Context_Active(t *testing.T) {
	ctx := context.Background()
	mux := security.Mux(ctx)
	req := httptest.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.That(t, "status code must be 200", w.Code, 200)
	// The body is empty in this example, but you can also check it if needed.
}

func TestServeMux_Has_Readiness_Check_When_Context_Canceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Immediately cancel the context.
	mux := security.Mux(ctx)
	req := httptest.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.That(t, "status code must be 503", w.Code, 503)
}

func TestServeMux_Unknown_Route(t *testing.T) {
	ctx := context.Background()
	mux := security.Mux(ctx)
	req := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.That(t, "status code must be 404", w.Code, 404)
}
