package mcp

import (
	"github.com/andygeiss/cloud-native-utils/service"
)

// ToolHandler is a function that handles tool calls.
// It follows the cloud-native-utils Function pattern.
type ToolHandler service.Function[ToolsCallParams, ToolsCallResult]

// ToolDefinition describes a tool's metadata for tools/list.
type ToolDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputSchema InputSchema `json:"inputSchema"`
}

// InputSchema represents the JSON Schema for tool input.
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

// Property represents a JSON Schema property.
type Property struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// Tool represents a registered tool with its definition and handler.
type Tool struct {
	Handler    ToolHandler
	Definition ToolDefinition
}

// NewTool creates a new tool with the given definition and handler.
func NewTool(name, description string, schema InputSchema, handler ToolHandler) Tool {
	return Tool{
		Definition: ToolDefinition{
			Name:        name,
			Description: description,
			InputSchema: schema,
		},
		Handler: handler,
	}
}

// NewObjectSchema creates an object input schema.
func NewObjectSchema(properties map[string]Property, required []string) InputSchema {
	return InputSchema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
}

// NewStringProperty creates a string property.
func NewStringProperty(description string) Property {
	return Property{
		Type:        "string",
		Description: description,
	}
}

// NewNumberProperty creates a number property.
func NewNumberProperty(description string) Property {
	return Property{
		Type:        "number",
		Description: description,
	}
}

// NewBooleanProperty creates a boolean property.
func NewBooleanProperty(description string) Property {
	return Property{
		Type:        "boolean",
		Description: description,
	}
}
