package stability_test

import (
	"context"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func Test_Retry_With_AlwaysFailingFunction_Should_ReturnError(t *testing.T) {
	// Arrange
	fn := stability.Retry(mockAlwaysFails(), 3, 10*time.Millisecond)

	// Act
	res, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "result must be correct", res, 0)
}

func Test_Retry_With_SuccessAfterRetries_Should_ReturnResult(t *testing.T) {
	// Arrange
	fn := stability.Retry(mockFailsTimes(2), 3, 10*time.Millisecond)

	// Act
	res, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func Test_Retry_With_SuccessfulFunction_Should_ReturnResult(t *testing.T) {
	// Arrange
	fn := stability.Retry(mockAlwaysSucceeds(), 3, 10*time.Millisecond)

	// Act
	res, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}
