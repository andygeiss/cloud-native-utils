package stability_test

import (
	"context"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func Test_Breaker_With_CancelledContext_Should_ReturnError(t *testing.T) {
	// Arrange
	fn := stability.Breaker[int](mockAlwaysSucceeds(), 3)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	_, err := fn(ctx, 42)

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be context cancelled", err.Error(), "context canceled")
}

func Test_Breaker_With_ConcurrentCalls_Should_AllSucceed(t *testing.T) {
	// Arrange
	const goroutines = 10
	threshold := 3
	fn := stability.Breaker[int](mockAlwaysSucceeds(), threshold)
	errs := make(chan error, goroutines)

	// Act
	for i := 0; i < goroutines; i++ {
		go func() {
			_, err := fn(context.Background(), 42)
			errs <- err
		}()
	}

	// Assert
	for i := 0; i < goroutines; i++ {
		err := <-errs
		assert.That(t, "err must be nil", err == nil, true)
	}
}

func Test_Breaker_With_FailingCall_Should_ReturnError(t *testing.T) {
	// Arrange
	fn := stability.Breaker[int](mockAlwaysFails(), 3)

	// Act
	_, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be correct", err.Error(), "error")
}

func Test_Breaker_With_Recovery_Should_SucceedAfterWait(t *testing.T) {
	// Arrange
	threshold := 3
	fn := stability.Breaker[int](mockFailsTimes(threshold), threshold)

	// Act - Exceed the failure threshold to trip the breaker
	for i := 0; i < threshold; i++ {
		_, _ = fn(context.Background(), 42)
	}
	_, err := fn(context.Background(), 42)
	assert.That(t, "err must be ServiceUnavailable", err.Error(), stability.ErrorBreakerServiceUnavailable.Error())

	// Wait for the breaker to recover
	time.Sleep(2 * time.Second)

	// Act - Call again after recovery period
	res, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be nil after recovery", err == nil, true)
	assert.That(t, "result must be correct after recovery", res, 42)
}

func Test_Breaker_With_SuccessfulCall_Should_ReturnResult(t *testing.T) {
	// Arrange
	fn := stability.Breaker[int](mockAlwaysSucceeds(), 3)

	// Act
	res, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func Test_Breaker_With_ThresholdExceeded_Should_ReturnServiceUnavailable(t *testing.T) {
	// Arrange
	threshold := 3
	fn := stability.Breaker[int](mockAlwaysFails(), threshold)

	// Act
	for i := 0; i < threshold; i++ {
		_, _ = fn(context.Background(), 42)
	}
	_, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be correct", err.Error(), stability.ErrorBreakerServiceUnavailable.Error())
}

func Test_Breaker_With_TimeoutContext_Should_ReturnError(t *testing.T) {
	// Arrange
	fn := stability.Breaker[int](mockAlwaysSucceeds(), 3)
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	// Act
	_, err := fn(ctx, 42)

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
}
