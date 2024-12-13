package stability_test

import (
	"cloud-native/utils/stability"
	"context"
	"testing"
	"time"
)

func TestThrottle_Call_Once_Succeeds(t *testing.T) {
	fn := stability.Throttle[int](mockAlwaysSucceeds(), 3, 3, 100*time.Millisecond)
	res, err := fn(context.Background(), 42)
	if err != nil {
		t.Fatalf("error must be nil, but got %v", err)
	}
	if res != 42 {
		t.Fatalf("result must be %d, but got %d", 42, res)
	}
}

func TestThrottle_Call_Twice_Returns_Error(t *testing.T) {
	fn := stability.Throttle[int](mockSucceedsTimes(1), 3, 3, 100*time.Millisecond)
	fn(context.Background(), 42)
	_, err := fn(context.Background(), 42)
	if err == nil {
		t.Fatal("error must be not nil")
	}
	if err.Error() != "error" {
		t.Fatalf("error must be correct, but got %v", err)
	}
}

func TestThrottle_Call_Reaches_Max_Tokens(t *testing.T) {
	fn := stability.Throttle[int](mockSucceedsTimes(10), 3, 1, 100*time.Millisecond)
	fn(context.Background(), 42)
	fn(context.Background(), 42)
	fn(context.Background(), 42)
	_, err := fn(context.Background(), 42)
	if err == nil {
		t.Fatal("error must be not nil")
	}
	if err.Error() != "too many calls" {
		t.Fatalf("error must be correct, but got %v", err)
	}
}

func TestThrottle_Call_Refills_After_Duration(t *testing.T) {
	fn := stability.Throttle[int](mockSucceedsTimes(10), 3, 1, 100*time.Millisecond)
	fn(context.Background(), 42)
	fn(context.Background(), 42)
	fn(context.Background(), 42)
	time.Sleep(150 * time.Millisecond)
	res, err := fn(context.Background(), 42)
	if err != nil {
		t.Fatalf("error must be nil, but got %v", err)
	}
	if res != 42 {
		t.Fatalf("result must be %d, but got %d", 42, res)
	}
}
