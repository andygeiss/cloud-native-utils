package mcp_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/andygeiss/cloud-native-utils/mcp"
)

// FuzzServerServe feeds arbitrary bytes down the pipe an MCP client owns. The
// server must answer or reject every line without panicking.
func FuzzServerServe(f *testing.F) {
	f.Add(`{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n")
	f.Add(`{"jsonrpc":"2.0","id":2,"method":"tools/list"}` + "\n")
	f.Add(`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"echo","arguments":{"text":"hi"}}}` + "\n")
	f.Add(`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"missing"}}` + "\n")
	f.Add(`{"jsonrpc":"2.0","method":"nope"}` + "\n")
	f.Add("not json\n")
	f.Add("{\n")
	f.Add("")
	f.Add("\n\n\n")

	f.Fuzz(func(t *testing.T, input string) {
		server := mcp.NewServerWithIO("fuzz", "1.0.0", strings.NewReader(input), io.Discard)
		server.RegisterTool(mcp.NewTool(
			"echo",
			"Echo the text back",
			mcp.NewObjectSchema(
				map[string]mcp.Property{"text": mcp.NewStringProperty("text to echo")},
				[]string{"text"},
			),
			func(ctx context.Context, params mcp.ToolsCallParams) (mcp.ToolsCallResult, error) {
				text, _ := params.Arguments["text"].(string)
				return mcp.ToolsCallResult{Content: []mcp.ContentBlock{mcp.NewTextContent(text)}}, nil
			},
		))

		// Serve returns nil once the reader is exhausted. A panic is the bug.
		_ = server.Serve(context.Background())
	})
}
