package stability_test

import (
	"cloud-native/utils/stability"
	"context"
	"errors"
	"testing"
	"time"
)

func TestTimeout_SuccessfulFunction(t *testing.T) {
	fn := stability.Timeout(mockAlwaysSucceeds(), 1*time.Second)
	ctx := context.Background()
	output, err := fn(ctx, 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output != 42 {
		t.Fatalf("expected output to be 10, got %v", output)
	}
}

func TestTimeout_TimeoutFunction(t *testing.T) {
	fn := stability.Timeout(mockSlow(500*time.Millisecond), 200*time.Millisecond)
	ctx := context.Background()
	_, err := fn(ctx, 5)
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected error to be context.DeadlineExceeded, got %v", err)
	}
}

func TestTimeout_ContextCancellation(t *testing.T) {
	fn := stability.Timeout(mockSlow(250*time.Millisecond), 500*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately
	_, err := fn(ctx, 5)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected error to be context.Canceled, got %v", err)
	}
}

func TestTimeout_NoTimeout(t *testing.T) {
	fn := stability.Timeout(mockSlow(250*time.Millisecond), 300*time.Millisecond)
	ctx := context.Background()
	output, err := fn(ctx, 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output != 10 {
		t.Fatalf("expected output to be 10, got %v", output)
	}
}
