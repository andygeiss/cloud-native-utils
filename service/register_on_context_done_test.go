package service_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/service"
)

func Test_RegisterOnContextDone_With_CancelledContext_Should_CallFunction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	called := false

	// Act
	service.RegisterOnContextDone(ctx, func() {
		called = true
		wg.Done()
	})
	cancel()

	// Assert
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		assert.That(t, "function must be called", called, true)
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for function to be called")
	}
}

func Test_RegisterOnContextDone_With_ValidContext_Should_NotCallFunctionImmediately(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	called := false

	// Act
	service.RegisterOnContextDone(ctx, func() {
		called = true
	})

	// Assert
	time.Sleep(100 * time.Millisecond)
	assert.That(t, "function must not be called immediately", called, false)
}
