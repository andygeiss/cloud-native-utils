package stability_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func TestRetry_Succeeds(t *testing.T) {
	fn := stability.Retry(mockAlwaysSucceeds(), 3, 10*time.Millisecond)
	err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
}

func TestRetry_Succeeds_With_Retries(t *testing.T) {
	fn := stability.Retry(mockFailsTimes(2), 3, 10*time.Millisecond)
	err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
}

func TestRetry_Fails_After_Max_Retries(t *testing.T) {
	fn := stability.Retry(mockAlwaysFails(), 3, 10*time.Millisecond)
	err := fn(context.Background(), 42)
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must match expected", err.Error(), "error")
}

func TestRetry_Stops_On_Context_Cancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer cancel()
	fn := stability.Retry(mockSlow(20*time.Millisecond), 3, 10*time.Millisecond)
	err := fn(ctx, 42)
	assert.That(t, "err must match context error", err, context.DeadlineExceeded)
}

func TestRetry2_Succeeds(t *testing.T) {
	fn := stability.Retry2(mockAlwaysSucceeds2(), 3, 10*time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func TestRetry2_Succeeds_With_Retries(t *testing.T) {
	fn := stability.Retry2(mockFailsTimes2(2), 3, 10*time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func TestRetry2_Fails_After_Max_Retries(t *testing.T) {
	fn := stability.Retry2(mockAlwaysFails2(), 3, 10*time.Millisecond)
	_, err := fn(context.Background(), 42)
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must match expected", err.Error(), "error")
}

func TestRetry2_Stops_On_Context_Cancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer cancel()
	fn := stability.Retry2(mockSlow2(20*time.Millisecond), 3, 10*time.Millisecond)
	_, err := fn(ctx, 42)
	assert.That(t, "err must match context error", err, context.DeadlineExceeded)
}

func TestRetry2_Concurrent(t *testing.T) {
	const goroutines = 10
	numberOfCalls := int32(0)
	fn := stability.Retry2(func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 3, 10*time.Millisecond)

	results := make(chan int, goroutines)
	errs := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			res, err := fn(context.Background(), 42)
			results <- res
			errs <- err
		}()
	}

	for i := 0; i < goroutines; i++ {
		err := <-errs
		res := <-results
		assert.That(t, "err must be nil", err == nil, true)
		assert.That(t, "result must be correct", res, 42)
	}

	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be correct", n, int32(goroutines))
}
