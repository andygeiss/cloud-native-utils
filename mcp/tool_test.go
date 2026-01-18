package mcp_test

import (
	"context"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/mcp"
)

func Test_NewTool_With_ValidParams_Should_ReturnToolWithDefinition(t *testing.T) {
	// Arrange
	schema := mcp.NewObjectSchema(map[string]mcp.Property{
		"name": mcp.NewStringProperty("The name"),
	}, []string{"name"})
	handler := func(ctx context.Context, params mcp.ToolsCallParams) (mcp.ToolsCallResult, error) {
		return mcp.ToolsCallResult{}, nil
	}

	// Act
	tool := mcp.NewTool("test-tool", "A test description", schema, handler)

	// Assert
	assert.That(t, "tool name must be correct", tool.Definition.Name, "test-tool")
	assert.That(t, "tool description must be correct", tool.Definition.Description, "A test description")
	assert.That(t, "tool handler must not be nil", tool.Handler != nil, true)
}

func Test_NewObjectSchema_With_Properties_Should_ReturnObjectSchema(t *testing.T) {
	// Arrange
	properties := map[string]mcp.Property{
		"name": mcp.NewStringProperty("The name"),
		"age":  mcp.NewNumberProperty("The age"),
	}
	required := []string{"name"}

	// Act
	schema := mcp.NewObjectSchema(properties, required)

	// Assert
	assert.That(t, "schema type must be object", schema.Type, "object")
	assert.That(t, "schema properties count must be 2", len(schema.Properties), 2)
	assert.That(t, "schema required count must be 1", len(schema.Required), 1)
	assert.That(t, "schema required[0] must be name", schema.Required[0], "name")
}

func Test_NewObjectSchema_With_NilProperties_Should_ReturnEmptySchema(t *testing.T) {
	// Act
	schema := mcp.NewObjectSchema(nil, nil)

	// Assert
	assert.That(t, "schema type must be object", schema.Type, "object")
	assert.That(t, "schema properties must be nil", schema.Properties == nil, true)
	assert.That(t, "schema required must be nil", schema.Required == nil, true)
}

func Test_NewStringProperty_With_Description_Should_ReturnStringProperty(t *testing.T) {
	// Act
	prop := mcp.NewStringProperty("A string description")

	// Assert
	assert.That(t, "property type must be string", prop.Type, "string")
	assert.That(t, "property description must be correct", prop.Description, "A string description")
}

func Test_NewNumberProperty_With_Description_Should_ReturnNumberProperty(t *testing.T) {
	// Act
	prop := mcp.NewNumberProperty("A number description")

	// Assert
	assert.That(t, "property type must be number", prop.Type, "number")
	assert.That(t, "property description must be correct", prop.Description, "A number description")
}

func Test_NewBooleanProperty_With_Description_Should_ReturnBooleanProperty(t *testing.T) {
	// Act
	prop := mcp.NewBooleanProperty("A boolean description")

	// Assert
	assert.That(t, "property type must be boolean", prop.Type, "boolean")
	assert.That(t, "property description must be correct", prop.Description, "A boolean description")
}

func Test_NewStringProperty_With_EmptyDescription_Should_ReturnPropertyWithEmptyDescription(t *testing.T) {
	// Act
	prop := mcp.NewStringProperty("")

	// Assert
	assert.That(t, "property type must be string", prop.Type, "string")
	assert.That(t, "property description must be empty", prop.Description, "")
}
