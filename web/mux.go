package web

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/andygeiss/cloud-native-utils/efficiency"
)

// NewServeMux creates a new mux with the liveness check endpoint (/liveness)
// and the readiness check endpoint (/readiness).
// The mux is returned along with a new ServerSessions instance.
func NewServeMux(ctx context.Context, efs fs.FS) (*http.ServeMux, *ServerSessions) {
	// Create a new mux with liveness and readyness endpoint.
	mux := http.NewServeMux()

	// Create an in-memory store for the server sessions.
	serverSessions := NewServerSessions()

	// Chroot into the assets directory for static files.
	staticFS, err := fs.Sub(efs, "assets")
	if err != nil {
		panic(err)
	}

	// Embed static files into the mux.
	mux.Handle("/static/", efficiency.WithCompression(http.FileServerFS(staticFS)))

	// Add OpenID Connect endpoints to the mux.
	mux.Handle("GET /auth/callback", IdentityProvider.Callback(serverSessions))
	mux.Handle("GET /auth/login", IdentityProvider.Login())
	mux.Handle("GET /auth/logout/{session_id}", IdentityProvider.Logout(serverSessions))

	// Health endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Add a liveness check endpoint to the mux.
	mux.HandleFunc("GET /liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
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
