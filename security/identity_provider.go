package security

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/andygeiss/cloud-native-utils/resource"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// identityProvider represents an identity provider.
type identityProvider struct {
	oauth2Config       *oauth2.Config
	oidcConfig         *oidc.Config
	provider           *oidc.Provider
	stateCodeVerifiers resource.Access[string, string]
}

type ContextKey string

const (
	ContextSessionID ContextKey = "session_id"
	ContextEmail     ContextKey = "email"
	ContextIssuer    ContextKey = "issuer"
	ContextName      ContextKey = "name"
	ContextSubject   ContextKey = "subject"
	ContextVerified  ContextKey = "verified"
)

// IdentityTokenClaims represents the claims of an identity token.
type IdentityTokenClaims struct {
	Email    string `json:"email"`
	Issuer   string `json:"iss"`
	Name     string `json:"name"`
	Subject  string `json:"sub"`
	Verified bool   `json:"email_verified"`
}

// NewIdentityProvider creates a new identity provider.
func NewIdentityProvider() *identityProvider {
	return &identityProvider{
		stateCodeVerifiers: resource.NewInMemoryAccess[string, string](),
	}
}

// IdentityProvider is a singleton instance of the identity provider.
var IdentityProvider = NewIdentityProvider()

// Callback returns a handler function for the identity provider's callback endpoint.
func (a *identityProvider) Callback(sessions *ServerSessions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the code and state parameters from the request URL.
		ctx := context.Background()
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		// Retrieve the code verifier from the state.
		codeVerifier, _ := a.stateCodeVerifiers.Read(state)
		if codeVerifier == nil {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}

		// Exchange the authorization code for an access token.
		token, err := a.oauth2Config.Exchange(ctx, code,
			oauth2.SetAuthURLParam("code_verifier", *codeVerifier),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("token exchange failed: %v", err), http.StatusBadRequest)
			return
		}

		// Parse and verify ID Token
		rawToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "no id_token in token response", http.StatusBadRequest)
			return
		}

		// Verify the ID token using the provider's verifier.
		verifier := a.provider.Verifier(a.oidcConfig)
		idToken, err := verifier.Verify(ctx, rawToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("verify id_token: %v", err), http.StatusUnauthorized)
			return
		}

		// Extract the claims from the ID token.
		var claims IdentityTokenClaims
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, fmt.Sprintf("read claims: %v", err), http.StatusInternalServerError)
			return
		}

		// Generate a unique session ID and create a session with the claims as data.
		sessionId := GenerateID()[:32]
		sessions.Create(sessionId, claims)
		redirectUrl := fmt.Sprintf("%s?session_id=%s", os.Getenv("REDIRECT_URL"), sessionId)

		// Redirect the user to the redirect URL.
		http.Redirect(w, r, redirectUrl, http.StatusFound)
	}
}

// Login returns a handler function for the identity provider's login endpoint.
func (a *identityProvider) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure that the identity provider is properly configured.
		if a.oauth2Config == nil {
			if err := a.setup(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Generate a unique state identifier and a PKCE code verifier/challenge pair.
		state := GenerateID()
		codeVerifier, challenge := GeneratePKCE()

		// Store the state and code verifier for further use.
		a.stateCodeVerifiers.Create(state, codeVerifier)

		// Create the authorization URL with the PKCE parameters.
		authUrl := a.oauth2Config.AuthCodeURL(state,
			oauth2.SetAuthURLParam("code_challenge", challenge),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		)

		// Redirect the user to the authorization URL.
		http.Redirect(w, r, authUrl, http.StatusFound)
	}
}

// Logout handles the logout request.
func (a *identityProvider) Logout(sessions *ServerSessions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the session ID from the request.
		sessionId := r.URL.Query().Get("session_id")

		// Delete the session from the server.
		sessions.Delete(sessionId)

		// Redirect the user to the logout URL.
		http.Redirect(w, r, os.Getenv("REDIRECT_URL"), http.StatusFound)
	}
}

// WithAuth adds authentication information to the context.
func (a *identityProvider) WithAuth(sessions *ServerSessions, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the current context.
		ctx := r.Context()

		// Retrieve the session ID from the request URL.
		sessionId := r.URL.Query().Get("session_id")

		if sessionId != "" {
			// Retrieve the session by using the session ID.
			if session, ok := sessions.Read(sessionId); ok {
				claims, _ := session.Data.(IdentityTokenClaims)

				// Add the claims to the context.
				ctx = context.WithValue(ctx, ContextEmail, claims.Email)
				ctx = context.WithValue(ctx, ContextIssuer, claims.Issuer)
				ctx = context.WithValue(ctx, ContextName, claims.Name)
				ctx = context.WithValue(ctx, ContextSubject, claims.Subject)
			}

			// Add the session ID to the context.
			ctx = context.WithValue(ctx, ContextSessionID, sessionId)
		}

		// Call the next http handler with context.
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (a *identityProvider) setup() (err error) {
	// Create a new context.
	ctx := context.Background()

	// Initialize the identity provider.
	provider, err := oidc.NewProvider(ctx, os.Getenv("OIDC_ISSUER"))
	if err != nil {
		return err
	}

	// Create and store the OAuth2 configuration.
	a.oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("OIDC_CLIENT_ID"),
		ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  os.Getenv("OIDC_REDIRECT_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	// Create and store the OpenID Connect configuration.
	a.oidcConfig = &oidc.Config{
		ClientID: os.Getenv("OIDC_CLIENT_ID"),
	}

	// Store the identity provider.
	a.provider = provider

	return nil
}
