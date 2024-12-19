package stability_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func TestThrottle_Call_Once_Succeeds(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Throttle(func(ctx context.Context, in int) error {
		atomic.AddInt32(&numberOfCalls, 1)
		return nil
	}, 3, 1, 100*time.Millisecond)

	err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 1", n, int32(1))
}

func TestThrottle_Call_Twice_Returns_Error(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Throttle(func(ctx context.Context, in int) error {
		atomic.AddInt32(&numberOfCalls, 1)
		return nil
	}, 1, 1, 100*time.Millisecond)

	_ = fn(context.Background(), 42)
	err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err, stability.ErrorThrottleTooManyCalls)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 1", n, int32(1))
}

func TestThrottle_Call_Reaches_Max_Tokens(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Throttle(func(ctx context.Context, in int) error {
		atomic.AddInt32(&numberOfCalls, 1)
		return nil
	}, 3, 1, 100*time.Millisecond)

	_ = fn(context.Background(), 42)
	_ = fn(context.Background(), 42)
	_ = fn(context.Background(), 42)
	err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err, stability.ErrorThrottleTooManyCalls)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 3", n, int32(3))
}

func TestThrottle_Call_Refills_After_Duration(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Throttle(func(ctx context.Context, in int) error {
		atomic.AddInt32(&numberOfCalls, 1)
		return nil
	}, 3, 1, 100*time.Millisecond)

	_ = fn(context.Background(), 42)
	_ = fn(context.Background(), 42)
	_ = fn(context.Background(), 42)
	time.Sleep(150 * time.Millisecond)
	err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 4", n, int32(4))
}

func TestThrottle2_Call_Once_Succeeds(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Throttle2(func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 3, 1, 100*time.Millisecond)

	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 1", n, int32(1))
}

func TestThrottle2_Call_Twice_Returns_Error(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Throttle2(func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 1, 1, 100*time.Millisecond)

	_, _ = fn(context.Background(), 42)
	_, err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err, stability.ErrorThrottleTooManyCalls)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 1", n, int32(1))
}

func TestThrottle2_Call_Reaches_Max_Tokens(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Throttle2(func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 3, 1, 100*time.Millisecond)

	_, _ = fn(context.Background(), 42)
	_, _ = fn(context.Background(), 42)
	_, _ = fn(context.Background(), 42)
	_, err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err, stability.ErrorThrottleTooManyCalls)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 3", n, int32(3))
}

func TestThrottle2_Call_Refills_After_Duration(t *testing.T) {
	numberOfCalls := int32(0)
	fn := stability.Throttle2(func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 3, 1, 100*time.Millisecond)

	_, _ = fn(context.Background(), 42)
	_, _ = fn(context.Background(), 42)
	_, _ = fn(context.Background(), 42)
	time.Sleep(150 * time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
	n := atomic.LoadInt32(&numberOfCalls)
	assert.That(t, "number of calls must be 4", n, int32(4))
}
