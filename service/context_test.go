package service_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/service"
)

func Test_Context_With_NoArgs_Should_ReturnCancelFunc(t *testing.T) {
	// Arrange & Act
	_, cancel := service.Context()
	defer cancel()

	// Assert
	assert.That(t, "cancel must not be nil", cancel != nil, true)
}

func Test_Context_With_NoArgs_Should_ReturnContext(t *testing.T) {
	// Arrange & Act
	ctx, cancel := service.Context()
	defer cancel()

	// Assert
	assert.That(t, "ctx must not be nil", ctx != nil, true)
}
