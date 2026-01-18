package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/mcp"
	"github.com/andygeiss/cloud-native-utils/security"
	"github.com/andygeiss/cloud-native-utils/web"
)

func Test_WithAuth_With_ValidSession_Should_SetSessionIDInContext(t *testing.T) {
	// Arrange
	sessions := web.NewServerSessions()
	sessionID := security.GenerateID()[:32]
	sessions.Create(sessionID, web.IdentityTokenClaims{Name: "John Doe"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ui/", nil)
	r.AddCookie(&http.Cookie{Name: "sid", Value: sessionID})
	var got string
	next := func(w http.ResponseWriter, r *http.Request) {
		got = r.Context().Value(web.ContextSessionID).(string)
		w.WriteHeader(http.StatusOK)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ui/", web.WithAuth(sessions, next))

	// Act
	mux.ServeHTTP(w, r)

	// Assert
	assert.That(t, "status code must be 200", w.Code, http.StatusOK)
	assert.That(t, "session_id must be correct", sessionID, got)
}

func Test_WithAuth_With_ValidSession_Should_SetClaimsInContext(t *testing.T) {
	// Arrange
	sessions := web.NewServerSessions()
	sessionID := security.GenerateID()[:32]
	claims := web.IdentityTokenClaims{
		Email:    "john@example.com",
		Issuer:   "https://issuer.example.com",
		Name:     "John Doe",
		Subject:  "user-123",
		Verified: true,
	}
	sessions.Create(sessionID, claims)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ui/", nil)
	r.AddCookie(&http.Cookie{Name: "sid", Value: sessionID})
	var gotEmail, gotIssuer, gotName, gotSubject string
	var gotVerified bool
	next := func(w http.ResponseWriter, r *http.Request) {
		gotEmail = r.Context().Value(web.ContextEmail).(string)
		gotIssuer = r.Context().Value(web.ContextIssuer).(string)
		gotName = r.Context().Value(web.ContextName).(string)
		gotSubject = r.Context().Value(web.ContextSubject).(string)
		gotVerified = r.Context().Value(web.ContextVerified).(bool)
		w.WriteHeader(http.StatusOK)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ui/", web.WithAuth(sessions, next))

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

func Test_WithBearerAuth_With_MissingAuthHeader_Should_ReturnJSONRPCError(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Act
	web.WithBearerAuth(nil, next)(w, r)

	// Assert
	assert.That(t, "status code must be 401", w.Code, http.StatusUnauthorized)
	assert.That(t, "content-type must be application/json", w.Header().Get("Content-Type"), "application/json")

	var resp mcp.Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.That(t, "jsonrpc must be 2.0", resp.JSONRPC, "2.0")
	assert.That(t, "error must not be nil", resp.Error != nil, true)
	assert.That(t, "error message must be correct", resp.Error.Message, "Missing Authorization header")
}

func Test_WithBearerAuth_With_InvalidAuthScheme_Should_ReturnJSONRPCError(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	r.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Act
	web.WithBearerAuth(nil, next)(w, r)

	// Assert
	assert.That(t, "status code must be 401", w.Code, http.StatusUnauthorized)

	var resp mcp.Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.That(t, "error message must be correct", resp.Error.Message, "Invalid Authorization scheme, expected Bearer")
}

func Test_WithBearerAuth_With_EmptyToken_Should_ReturnJSONRPCError(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	r.Header.Set("Authorization", "Bearer ")
	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Act
	web.WithBearerAuth(nil, next)(w, r)

	// Assert
	assert.That(t, "status code must be 401", w.Code, http.StatusUnauthorized)

	var resp mcp.Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.That(t, "error message must be correct", resp.Error.Message, "Empty Bearer token")
}

func Test_WithBearerAuth_With_NilVerifier_Should_ReturnJSONRPCError(t *testing.T) {
	// Arrange
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	r.Header.Set("Authorization", "Bearer some-token")
	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Act
	web.WithBearerAuth(nil, next)(w, r)

	// Assert
	assert.That(t, "status code must be 401", w.Code, http.StatusUnauthorized)

	var resp mcp.Response
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.That(t, "error message must be correct", resp.Error.Message, "Invalid or expired token")
}
