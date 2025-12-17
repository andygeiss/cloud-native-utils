package extensibility_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/extensibility"
)

type TestPort interface {
	FindByID(id string) (name string, err error)
}

func Test_LoadPlugin_With_InvalidPluginPath_Should_ReturnError(t *testing.T) {
	// Arrange & Act
	_, err := extensibility.LoadPlugin[TestPort]("testdata/adapter2.so", "Adapter")

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
}

func Test_LoadPlugin_With_InvalidSymbol_Should_ReturnError(t *testing.T) {
	// Arrange & Act
	_, err := extensibility.LoadPlugin[TestPort]("testdata/adapter.so", "Adapter2")

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
}

func Test_LoadPlugin_With_ValidPlugin_Should_ReturnAdapter(t *testing.T) {
	// Arrange & Act
	adapter, err := extensibility.LoadPlugin[TestPort]("testdata/adapter.so", "Adapter")
	name, err := adapter.FindByID("1")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "adapter.FindByID must return 'Andy'", name, "Andy")
}
