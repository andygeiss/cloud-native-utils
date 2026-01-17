package mcp

import "errors"

// JSON-RPC 2.0 standard error codes.
const (
	ErrorCodeInternal       = -32603
	ErrorCodeInvalidParams  = -32602
	ErrorCodeInvalidRequest = -32600
	ErrorCodeMethodNotFound = -32601
	ErrorCodeParse          = -32700
)

// MCP-specific errors.
var (
	ErrMCPInvalidToolParams = errors.New("invalid tool parameters")
	ErrMCPNotInitialized    = errors.New("server not initialized")
	ErrMCPToolNotFound      = errors.New("tool not found")
	ErrMCPTransportClosed   = errors.New("transport closed")
)
