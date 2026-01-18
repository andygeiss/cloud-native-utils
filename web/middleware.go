package web

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/andygeiss/cloud-native-utils/mcp"
	"github.com/coreos/go-oidc/v3/oidc"
)

// ContextKey is a type for context keys used in the web package.
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
		// Use the request context.
		ctx := r.Context()

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
			if session, ok := sessions.Read(sessionID); ok && session != nil {
				if claims, ok := session.Data.(IdentityTokenClaims); ok {
					email = claims.Email
					issuer = claims.Issuer
					name = claims.Name
					subject = claims.Subject
					verified = claims.Verified
				}
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

// WithBearerAuth validates OAuth 2.1 Bearer tokens for MCP endpoints.
// It extracts the token from the Authorization header, verifies it against
// the OIDC provider, and populates the request context with user claims.
// On failure, it returns a JSON-RPC 2.0 error response.
func WithBearerAuth(verifier *oidc.IDTokenVerifier, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract Bearer token from Authorization header.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONRPCError(w, "Missing Authorization header")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeJSONRPCError(w, "Invalid Authorization scheme, expected Bearer")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			writeJSONRPCError(w, "Empty Bearer token")
			return
		}

		// Check for nil verifier (configuration error).
		if verifier == nil {
			writeJSONRPCError(w, "Invalid or expired token")
			return
		}

		// Verify token with OIDC provider.
		idToken, err := verifier.Verify(r.Context(), tokenString)
		if err != nil {
			writeJSONRPCError(w, "Invalid or expired token")
			return
		}

		// Extract claims from token.
		var claims IdentityTokenClaims
		if err := idToken.Claims(&claims); err != nil {
			writeJSONRPCError(w, "Failed to parse token claims")
			return
		}

		// Populate context with authenticated user info.
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextEmail, claims.Email)
		ctx = context.WithValue(ctx, ContextIssuer, claims.Issuer)
		ctx = context.WithValue(ctx, ContextName, claims.Name)
		ctx = context.WithValue(ctx, ContextSubject, claims.Subject)
		ctx = context.WithValue(ctx, ContextVerified, claims.Verified)

		// Set security headers (consistent with WithAuth).
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")

		next(w, r.WithContext(ctx))
	}
}

// writeJSONRPCError writes a JSON-RPC 2.0 error response with 401 status.
func writeJSONRPCError(w http.ResponseWriter, message string) {
	resp := mcp.NewErrorResponse(nil, mcp.ErrorCodeInvalidRequest, message)
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write(data)
}
