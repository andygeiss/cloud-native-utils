package stability_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func Test_Debounce_With_ConcurrentCalls_Should_CallOnce(t *testing.T) {
	// Arrange
	const goroutines = 10
	numberOfCalls := int32(0)
	fn := stability.Debounce[int, int](func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 50*time.Millisecond)
	errs := make(chan error, goroutines)

	// Act
	for range goroutines {
		go func() {
			_, err := fn(context.Background(), 42)
			errs <- err
		}()
	}
	for range goroutines {
		err := <-errs
		assert.That(t, "error must be nil", err == nil, true)
	}

	// Assert
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be correct", n, int32(1))
}

func Test_Debounce_With_SingleCall_Should_Succeed(t *testing.T) {
	// Arrange
	fn := stability.Debounce[int](mockAlwaysSucceeds(), 50*time.Millisecond)

	// Act
	res, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func Test_Debounce_With_TwoCallsAfterDuration_Should_ReturnError(t *testing.T) {
	// Arrange
	fn := stability.Debounce[int](mockSucceedsTimes(1), 50*time.Millisecond)

	// Act
	_, _ = fn(context.Background(), 42)
	time.Sleep(100 * time.Millisecond)
	_, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "error must be correct", err.Error(), "error")
}

func Test_Debounce_With_TwoCalls_Should_ReturnLastResult(t *testing.T) {
	// Arrange
	fn := stability.Debounce[int](mockSucceedsTimes(2), 50*time.Millisecond)

	// Act
	_, _ = fn(context.Background(), 21)
	res, err := fn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "result must be correct", res, 42)
}
