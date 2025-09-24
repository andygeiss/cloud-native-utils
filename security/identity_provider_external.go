package security

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/andygeiss/cloud-native-utils/resource"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// externalIdentityProvider represents an external identity provider.
type externalIdentityProvider struct {
	config     *oauth2.Config
	stateCodes resource.Access[string, string]
	verifier   *oidc.IDTokenVerifier
}

// NewExternalIdentityProvider creates a new external identity provider.
func NewExternalIdentityProvider() IdentityProvider {
	return &externalIdentityProvider{
		stateCodes: resource.NewInMemoryAccess[string, string](),
	}
}

// GetLoginURL returns the login URL for the identity provider.
func (a *externalIdentityProvider) GetLogin(ctx context.Context) (res GetLoginResponse, err error) {
	// Generate PKCE verifier and challenge.
	codeVerifier := GenerateID() // 43–128 chars recommended; ensure yours does that.
	sum := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(sum[:])
	state := GenerateID()[:32]

	// Store the *verifier* (not the challenge) for later retrieval.
	if err := a.stateCodes.Create(state, codeVerifier); err != nil {
		return GetLoginResponse{}, err
	}

	// Lazy-init config/verifier.
	if a.config == nil || a.verifier == nil {
		a.config, a.verifier, err = getIdentityProviderConfig(ctx)
		if err != nil {
			return GetLoginResponse{}, err
		}
	}

	// Build auth URL with PKCE params.
	url := a.config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	return GetLoginResponse{
		CodeVerifier: codeVerifier,
		State:        state,
		URL:          url,
	}, nil
}

// VerifyLogin verifies the login code and state.
func (a *externalIdentityProvider) VerifyLogin(ctx context.Context, req VerifyLoginRequest) (res VerifyLoginResponse, err error) {
	if req.State == "" || req.Code == "" {
		return VerifyLoginResponse{}, errors.New("missing state or code")
	}

	// Lazy-init (in case Verify is called first in some flows/tests).
	if a.config == nil || a.verifier == nil {
		a.config, a.verifier, err = getIdentityProviderConfig(ctx)
		if err != nil {
			return VerifyLoginResponse{}, err
		}
	}

	// Load the stored PKCE verifier using the state.
	codeVerifierPtr, err := a.stateCodes.Read(req.State)
	if err != nil || codeVerifierPtr == nil {
		return VerifyLoginResponse{}, errors.New("invalid state or verifier not found")
	}
	codeVerifier := *codeVerifierPtr

	// Delete the stored PKCE verifier after use.
	if err := a.stateCodes.Delete(req.State); err != nil {
		return VerifyLoginResponse{}, err
	}

	// Exchange the code for tokens using the original verifier.
	token, err := a.config.Exchange(ctx, req.Code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
	if err != nil {
		return VerifyLoginResponse{}, err
	}
	if token == nil {
		return VerifyLoginResponse{}, errors.New("token is nil")
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return VerifyLoginResponse{}, errors.New("no id_token")
	}

	// Verify the ID token signature & standard claims.
	idt, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return VerifyLoginResponse{}, err
	}

	// Extract claims from the *ID token*.
	var claims struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := idt.Claims(&claims); err != nil {
		return VerifyLoginResponse{}, fmt.Errorf("parsing id token claims: %w", err)
	}

	return VerifyLoginResponse{
		Email: claims.Email,
		Name:  claims.Name,
	}, nil
}

// getIdentityProvider returns the identity provider endpoints.
func getIdentityProvider(ctx context.Context) (provider *oidc.Provider, err error) {

	// Get the identity provider endpoint from the configuration.
	issuer := fmt.Sprintf("%s/realms/%s", os.Getenv("IDP_URL"), os.Getenv("IDP_REALM"))

	// Send a request to the provider in the background to get the endpoints.
	return oidc.NewProvider(ctx, issuer)
}

// getIdentityProviderConfig returns the identity provider configuration.
func getIdentityProviderConfig(ctx context.Context) (config *oauth2.Config, verifier *oidc.IDTokenVerifier, err error) {
	// Get the identity provider endpoint from the configuration.
	provider, err := getIdentityProvider(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Use the same provider for the verifier.
	clientID := os.Getenv("IDP_CLIENT_ID")
	verifier = provider.Verifier(&oidc.Config{ClientID: clientID})

	// Return the identity provider configuration.
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: os.Getenv("IDP_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("IDP_REDIRECT_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		Endpoint:     provider.Endpoint(),
	}, verifier, nil
}
