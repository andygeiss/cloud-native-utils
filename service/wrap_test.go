package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/service"
)

func TestWrap_Succeeds(t *testing.T) {
	succeedsFn := func(x int) error {
		return nil
	}
	wrappedFn := service.Wrap[int](succeedsFn)
	err := wrappedFn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
}

func TestWrap_Fails_With_Timeout(t *testing.T) {
	failsFn := func(x int) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	wrappedFn := service.Wrap[int](failsFn)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := wrappedFn(ctx, 42)
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be correct", err.Error(), "context deadline exceeded")
}

func TestWrap2_Succeeds(t *testing.T) {
	succeedsFn := func(x int) (int, error) {
		return x, nil
	}
	wrappedFn := service.Wrap2[int, int](succeedsFn)
	res, err := wrappedFn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func TestWrap2_Fails_With_Timeout(t *testing.T) {
	failsFn := func(x int) (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	}
	wrappedFn := service.Wrap2[int, int](failsFn)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := wrappedFn(ctx, 42)
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be correct", err.Error(), "context deadline exceeded")
}
