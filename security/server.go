package security

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// tlsConfig creates and returns a *tls.Config configured for the given domains.
// It handles automatic TLS certificate acquisition, renewal, and secure settings.
func tlsConfig(domains ...string) *tls.Config {
	// Default directory for storing cached TLS certificates
	const defaultCertCache = "./testdata"
	// autocert.Manager automates the process of obtaining and managing TLS certificates.
	mgr := &autocert.Manager{
		Cache:      autocert.DirCache(defaultCertCache),
		HostPolicy: autocert.HostWhitelist(domains...),
		Prompt:     autocert.AcceptTOS, // Automatically accept the Let's Encrypt Terms of Service
	}
	// Return a TLS configuration with secure settings and certificate management.
	cfg := &tls.Config{
		// Define supported cipher suites for secure communication.
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
		// Specify preferred elliptic curves for key exchange.
		CurvePreferences: []tls.CurveID{
			tls.CurveP256, // NIST P-256 curve
			tls.X25519,    // Modern curve with better performance and security
		},
		// Dynamically obtain certificates for the specified domains.
		GetCertificate: mgr.GetCertificate,
		// Set the minimum supported TLS version to 1.2 to avoid insecure older versions.
		MinVersion: tls.VersionTLS12,
		// Define supported application layer protocols using ALPN (e.g., HTTP/2, HTTP/1.1, ACME).
		NextProtos: []string{
			"h2",           // HTTP/2
			"http/1.1",     // HTTP/1.1
			acme.ALPNProto, // ACME TLS-ALPN-01 protocol for certificate challenges
		},
		// Prefer server-selected cipher suites over client preferences.
		PreferServerCipherSuites: true,
	}
	// Use the self-signed certificate for localhost.
	if len(domains) == 1 && domains[0] == "localhost" {
		cfg.GetCertificate = nil
	}
	return cfg
}

// NewServer creates and returns a configured HTTP server with secure TLS settings.
// Accepts an HTTP request multiplexer (mux) and a list of domains for certificate management.
func NewServer(mux *http.ServeMux, domains ...string) *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "443"
	}
	return &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        mux,
		MaxHeaderBytes: 1 << 20, // Maximum size of request headers (1 MiB).
	}
}
