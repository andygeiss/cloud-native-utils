package service_test

import (
	"cloud-native/service"
	"context"
	"testing"
	"time"
)

func TestWrap_Succeeds(t *testing.T) {
	succeedsFn := func() (*int, error) {
		value := 42
		return &value, nil
	}
	wrappedFn := service.Wrap[int](succeedsFn)
	res, err := wrappedFn(context.Background())
	if err != nil {
		t.Fatalf("error must be nil, but got %v", err)
	}
	if res == nil {
		t.Fatal("result must be not nil")
	}
	if *res != 42 {
		t.Fatalf("result must be %d, but got %d", 42, *res)
	}
}

func TestWrap_Fails_With_Timeout(t *testing.T) {
	failsFn := func() (*int, error) {
		value := 42
		time.Sleep(100 * time.Millisecond)
		return &value, nil
	}
	wrappedFn := service.Wrap[int](failsFn)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := wrappedFn(ctx)
	if err == nil {
		t.Fatal("error must be not nil")
	}
	if err.Error() != "context deadline exceeded" {
		t.Fatalf("error must be correct, but got %v", err)
	}
}
