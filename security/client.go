package security

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"time"
)

// TLSClientConfig creates and returns a *tls.Config configured for mutual TLS authentication.
// It loads client specific certificates and adds server specific root CA certificates.
func TLSClientConfig(certFile, keyFile, caFile string) *tls.Config {
	// Load client specific certificates for mutual TLS authentication.
	var clientCerts []tls.Certificate
	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err == nil {
			clientCerts = append(clientCerts, cert)
		}
	}

	// Add server specific root CA certificates to the pool of trusted certificates.
	var rootCAs *x509.CertPool
	if caFile != "" {
		caPool := x509.NewCertPool()
		caCert, err := os.ReadFile(caFile)
		if err == nil {
			caPool.AppendCertsFromPEM(caCert)
			rootCAs = caPool
		}
	}

	return &tls.Config{
		Certificates: clientCerts,
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
		// Set the minimum supported TLS version to 1.2 to avoid insecure older versions.
		MinVersion: tls.VersionTLS12,
		// Prefer server-selected cipher suites over client preferences.
		PreferServerCipherSuites: true,
		// Use the system's root CA certificates to verify server certificates by default (nil)
		// or provide a custom pool of trusted certificates.
		RootCAs: rootCAs,
	}
}

// NewClient creates and returns a new *http.Client with a default timeout of 5 seconds.
// The timeout can be adjusted by setting the CLIENT_TIMEOUT environment variable.
func NewClient() *http.Client {
	// Create a new *http.Client with a default timeout of 5 seconds.
	return &http.Client{
		Timeout: ParseDuration("CLIENT_TIMEOUT", 5*time.Second),
	}
}

// NewClientWithTLS creates and returns a new *http.Client with mutual TLS authentication.
func NewClientWithTLS(certFile, keyFile, caFile string) *http.Client {
	client := NewClient()
	client.Transport = &http.Transport{
		TLSClientConfig: TLSClientConfig(certFile, keyFile, caFile),
	}
	return client
}
