package security

import (
	"context"
	"net/http"
)

type ContextKey string

const (
	ContextSessionID ContextKey = "session_id"
	ContextEmail     ContextKey = "email"
	ContextIssuer    ContextKey = "issuer"
	ContextName      ContextKey = "name"
	ContextSubject   ContextKey = "subject"
	ContextVerified  ContextKey = "verified"
)

// WithAuth adds authentication information to the context.
func WithAuth(sessions *ServerSessions, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a new context.
		ctx := context.Background()

		// Retrieve the session ID from the request URL.
		sessionId := r.PathValue("session_id")

		// Define the claims
		var email, issuer, name, subject string
		var verified bool

		if sessionId != "" {
			// Retrieve the session by using the session ID.
			if session, ok := sessions.Read(sessionId); ok {
				claims, _ := session.Data.(IdentityTokenClaims)
				email = claims.Email
				issuer = claims.Issuer
				name = claims.Name
				subject = claims.Subject
				verified = claims.Verified
			}
		}

		// Add the authentication informations.
		ctx = context.WithValue(ctx, ContextEmail, email)
		ctx = context.WithValue(ctx, ContextIssuer, issuer)
		ctx = context.WithValue(ctx, ContextName, name)
		ctx = context.WithValue(ctx, ContextSessionID, sessionId)
		ctx = context.WithValue(ctx, ContextSubject, subject)
		ctx = context.WithValue(ctx, ContextVerified, verified)

		// Call the next http handler with context.
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// WithNoStoreNoReferrer applies response headers that reduce data leakage.
//
// It sets:
// - Cache-Control: no-store
// - Referrer-Policy: no-referrer
func WithNoStoreNoReferrer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	}
}

// WithAuthenticatedSecurityHeaders applies server-side security headers that MUST
// be present on authenticated pages.
//
// Deprecated: use WithNoStoreNoReferrer.
func WithAuthenticatedSecurityHeaders(next http.HandlerFunc) http.HandlerFunc {
	return WithNoStoreNoReferrer(next)
}
