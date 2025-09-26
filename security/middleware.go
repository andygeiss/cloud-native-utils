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
		sessionId := r.URL.Query().Get("session_id")

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
