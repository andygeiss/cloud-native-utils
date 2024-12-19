package stability_test

import (
	"context"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native/utils/assert"
	"github.com/andygeiss/cloud-native/utils/stability"
)

func TestRetry_Succeeds(t *testing.T) {
	fn := stability.Retry[int](mockAlwaysSucceeds(), 3, 10*time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func TestRetry_Suceeds_With_Retries(t *testing.T) {
	fn := stability.Retry[int](mockFailsTimes(2), 3, 10*time.Millisecond)
	res, err := fn(context.Background(), 42)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}
