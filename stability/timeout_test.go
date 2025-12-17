package stability_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func Test_Timeout_With_CancelledContext_Should_ReturnCancelled(t *testing.T) {
	// Arrange
	fn := stability.Timeout(mockSlow(250*time.Millisecond), 500*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	_, err := fn(ctx, 5)

	// Assert
	assert.That(t, "err must be correct", errors.Is(err, context.Canceled), true)
}

func Test_Timeout_With_FastFunction_Should_ReturnResult(t *testing.T) {
	// Arrange
	fn := stability.Timeout(mockSlow(250*time.Millisecond), 300*time.Millisecond)
	ctx := context.Background()

	// Act
	result, err := fn(ctx, 5)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", result, 10)
}

func Test_Timeout_With_SlowFunction_Should_ReturnDeadlineExceeded(t *testing.T) {
	// Arrange
	fn := stability.Timeout(mockSlow(500*time.Millisecond), 200*time.Millisecond)
	ctx := context.Background()

	// Act
	_, err := fn(ctx, 5)

	// Assert
	assert.That(t, "err must be correct", errors.Is(err, context.DeadlineExceeded), true)
}

func Test_Timeout_With_SuccessfulFunction_Should_ReturnResult(t *testing.T) {
	// Arrange
	fn := stability.Timeout(mockAlwaysSucceeds(), 1*time.Second)
	ctx := context.Background()

	// Act
	result, err := fn(ctx, 5)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", result, 42)
}
