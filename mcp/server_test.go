package mcp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/mcp"
)

func runServer(input string, registerTool *mcp.Tool) ([]string, error) {
	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}
	server := mcp.NewServerWithIO("test-server", "1.0.0", reader, writer)
	if registerTool != nil {
		server.RegisterTool(*registerTool)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err := server.Serve(ctx)
	return strings.Split(strings.TrimSpace(writer.String()), "\n"), err
}

func parseResponse(line string) mcp.Response {
	var resp mcp.Response
	_ = json.Unmarshal([]byte(line), &resp)
	return resp
}

func Test_Server_With_Initialize_Should_ReturnCapabilities(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}` + "\n"
	lines, _ := runServer(input, nil)
	resp := parseResponse(lines[0])
	assert.That(t, "jsonrpc must be 2.0", resp.JSONRPC, "2.0")
	assert.That(t, "error must be nil", resp.Error == nil, true)
	assert.That(t, "result must not be nil", resp.Result != nil, true)
}

func Test_Server_With_ToolsList_Should_ReturnRegisteredTools(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}` + "\n" +
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}` + "\n"
	tool := mcp.NewTool("test-tool", "A test tool", mcp.NewObjectSchema(nil, nil), mockToolHandler())
	lines, _ := runServer(input, &tool)
	assert.That(t, "should have 2 responses", len(lines), 2)
	resp := parseResponse(lines[1])
	assert.That(t, "error must be nil", resp.Error == nil, true)
	result, ok := resp.Result.(map[string]any)
	assert.That(t, "result must be map", ok, true)
	tools, ok := result["tools"].([]any)
	assert.That(t, "tools must be array", ok, true)
	assert.That(t, "tools must have 1 element", len(tools), 1)
}

func Test_Server_With_ToolsCall_Should_ExecuteHandler(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}` + "\n" +
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"test-tool","arguments":{}}}` + "\n"
	tool := mcp.NewTool("test-tool", "A test tool", mcp.NewObjectSchema(nil, nil), mockToolHandler())
	lines, _ := runServer(input, &tool)
	assert.That(t, "should have 2 responses", len(lines), 2)
	resp := parseResponse(lines[1])
	assert.That(t, "error must be nil", resp.Error == nil, true)
	result, ok := resp.Result.(map[string]any)
	assert.That(t, "result must be map", ok, true)
	content, ok := result["content"].([]any)
	assert.That(t, "content must be array", ok, true)
	assert.That(t, "content must have 1 element", len(content), 1)
}

func Test_Server_With_UnknownMethod_Should_ReturnMethodNotFound(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":1,"method":"unknown/method"}` + "\n"
	lines, _ := runServer(input, nil)
	resp := parseResponse(lines[0])
	assert.That(t, "error must not be nil", resp.Error != nil, true)
	assert.That(t, "error code must be -32601", resp.Error.Code, mcp.ErrorCodeMethodNotFound)
}

func Test_Server_With_ToolsCallBeforeInit_Should_ReturnError(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"test-tool","arguments":{}}}` + "\n"
	tool := mcp.NewTool("test-tool", "A test tool", mcp.NewObjectSchema(nil, nil), mockToolHandler())
	lines, _ := runServer(input, &tool)
	resp := parseResponse(lines[0])
	assert.That(t, "error must not be nil", resp.Error != nil, true)
	assert.That(t, "error code must be -32603", resp.Error.Code, mcp.ErrorCodeInternal)
}

func Test_Server_With_ToolHandlerError_Should_ReturnIsErrorTrue(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}` + "\n" +
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"error-tool","arguments":{}}}` + "\n"
	tool := mcp.NewTool("error-tool", "A tool that errors", mcp.NewObjectSchema(nil, nil), mockToolHandlerWithError())
	lines, _ := runServer(input, &tool)
	assert.That(t, "should have 2 responses", len(lines), 2)
	resp := parseResponse(lines[1])
	assert.That(t, "error must be nil", resp.Error == nil, true)
	result, ok := resp.Result.(map[string]any)
	assert.That(t, "result must be map", ok, true)
	isError, ok := result["isError"].(bool)
	assert.That(t, "isError must be bool", ok, true)
	assert.That(t, "isError must be true", isError, true)
}

func Test_Server_With_ToolNotFound_Should_ReturnError(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}` + "\n" +
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"nonexistent","arguments":{}}}` + "\n"
	lines, _ := runServer(input, nil)
	assert.That(t, "should have 2 responses", len(lines), 2)
	resp := parseResponse(lines[1])
	assert.That(t, "error must not be nil", resp.Error != nil, true)
	assert.That(t, "error code must be -32602", resp.Error.Code, mcp.ErrorCodeInvalidParams)
}
