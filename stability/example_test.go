package stability_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andygeiss/cloud-native-utils/service"
	"github.com/andygeiss/cloud-native-utils/stability"
)

func ExampleRetry() {
	attempts := 0
	flaky := service.Wrap(func(in string) (string, error) {
		attempts++
		if attempts < 3 {
			return "", errors.New("temporary failure")
		}
		return "ok: " + in, nil
	})

	out, err := stability.Retry(flaky, 5, time.Millisecond)(context.Background(), "ping")
	fmt.Println(out, err, attempts)
	// Output: ok: ping <nil> 3
}

func ExampleBreaker() {
	failing := service.Wrap(func(in string) (string, error) {
		return "", errors.New("backend down")
	})

	// After two failures the breaker opens and stops calling the backend.
	guarded := stability.Breaker(failing, 2)
	ctx := context.Background()
	for range 3 {
		_, err := guarded(ctx, "ping")
		fmt.Println(err)
	}
	// Output:
	// backend down
	// backend down
	// service unavailable
}

func ExampleTimeout() {
	slow := service.Wrap(func(in string) (string, error) {
		time.Sleep(50 * time.Millisecond)
		return in, nil
	})

	_, err := stability.Timeout(slow, time.Millisecond)(context.Background(), "ping")
	fmt.Println(errors.Is(err, context.DeadlineExceeded))
	// Output: true
}

// The wrappers compose: each one takes a Function and returns a Function.
func Example_composition() {
	calls := 0
	backend := service.Wrap(func(in string) (string, error) {
		calls++
		return "pong", nil
	})

	fn := stability.Retry(backend, 3, time.Millisecond)
	fn = stability.Timeout(fn, time.Second)

	out, err := fn(context.Background(), "ping")
	fmt.Println(out, err, calls)
	// Output: pong <nil> 1
}
