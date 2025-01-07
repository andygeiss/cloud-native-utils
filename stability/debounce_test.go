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
	fn := stability.Debounce[int](mockAlwaysSucceeds(), 50*time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func TestDebounce_Call_Twice_Returns_Last_Result(t *testing.T) {
	fn := stability.Debounce[int](mockSucceedsTimes(2), 50*time.Millisecond)
	fn(context.Background(), 21)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "result must be correct", res, 42)
}

func TestDebounce_Call_Twice_Returns_Error(t *testing.T) {
	fn := stability.Debounce[int](mockSucceedsTimes(1), 50*time.Millisecond)
	fn(context.Background(), 42)
	// Wait to trigger new service all next time.
	time.Sleep(100 * time.Millisecond)
	_, err := fn(context.Background(), 42)
	assert.That(t, "error must be correct", err.Error(), "error")
}

func TestDebounce_Call_Concurrent(t *testing.T) {
	const goroutines = 10
	numberOfCalls := int32(0)
	fn := stability.Debounce[int, int](func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 50*time.Millisecond)
	errs := make(chan error, 3)
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
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be correct", n, int32(1))
}
