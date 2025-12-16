package security_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_Middleware_WithAuth_Should_Contain_SessionID(t *testing.T) {
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
