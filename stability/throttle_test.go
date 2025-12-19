package stability_test

import (
	"context"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func Test_Throttle_With_MaxTokensReached_Should_ReturnTooManyCalls(t *testing.T) {
	// Arrange
	fn := stability.Throttle[int](mockSucceedsTimes(10), 3, 1, 100*time.Millisecond)

	// Act
	fn(context.Background(), 42)
	fn(context.Background(), 42)
	fn(context.Background(), 42)
	_, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be correct", err.Error(), "Too many calls")
}

func Test_Throttle_With_RefillAfterDuration_Should_AllowNewCalls(t *testing.T) {
	// Arrange
	fn := stability.Throttle[int](mockSucceedsTimes(10), 3, 1, 100*time.Millisecond)
	fn(context.Background(), 42)
	fn(context.Background(), 42)
	fn(context.Background(), 42)

	// Act
	time.Sleep(150 * time.Millisecond)
	res, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func Test_Throttle_With_SingleCall_Should_Succeed(t *testing.T) {
	// Arrange
	fn := stability.Throttle[int](mockAlwaysSucceeds(), 3, 3, 100*time.Millisecond)

	// Act
	res, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "res must be correct", res, 42)
}

func Test_Throttle_With_TwoCalls_Should_ReturnError(t *testing.T) {
	// Arrange
	fn := stability.Throttle[int](mockSucceedsTimes(1), 3, 3, 100*time.Millisecond)

	// Act
	fn(context.Background(), 42)
	_, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be correct", err.Error(), "error")
}
