package security_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_WithAuth_With_ValidSession_Should_SetSessionIDInContext(t *testing.T) {
	// Arrange
	sessions := security.NewServerSessions()
	sessionID := security.GenerateID()[:32]
	sessions.Create(sessionID, security.IdentityTokenClaims{Name: "John Doe"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", fmt.Sprintf("/ui/%s/", sessionID), nil)
	var got string
	next := func(w http.ResponseWriter, r *http.Request) {
		got = r.Context().Value(security.ContextSessionID).(string)
		w.WriteHeader(http.StatusOK)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ui/{session_id}/", security.WithAuth(sessions, next))

	// Act
	mux.ServeHTTP(w, r)

	// Assert
	assert.That(t, "status code must be 200", w.Code, http.StatusOK)
	assert.That(t, "session_id must be correct", sessionID, got)
}

func Test_WithNoStoreNoReferrer_Should_Set_NoStore_And_NoReferrer(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ui/test/", nil)
	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ui/test/", security.WithNoStoreNoReferrer(next))

	// Act
	mux.ServeHTTP(w, r)

	// Assert
	assert.That(t, "status code must be 200", w.Code, http.StatusOK)
	assert.That(t, "Cache-Control must be no-store", w.Header().Get("Cache-Control"), "no-store")
	assert.That(t, "Referrer-Policy must be no-referrer", w.Header().Get("Referrer-Policy"), "no-referrer")
}
