package stability_test

import (
	"cloud-native/utils/stability"
	"context"
	"testing"
	"time"
)

func TestRetry_Succeeds(t *testing.T) {
	fn := stability.Retry[int](mockAlwaysSucceeds(), 3, 10*time.Millisecond)
	res, err := fn(context.Background(), 42)
	if err != nil {
		t.Fatalf("error must be nil, but got %v", err)
	}
	if res != 42 {
		t.Fatalf("result must be %d, but got %d", 42, res)
	}
}

func TestRetry_Suceeds_With_Retries(t *testing.T) {
	fn := stability.Retry[int](mockFailsTimes(2), 3, 10*time.Millisecond)
	res, err := fn(context.Background(), 42)
	if err != nil {
		t.Fatalf("error must be nil, but got %v", err)
	}
	if res != 42 {
		t.Fatalf("result must be %d, but got %d", 42, res)
	}
}
