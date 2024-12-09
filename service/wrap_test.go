package service_test

import (
	"cloud-native/service"
	"context"
	"testing"
	"time"
)

func TestWrap_Succeeds(t *testing.T) {
	succeedsFn := func(x int) (int, error) {
		return x, nil
	}
	wrappedFn := service.Wrap[int](succeedsFn)
	res, err := wrappedFn(context.Background(), 42)
	if err != nil {
		t.Fatalf("error must be nil, but got %v", err)
	}
	if res != 42 {
		t.Fatalf("result must be %d, but got %d", 42, res)
	}
}

func TestWrap_Fails_With_Timeout(t *testing.T) {
	failsFn := func(x int) (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	}
	wrappedFn := service.Wrap[int](failsFn)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := wrappedFn(ctx, 42)
	if err == nil {
		t.Fatal("error must be not nil")
	}
	if err.Error() != "context deadline exceeded" {
		t.Fatalf("error must be correct, but got %v", err)
	}
}
