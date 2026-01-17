package mcp

// InitializeParams represents the parameters for the initialize method.
type InitializeParams struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    Capabilities   `json:"capabilities"`
	ClientInfo      Implementation `json:"clientInfo"`
}

// InitializeResult represents the result of the initialize method.
type InitializeResult struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    Capabilities   `json:"capabilities"`
	ServerInfo      Implementation `json:"serverInfo"`
}

// Capabilities represents MCP capabilities.
type Capabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

// ToolsCapability represents tools capability.
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// Implementation represents server or client implementation info.
type Implementation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ToolsListResult represents the result of tools/list.
type ToolsListResult struct {
	Tools []ToolDefinition `json:"tools"`
}

// ToolsCallParams represents the parameters for tools/call.
type ToolsCallParams struct {
	Arguments map[string]any `json:"arguments,omitempty"`
	Name      string         `json:"name"`
}

// ToolsCallResult represents the result of tools/call.
type ToolsCallResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

// ContentBlock represents a content block in tool results.
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// NewTextContent creates a text content block.
func NewTextContent(text string) ContentBlock {
	return ContentBlock{
		Type: "text",
		Text: text,
	}
}
