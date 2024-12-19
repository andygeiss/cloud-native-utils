package stability_test

import (
	"context"
	"testing"

	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func TestBreaker_Call_Once_Succeeds(t *testing.T) {
	fn := stability.Breaker[int](mockAlwaysSucceeds(), 3)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func TestBreaker_Call_Once_Fails(t *testing.T) {
	fn := stability.Breaker[int](mockAlwaysFails(), 3)
	_, err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err.Error(), "error")
}

func TestBreaker_Call_Threshold(t *testing.T) {
	threshold := 3
	fn := stability.Breaker[int](mockAlwaysFails(), threshold)
	for i := 0; i < threshold; i++ {
		_, _ = fn(context.Background(), 42)
	}
	_, err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err.Error(), stability.ErrorBreakerServiceUnavailable.Error())
}

func TestBreaker_Call_Concurrent(t *testing.T) {
	const goroutines = 10
	threshold := 3
	fn := stability.Breaker[int](mockAlwaysSucceeds(), threshold)
	errs := make(chan error, goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			_, err := fn(context.Background(), 42)
			errs <- err
		}()
	}
	for i := 0; i < goroutines; i++ {
		err := <-errs
		assert.That(t, "err must be nil", err == nil, true)
	}
}

func TestBreaker_Recovery(t *testing.T) {
	threshold := 3
	fn := stability.Breaker[int](mockFailsTimes(threshold), threshold)

	// Exceed the failure threshold to trip the breaker
	for i := 0; i < threshold; i++ {
		_, _ = fn(context.Background(), 42)
	}
	_, err := fn(context.Background(), 42)
	assert.That(t, "err must be ServiceUnavailable", err.Error(), stability.ErrorBreakerServiceUnavailable.Error())

	// Wait for the breaker to recover
	time.Sleep(2 * time.Second)

	// Call again after recovery period
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil after recovery", err == nil, true)
	assert.That(t, "result must be correct after recovery", res, 42)
}
