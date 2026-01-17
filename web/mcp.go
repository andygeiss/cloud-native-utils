package web

import (
	"bytes"
	"io"
	"net/http"

	"github.com/andygeiss/cloud-native-utils/mcp"
)

// MCPHandler provides HTTP transport for MCP servers.
type MCPHandler struct {
	server *mcp.Server
}

// NewMCPHandler creates a handler that bridges HTTP to an MCP server.
func NewMCPHandler(server *mcp.Server) *MCPHandler {
	return &MCPHandler{
		server: server,
	}
}

// Handler returns an http.HandlerFunc for POST /mcp requests.
func (h *MCPHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request", http.StatusBadRequest)
			return
		}

		// MCP expects newline-delimited JSON
		if len(body) > 0 && body[len(body)-1] != '\n' {
			body = append(body, '\n')
		}

		// Create I/O buffers for MCP server
		input := bytes.NewReader(body)
		output := &bytes.Buffer{}

		// Create per-request MCP server with HTTP I/O
		server := mcp.NewServerWithIO(
			h.server.Name(),
			h.server.Version(),
			input,
			output,
		)

		// Copy tool registrations from the original server
		for _, tool := range h.server.Tools() {
			server.RegisterTool(tool)
		}

		// Serve processes the request
		_ = server.Serve(r.Context())

		// Write response
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(output.Bytes())
	}
}
