package mcp

import "encoding/json"

// Request represents a JSON-RPC 2.0 request.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response.
type Response struct {
	Result  any             `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
}

// ResponseError represents a JSON-RPC 2.0 error object.
type ResponseError struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewResponse creates a successful JSON-RPC 2.0 response.
func NewResponse(id json.RawMessage, result any) Response {
	return Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// NewErrorResponse creates a JSON-RPC 2.0 error response.
func NewErrorResponse(id json.RawMessage, code int, message string) Response {
	return Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &ResponseError{
			Code:    code,
			Message: message,
		},
	}
}
