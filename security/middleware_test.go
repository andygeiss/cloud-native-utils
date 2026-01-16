package security_test

import (
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
	r := httptest.NewRequest("GET", "/ui/", nil)
	r.AddCookie(&http.Cookie{Name: "sid", Value: sessionID})
	var got string
	next := func(w http.ResponseWriter, r *http.Request) {
		got = r.Context().Value(security.ContextSessionID).(string)
		w.WriteHeader(http.StatusOK)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ui/", security.WithAuth(sessions, next))

	// Act
	mux.ServeHTTP(w, r)

	// Assert
	assert.That(t, "status code must be 200", w.Code, http.StatusOK)
	assert.That(t, "session_id must be correct", sessionID, got)
}
