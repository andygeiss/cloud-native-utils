package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/service"
)

func Test_Wrap_With_SuccessfulFunction_Should_ReturnResult(t *testing.T) {
	// Arrange
	succeedsFn := func(x int) (int, error) {
		return x, nil
	}
	wrappedFn := service.Wrap[int](succeedsFn)

	// Act
	res, err := wrappedFn(context.Background(), 42)

	// Assert
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "result must be correct", res, 42)
}

func Test_Wrap_With_TimeoutContext_Should_ReturnError(t *testing.T) {
	// Arrange
	failsFn := func(x int) (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	}
	wrappedFn := service.Wrap[int](failsFn)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Act
	_, err := wrappedFn(ctx, 42)

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be correct", err.Error(), "context deadline exceeded")
}
