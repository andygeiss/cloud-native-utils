package mcp_test

import (
	"context"
	"errors"

	"github.com/andygeiss/cloud-native-utils/mcp"
)

func mockToolHandler() mcp.ToolHandler {
	return func(_ context.Context, _ mcp.ToolsCallParams) (mcp.ToolsCallResult, error) {
		return mcp.ToolsCallResult{
			Content: []mcp.ContentBlock{mcp.NewTextContent("mock result")},
		}, nil
	}
}

func mockToolHandlerWithError() mcp.ToolHandler {
	return func(_ context.Context, _ mcp.ToolsCallParams) (mcp.ToolsCallResult, error) {
		return mcp.ToolsCallResult{}, errors.New("mock error")
	}
}
