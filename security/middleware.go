package security

import (
	"context"
	"net/http"
	"strings"
)

type ContextKey string

const (
	ContextEmail     ContextKey = "email"
	ContextIssuer    ContextKey = "issuer"
	ContextName      ContextKey = "name"
	ContextSessionID ContextKey = "session_id"
	ContextSubject   ContextKey = "subject"
	ContextVerified  ContextKey = "verified"
)

// WithAuth adds authentication information to the context.
func WithAuth(sessions *ServerSessions, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a new context.
		ctx := context.Background()

		// Read session ID from cookie.
		var sessionID string
		if c, err := r.Cookie("sid"); err == nil {
			sessionID = c.Value
		}

		// Define the claims
		var email, issuer, name, subject string
		var verified bool

		if sessionID != "" {
			// Retrieve the session by using the session ID.
			if session, ok := sessions.Read(sessionID); ok {
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
		ctx = context.WithValue(ctx, ContextSessionID, sessionID)
		ctx = context.WithValue(ctx, ContextSubject, subject)
		ctx = context.WithValue(ctx, ContextVerified, verified)

		// Prevent sensitive information leakage via Referer header.
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Prevent content type sniffing.
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking.
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevent caching of dynamic/sensitive responses.
		// If you serve static assets from /static, let them be cached.
		if !strings.HasPrefix(r.URL.Path, "/static/") {
			w.Header().Set("Cache-Control", "no-store")
		}

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
