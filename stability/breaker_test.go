package stability_test

import (
	"cloud-native/stability"
	"context"
	"testing"
)

func TestBreaker_Call_Once_Succeeds(t *testing.T) {
	fn := stability.Breaker[int](mockAlwaysSucceeds(), 3)
	res, err := fn(context.Background())
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

func TestBreaker_Call_Once_Fails(t *testing.T) {
	fn := stability.Breaker[int](mockAlwaysFails(), 3)
	res, err := fn(context.Background())
	if res != nil {
		t.Fatal("result must be nil")
	}
	if err.Error() != "error" {
		t.Fatalf("error must be %s, but got %s", "error", err.Error())
	}
}

func TestBreaker_Call_Threshold(t *testing.T) {
	threshold := 3
	fn := stability.Breaker[int](mockAlwaysFails(), threshold)
	for range threshold {
		_, _ = fn(context.Background())
	}
	res, err := fn(context.Background())
	if res != nil {
		t.Fatal("result must be nil")
	}
	if err.Error() != stability.ErrorBreakerServiceUnavailable.Error() {
		t.Fatalf("error must be %s, but got %s", "error", err.Error())
	}
}

func TestBreaker_Call_Concurrent(t *testing.T) {
	const goroutines = 10
	threshold := 3
	fn := stability.Breaker[int](mockAlwaysSucceeds(), threshold)
	errs := make(chan error, 3)
	for range goroutines {
		go func() {
			_, err := fn(context.Background())
			errs <- err
		}()
	}
	for range goroutines {
		err := <-errs
		if err != nil {
			t.Fatalf("error must be nil, but got %v", err)
		}
	}
}
