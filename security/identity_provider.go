package security

import "context"

// GetLoginResponse represents a login request to an identity provider.
type GetLoginResponse struct {
	CodeVerifier string `json:"code_verifier"`
	State        string `json:"state"`
	URL          string `json:"url"`
}

// VerifyLoginRequest represents a login verification request to an identity provider.
type VerifyLoginRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// VerifyLoginResponse represents a login verification response from an identity provider.
type VerifyLoginResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// IdentityProvider is an interface for identity providers.
type IdentityProvider interface {
	GetLogin(ctx context.Context) (res GetLoginResponse, err error)
	VerifyLogin(ctx context.Context, req VerifyLoginRequest) (res VerifyLoginResponse, err error)
}
