package web_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/web"
)

func setupOIDCEnv(t *testing.T) {
	t.Helper()
	t.Setenv("OIDC_CLIENT_ID", "demo")
	t.Setenv("OIDC_CLIENT_SECRET", "rMCl4R5gBNulChi3bnwu5pp3zXIUKseQ")
	t.Setenv("OIDC_ISSUER", "http://localhost:8180/realms/local")
	t.Setenv("OIDC_REDIRECT_URL", "http://localhost:8080/auth/callback")
}

func Test_IdentityProviderCallback_With_MissingState_Should_ReturnBadRequest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	setupOIDCEnv(t)

	// Arrange
	sessions := web.NewServerSessions()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/auth/callback", nil)

	// Act
	web.IdentityProvider.Callback(sessions)(w, r)

	// Assert
	assert.That(t, "status code must be 400", w.Code, http.StatusBadRequest)
}

func Test_IdentityProviderCallback_With_ValidSession_Should_ProcessRequest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	setupOIDCEnv(t)

	// Arrange
	sessions := web.NewServerSessions()
	sessions.Create("test-id", "test-data")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/auth/login", nil)

	// Act
	web.IdentityProvider.Callback(sessions)(w, r)

	// Assert
	assert.That(t, "status code must be 400", w.Code, http.StatusBadRequest)
}

func Test_IdentityProviderLogin_With_ValidRequest_Should_RedirectWithOIDCParams(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	setupOIDCEnv(t)

	// Arrange
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/auth/login", nil)

	// Act
	web.IdentityProvider.Login()(w, r)

	// Assert
	location := w.Header().Get("Location")
	assert.That(t, "status code must be 302", w.Code, 302)
	assert.That(t, "location has client_id", strings.Contains(location, "client_id"), true)
	assert.That(t, "location has code_challenge", strings.Contains(location, "code_challenge"), true)
	assert.That(t, "location has code_challenge_method", strings.Contains(location, "code_challenge_method"), true)
	assert.That(t, "location has redirect_uri", strings.Contains(location, "redirect_uri"), true)
	assert.That(t, "location has scope", strings.Contains(location, "scope"), true)
	assert.That(t, "location has state", strings.Contains(location, "state"), true)
}

func Test_IdentityProviderLogout_With_ValidSession_Should_DeleteSessionAndRedirect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	setupOIDCEnv(t)

	// Arrange
	sessions := web.NewServerSessions()
	sessions.Create("test-id", "test-data")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	r.AddCookie(&http.Cookie{Name: "sid", Value: "test-id"})

	// Act
	web.IdentityProvider.Logout(sessions)(w, r)
	_, exists := sessions.Read("test-id")

	// Assert
	assert.That(t, "status code must be 302", w.Code, 302)
	assert.That(t, "session must be deleted", exists, false)
}
