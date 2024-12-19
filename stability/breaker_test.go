package stability_test

import (
	"context"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

// Tests for Breaker.
func TestBreaker_Call_Once_Succeeds(t *testing.T) {
	fn := stability.Breaker(mockAlwaysSucceeds(), 3)
	err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
}

func TestBreaker_Call_Once_Fails(t *testing.T) {
	fn := stability.Breaker(mockAlwaysFails(), 3)
	err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err.Error(), "error")
}

func TestBreaker_Call_Threshold(t *testing.T) {
	threshold := 3
	fn := stability.Breaker(mockAlwaysFails(), threshold)
	for i := 0; i < threshold; i++ {
		_ = fn(context.Background(), 42)
	}
	err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err, stability.ErrorBreakerServiceUnavailable)
}

func TestBreaker_Call_Recovers(t *testing.T) {
	threshold := 3
	fn := stability.Breaker(mockFailsTimes(threshold), threshold)
	for i := 0; i < threshold; i++ {
		_ = fn(context.Background(), 42)
	}
	err := fn(context.Background(), 42)
	assert.That(t, "breaker must open", err, stability.ErrorBreakerServiceUnavailable)
	time.Sleep(2 * time.Second)
	// Wait to allow the breaker to retry after backoff.
	err = fn(context.Background(), 42)
	assert.That(t, "breaker must close on success", err == nil, true)
}

func TestBreaker_Call_Concurrent(t *testing.T) {
	const goroutines = 10
	threshold := 3
	fn := stability.Breaker(mockAlwaysSucceeds(), threshold)
	errs := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			errs <- fn(context.Background(), 42)
		}()
	}

	for i := 0; i < goroutines; i++ {
		err := <-errs
		assert.That(t, "err must be nil", err == nil, true)
	}
}

// Tests for Breaker2.
func TestBreaker2_Call_Once_Succeeds(t *testing.T) {
	fn := stability.Breaker2(mockAlwaysSucceeds2(), 3)
	res, err := fn(context.Background(), 21)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func TestBreaker2_Call_Once_Fails(t *testing.T) {
	fn := stability.Breaker2(mockAlwaysFails2(), 3)
	_, err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err.Error(), "error")
}

func TestBreaker2_Call_Threshold(t *testing.T) {
	threshold := 3
	fn := stability.Breaker2(mockAlwaysFails2(), threshold)
	for i := 0; i < threshold; i++ {
		_, _ = fn(context.Background(), 42)
	}
	_, err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err, stability.ErrorBreakerServiceUnavailable)
}

func TestBreaker2_Call_Recovers(t *testing.T) {
	threshold := 3
	fn := stability.Breaker2(mockFailsTimes2(threshold), threshold)
	for i := 0; i < threshold; i++ {
		_, _ = fn(context.Background(), 42)
	}
	_, err := fn(context.Background(), 42)
	assert.That(t, "breaker must open", err, stability.ErrorBreakerServiceUnavailable)
	time.Sleep(2 * time.Second)
	// Wait to allow the breaker to retry after backoff.
	res, err := fn(context.Background(), 42)
	assert.That(t, "breaker must close on success", err, nil)
	assert.That(t, "result must be correct", res, 42)
}

func TestBreaker2_Call_Concurrent(t *testing.T) {
	const goroutines = 10
	threshold := 3
	fn := stability.Breaker2(mockAlwaysSucceeds2(), threshold)
	errs := make(chan error, goroutines)
	results := make(chan int, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			res, err := fn(context.Background(), 21)
			errs <- err
			results <- res
		}()
	}

	for i := 0; i < goroutines; i++ {
		err := <-errs
		assert.That(t, "err must be nil", err == nil, true)

		res := <-results
		assert.That(t, "result must be correct", res, 42)
	}
}
