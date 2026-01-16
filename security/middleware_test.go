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

func Test_WithAuth_With_ValidSession_Should_SetClaimsInContext(t *testing.T) {
	// Arrange
	sessions := security.NewServerSessions()
	sessionID := security.GenerateID()[:32]
	claims := security.IdentityTokenClaims{
		Email:    "john@example.com",
		Issuer:   "https://issuer.example.com",
		Name:     "John Doe",
		Subject:  "user-123",
		Verified: true,
	}
	sessions.Create(sessionID, claims)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ui/", nil)
	r.AddCookie(&http.Cookie{Name: "sid", Value: sessionID})
	var gotEmail, gotIssuer, gotName, gotSubject string
	var gotVerified bool
	next := func(w http.ResponseWriter, r *http.Request) {
		gotEmail = r.Context().Value(security.ContextEmail).(string)
		gotIssuer = r.Context().Value(security.ContextIssuer).(string)
		gotName = r.Context().Value(security.ContextName).(string)
		gotSubject = r.Context().Value(security.ContextSubject).(string)
		gotVerified = r.Context().Value(security.ContextVerified).(bool)
		w.WriteHeader(http.StatusOK)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ui/", security.WithAuth(sessions, next))

	// Act
	mux.ServeHTTP(w, r)

	// Assert
	assert.That(t, "status code must be 200", w.Code, http.StatusOK)
	assert.That(t, "email must be correct", gotEmail, "john@example.com")
	assert.That(t, "issuer must be correct", gotIssuer, "https://issuer.example.com")
	assert.That(t, "name must be correct", gotName, "John Doe")
	assert.That(t, "subject must be correct", gotSubject, "user-123")
	assert.That(t, "verified must be correct", gotVerified, true)
}
