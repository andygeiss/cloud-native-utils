package security

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
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

		// Prevent caching of dynamic/sensitive responses.
		// If you serve static assets from /static, let them be cached.
		if !strings.HasPrefix(r.URL.Path, "/static/") {
			w.Header().Set("Cache-Control", "no-store")
		}

		// Prevent sensitive information leakage via Referer header.
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Prevent content type sniffing.
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking.
		w.Header().Set("X-Frame-Options", "DENY")

		// Call the next http handler with context.
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// WithLogging logs the request with method, path and duration.
func WithLogging(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)

		// Log the request with method, path and duration.
		logger.Info(
			"http request handled",
			"method", r.Method,
			"path", r.RequestURI,
			"duration", time.Since(start),
		)
	}
}
