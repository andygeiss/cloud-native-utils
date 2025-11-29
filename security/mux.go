package security

import (
	"context"
	"embed"
	"net/http"

	"github.com/andygeiss/cloud-native-utils/efficiency"
)

// NewServeMux creates a new mux with the liveness check endpoint (/liveness)
// and the readiness check endpoint (/readiness).
// The mux is returned along with a new ServerSessions instance.
func NewServeMux(ctx context.Context, efs embed.FS) (mux *http.ServeMux, serverSessions *ServerSessions) {
	// Create a new mux with liveness and readyness endpoint.
	mux = http.NewServeMux()

	// Create an in-memory store for the server sessions.
	serverSessions = NewServerSessions()

	// Embed the assets into the mux.
	mux.Handle("/assets/", efficiency.WithCompression(http.FileServer(http.FS(efs))))

	// Add OpenID Connect endpoints to the mux.
	mux.Handle("GET /auth/callback", IdentityProvider.Callback(serverSessions))
	mux.Handle("GET /auth/login", IdentityProvider.Login())
	mux.Handle("GET /auth/logout", IdentityProvider.Logout(serverSessions))

	// Add a liveness check endpoint to the mux.
	mux.HandleFunc("GET /liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Add a readiness check endpoint to the mux.
	mux.HandleFunc("GET /readiness", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusServiceUnavailable)
		default:
			w.WriteHeader(http.StatusOK)
		}
	})

	return mux, serverSessions
}
