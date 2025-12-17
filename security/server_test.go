package security_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_NewServer_With_CustomPort_Should_UseCustomPort(t *testing.T) {
	// Arrange
	os.Setenv("PORT", "9090")
	defer os.Unsetenv("PORT")
	mux := http.NewServeMux()

	// Act
	server := security.NewServer(mux)

	// Assert
	assert.That(t, "server address must use custom port", server.Addr, ":9090")
}

func Test_NewServer_With_DefaultPort_Should_UsePort8080(t *testing.T) {
	// Arrange
	os.Unsetenv("PORT")
	mux := http.NewServeMux()

	// Act
	server := security.NewServer(mux)

	// Assert
	assert.That(t, "server address must use default port", server.Addr, ":8080")
}

func Test_NewServer_With_Mux_Should_SetHandler(t *testing.T) {
	// Arrange
	os.Unsetenv("PORT")
	mux := http.NewServeMux()

	// Act
	server := security.NewServer(mux)

	// Assert
	assert.That(t, "server handler must be set", server.Handler != nil, true)
}

func Test_NewServer_With_ValidMux_Should_ReturnNonNilServer(t *testing.T) {
	// Arrange
	os.Unsetenv("PORT")
	mux := http.NewServeMux()

	// Act
	server := security.NewServer(mux)

	// Assert
	assert.That(t, "server must not be nil", server != nil, true)
}

func Test_NewServer_With_ValidMux_Should_SetTimeouts(t *testing.T) {
	// Arrange
	os.Unsetenv("PORT")
	mux := http.NewServeMux()

	// Act
	server := security.NewServer(mux)

	// Assert
	assert.That(t, "idle timeout must be set", server.IdleTimeout > 0, true)
	assert.That(t, "read timeout must be set", server.ReadTimeout > 0, true)
	assert.That(t, "write timeout must be set", server.WriteTimeout > 0, true)
}
