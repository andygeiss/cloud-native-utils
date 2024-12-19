package stability

import (
	"context"
	"time"

	"github.com/andygeiss/cloud-native-utils/service"
)

// Timeout wraps a service.Function and enforces a timeout on its execution.
// If the function does not complete within the specified duration, the context is canceled,
// and the function returns an error indicating a timeout.
func Timeout[IN any](fn service.Function[IN], duration time.Duration) service.Function[IN] {
	return func(ctx context.Context, in IN) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Create a child context with a timeout.
		withTimeout, cancel := context.WithTimeout(ctx, duration)
		defer cancel() // Clean up the timeout context.

		// Create a channel to signal the function's completion.
		done := make(chan error, 1)

		// Run the function in a separate goroutine.
		go func() {
			done <- fn(withTimeout, in)
		}()

		// Wait for either the function's completion or the timeout.
		select {
		case err := <-done:
			return err
		case <-withTimeout.Done():
			return withTimeout.Err()
		}
	}
}

// Timeout2 wraps a service.Function2 and enforces a timeout on its execution.
// If the function does not complete within the specified duration, the context is canceled,
// and the function returns an error indicating a timeout.
func Timeout2[IN, OUT any](fn service.Function2[IN, OUT], duration time.Duration) service.Function2[IN, OUT] {
	type result struct {
		out OUT
		err error
	}

	return func(ctx context.Context, in IN) (out OUT, err error) {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}

		// Create a child context with a timeout.
		withTimeout, cancel := context.WithTimeout(ctx, duration)
		defer cancel() // Clean up the timeout context.

		// Create a channel to capture the result of the function execution.
		resCh := make(chan result, 1)

		// Run the function in a separate goroutine.
		go func() {
			out, err := fn(withTimeout, in)
			resCh <- result{out, err}
		}()

		// Wait for either the function's result or the timeout.
		select {
		case res := <-resCh:
			return res.out, res.err
		case <-withTimeout.Done():
			return out, withTimeout.Err()
		}
	}
}
