package stability_test

import (
	"context"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func TestRetry_Succeeds(t *testing.T) {
	fn := stability.Retry(mockAlwaysSucceeds(), 3, 10*time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func TestRetry_Succeeds_With_Retries(t *testing.T) {
	fn := stability.Retry(mockFailsTimes(2), 3, 10*time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func TestRetry_Fails(t *testing.T) {
	fn := stability.Retry(mockAlwaysFails(), 3, 10*time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "result must be correct", res, 0)
}
