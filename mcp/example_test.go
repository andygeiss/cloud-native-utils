package mcp_test

import (
	"context"
	"fmt"

	"github.com/andygeiss/cloud-native-utils/mcp"
)

func ExampleServer_RegisterTool() {
	server := mcp.NewServer("greeter", "1.0.0")

	server.RegisterTool(mcp.NewTool(
		"greet",
		"Greet someone by name",
		mcp.NewObjectSchema(
			map[string]mcp.Property{"name": mcp.NewStringProperty("who to greet")},
			[]string{"name"},
		),
		func(ctx context.Context, params mcp.ToolsCallParams) (mcp.ToolsCallResult, error) {
			name, _ := params.Arguments["name"].(string)
			return mcp.ToolsCallResult{
				Content: []mcp.ContentBlock{mcp.NewTextContent("Hello, " + name)},
			}, nil
		},
	))

	// server.Serve(ctx) speaks JSON-RPC over stdio. Calling the handler
	// directly is how a test reaches it.
	result, err := server.Tools()[0].Handler(context.Background(), mcp.ToolsCallParams{
		Name:      "greet",
		Arguments: map[string]any{"name": "Alice"},
	})

	fmt.Println(result.Content[0].Text, err)
	// Output: Hello, Alice <nil>
}
