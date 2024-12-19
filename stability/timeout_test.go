package stability_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func TestTimeout_SuccessfulFunction(t *testing.T) {
	fn := stability.Timeout(mockAlwaysSucceeds(), 1*time.Second)
	ctx := context.Background()
	result, err := fn(ctx, 5)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", result, 42)
}

func TestTimeout_TimeoutFunction(t *testing.T) {
	fn := stability.Timeout(mockSlow(500*time.Millisecond), 200*time.Millisecond)
	ctx := context.Background()
	_, err := fn(ctx, 5)
	assert.That(t, "err must be correct", errors.Is(err, context.DeadlineExceeded), true)
}

func TestTimeout_ContextCancellation(t *testing.T) {
	fn := stability.Timeout(mockSlow(250*time.Millisecond), 500*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately
	_, err := fn(ctx, 5)
	assert.That(t, "err must be correct", errors.Is(err, context.Canceled), true)
}

func TestTimeout_NoTimeout(t *testing.T) {
	fn := stability.Timeout(mockSlow(250*time.Millisecond), 300*time.Millisecond)
	ctx := context.Background()
	result, err := fn(ctx, 5)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", result, 10)
}
