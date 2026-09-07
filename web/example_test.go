package web_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andygeiss/cloud-native-utils/web"
)

// This example has no Output comment: it builds a server rather than running one.
func ExampleNewServeMux() {
	ctx := context.Background()

	// NewServeMux wires /static, /health, /liveness, /readiness and the OIDC
	// endpoints, then hands back the sessions store and identity provider it
	// built. Both belong to the caller from here on.
	mux, sessions, idp := web.NewServeMux(ctx, efs)

	// Sessions guard a browser route.
	mux.HandleFunc("GET /protected", web.WithAuth(sessions, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello %v", r.Context().Value(web.ContextEmail))
	}))

	// A bearer token guards an API route.
	mux.HandleFunc("POST /mcp", web.WithBearerAuth(idp.Verifier(), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// NewServer applies the timeouts and TLS settings; the caller starts it.
	server := web.NewServer(mux)
	_ = server.ListenAndServeTLS("server.crt", "server.key")
}
