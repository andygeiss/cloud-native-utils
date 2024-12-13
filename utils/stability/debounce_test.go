package stability_test

import (
	"cloud-native/utils/stability"
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestDebounce_Call_Once_Succeeds(t *testing.T) {
	fn := stability.Debounce[int](mockAlwaysSucceeds(), 50*time.Millisecond)
	res, err := fn(context.Background(), 42)
	if err != nil {
		t.Fatalf("error must be nil, but got %v", err)
	}
	if res != 42 {
		t.Fatalf("result must be %d, but got %d", 42, res)
	}
}

func TestDebounce_Call_Twice_Returns_Last_Result(t *testing.T) {
	fn := stability.Debounce[int](mockSucceedsTimes(1), 50*time.Millisecond)
	fn(context.Background(), 42)
	res, err := fn(context.Background(), 42)
	if err != nil {
		t.Fatalf("error must be nil, but got %v", err)
	}
	if res != 42 {
		t.Fatalf("result must be %d, but got %d", 42, res)
	}
}

func TestDebounce_Call_Twice_Returns_Error(t *testing.T) {
	fn := stability.Debounce[int](mockSucceedsTimes(1), 50*time.Millisecond)
	fn(context.Background(), 42)
	// Wait to trigger new service all next time.
	time.Sleep(100 * time.Millisecond)
	_, err := fn(context.Background(), 42)
	if err == nil {
		t.Fatal("error must be not nil")
	}
	if err.Error() != "error" {
		t.Fatalf("error must be correct, but got %v", err)
	}
}

func TestDebounce_Call_Concurrent(t *testing.T) {
	const goroutines = 10
	numberOfCalls := int32(0)
	fn := stability.Debounce[int, int](func(ctx context.Context, in int) (int, error) {
		atomic.AddInt32(&numberOfCalls, 1)
		return 42, nil
	}, 50*time.Millisecond)
	errs := make(chan error, 3)
	for range goroutines {
		go func() {
			_, err := fn(context.Background(), 42)
			errs <- err
		}()
	}
	for range goroutines {
		err := <-errs
		if err != nil {
			t.Fatalf("error must be nil, but got %v", err)
		}
	}
	n := atomic.LoadInt32(&numberOfCalls)
	if n != 1 {
		t.Fatalf("number of calls must be 1, but got %d", n)
	}
}
