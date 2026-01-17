package web_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/mcp"
	"github.com/andygeiss/cloud-native-utils/web"
)

func mockToolHandler() mcp.ToolHandler {
	return func(_ context.Context, _ mcp.ToolsCallParams) (mcp.ToolsCallResult, error) {
		return mcp.ToolsCallResult{
			Content: []mcp.ContentBlock{mcp.NewTextContent("mock result")},
		}, nil
	}
}

func parseSecondResponse(body string) mcp.Response {
	lines := strings.Split(strings.TrimSpace(body), "\n")
	var resp mcp.Response
	if len(lines) >= 2 {
		_ = json.Unmarshal([]byte(lines[1]), &resp)
	}
	return resp
}

func sendMCPRequest(handler *web.MCPHandler, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler.Handler()(rec, req)
	return rec
}

func Test_NewMCPHandler_Should_ReturnNonNilHandler(t *testing.T) {
	// Arrange
	server := mcp.NewServer("test-server", "1.0.0")

	// Act
	handler := web.NewMCPHandler(server)

	// Assert
	assert.That(t, "handler must not be nil", handler != nil, true)
}

func Test_MCPHandler_Handler_Should_ReturnHandlerFunc(t *testing.T) {
	// Arrange
	server := mcp.NewServer("test-server", "1.0.0")
	handler := web.NewMCPHandler(server)

	// Act
	handlerFunc := handler.Handler()

	// Assert
	assert.That(t, "handler func must not be nil", handlerFunc != nil, true)
}

func Test_MCPHandler_With_InitializeRequest_Should_ReturnCapabilities(t *testing.T) {
	// Arrange
	server := mcp.NewServer("test-server", "1.0.0")
	handler := web.NewMCPHandler(server)
	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}`
	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	rec := httptest.NewRecorder()

	// Act
	handler.Handler()(rec, req)

	// Assert
	assert.That(t, "status must be 200", rec.Code, http.StatusOK)
	assert.That(t, "content-type must be application/json", rec.Header().Get("Content-Type"), "application/json")
	var resp mcp.Response
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.That(t, "jsonrpc must be 2.0", resp.JSONRPC, "2.0")
	assert.That(t, "error must be nil", resp.Error == nil, true)
}

func Test_MCPHandler_With_ToolsList_Should_ReturnRegisteredTools(t *testing.T) {
	// Arrange
	server := mcp.NewServer("test-server", "1.0.0")
	tool := mcp.NewTool("test-tool", "A test tool", mcp.NewObjectSchema(nil, nil), mockToolHandler())
	server.RegisterTool(tool)
	handler := web.NewMCPHandler(server)
	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}` + "\n" +
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}` + "\n"

	// Act
	rec := sendMCPRequest(handler, body)

	// Assert
	assert.That(t, "status must be 200", rec.Code, http.StatusOK)
	resp := parseSecondResponse(rec.Body.String())
	assert.That(t, "error must be nil", resp.Error == nil, true)
	result, ok := resp.Result.(map[string]any)
	assert.That(t, "result must be map", ok, true)
	tools, ok := result["tools"].([]any)
	assert.That(t, "tools must be array", ok, true)
	assert.That(t, "tools must have 1 element", len(tools), 1)
}

func Test_MCPHandler_With_ToolsCall_Should_ExecuteHandler(t *testing.T) {
	// Arrange
	server := mcp.NewServer("test-server", "1.0.0")
	tool := mcp.NewTool("test-tool", "A test tool", mcp.NewObjectSchema(nil, nil), mockToolHandler())
	server.RegisterTool(tool)
	handler := web.NewMCPHandler(server)
	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}` + "\n" +
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"test-tool","arguments":{}}}` + "\n"

	// Act
	rec := sendMCPRequest(handler, body)

	// Assert
	assert.That(t, "status must be 200", rec.Code, http.StatusOK)
	resp := parseSecondResponse(rec.Body.String())
	assert.That(t, "error must be nil", resp.Error == nil, true)
	result, ok := resp.Result.(map[string]any)
	assert.That(t, "result must be map", ok, true)
	content, ok := result["content"].([]any)
	assert.That(t, "content must be array", ok, true)
	assert.That(t, "content must have 1 element", len(content), 1)
}

func Test_MCPHandler_With_EmptyBody_Should_NotPanic(t *testing.T) {
	// Arrange
	server := mcp.NewServer("test-server", "1.0.0")
	handler := web.NewMCPHandler(server)
	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(""))
	rec := httptest.NewRecorder()

	// Act & Assert - should not panic
	handler.Handler()(rec, req)
	assert.That(t, "status must be 200", rec.Code, http.StatusOK)
}
