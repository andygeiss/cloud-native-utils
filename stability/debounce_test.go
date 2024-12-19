package stability_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func TestDebounce_Call_Once_Succeeds(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Debounce(func(ctx context.Context, in int) error {
		atomic.AddInt32(&numberOfCalls, 1)
		return nil
	}, 50*time.Millisecond)

	err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 1", n, int32(1))
}

func TestDebounce_Call_Twice_Returns_Last_Result(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Debounce(func(ctx context.Context, in int) error {
		atomic.AddInt32(&numberOfCalls, 1)
		return nil
	}, 50*time.Millisecond)

	_ = fn(context.Background(), 42)
	_ = fn(context.Background(), 42)

	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 1", n, int32(1))
}

func TestDebounce_Call_Twice_With_Delay(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Debounce(func(ctx context.Context, in int) error {
		atomic.AddInt32(&numberOfCalls, 1)
		return nil
	}, 50*time.Millisecond)

	_ = fn(context.Background(), 42)
	time.Sleep(100 * time.Millisecond)
	_ = fn(context.Background(), 42)

	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 2", n, int32(2))
}

func TestDebounce_Call_Concurrent(t *testing.T) {
	const goroutines = 10
	numberOfCalls := int32(0)
	fn := stability.Debounce(func(ctx context.Context, in int) error {
		atomic.AddInt32(&numberOfCalls, 1)
		return nil
	}, 50*time.Millisecond)

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

	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 1", n, int32(1))
}

func TestDebounce2_Call_Once_Succeeds(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Debounce2(func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 50*time.Millisecond)

	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 1", n, int32(1))
}

func TestDebounce2_Call_Twice_Returns_Last_Result(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Debounce2(func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 50*time.Millisecond)

	_, _ = fn(context.Background(), 42)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 1", n, int32(1))
}

func TestDebounce2_Call_Twice_With_Delay(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Debounce2(func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 50*time.Millisecond)

	_, _ = fn(context.Background(), 42)
	time.Sleep(100 * time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 2", n, int32(2))
}

func TestDebounce2_Call_Concurrent(t *testing.T) {
	const goroutines = 10
	numberOfCalls := int32(0)
	fn := stability.Debounce2(func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 50*time.Millisecond)

	errs := make(chan error, goroutines)
	results := make(chan int, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			res, err := fn(context.Background(), 42)
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

	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 1", n, int32(1))
}
