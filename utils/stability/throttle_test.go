package stability_test

import (
	"context"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native/utils/assert"
	"github.com/andygeiss/cloud-native/utils/stability"
)

func TestThrottle_Call_Once_Succeeds(t *testing.T) {
	fn := stability.Throttle[int](mockAlwaysSucceeds(), 3, 3, 100*time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "res must be correct", res, 42)
}

func TestThrottle_Call_Twice_Returns_Error(t *testing.T) {
	fn := stability.Throttle[int](mockSucceedsTimes(1), 3, 3, 100*time.Millisecond)
	fn(context.Background(), 42)
	_, err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err.Error(), "error")
}

func TestThrottle_Call_Reaches_Max_Tokens(t *testing.T) {
	fn := stability.Throttle[int](mockSucceedsTimes(10), 3, 1, 100*time.Millisecond)
	fn(context.Background(), 42)
	fn(context.Background(), 42)
	fn(context.Background(), 42)
	_, err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err.Error(), "too many calls")
}

func TestThrottle_Call_Refills_After_Duration(t *testing.T) {
	fn := stability.Throttle[int](mockSucceedsTimes(10), 3, 1, 100*time.Millisecond)
	fn(context.Background(), 42)
	fn(context.Background(), 42)
	fn(context.Background(), 42)
	time.Sleep(150 * time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}
