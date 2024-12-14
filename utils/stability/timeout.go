package stability

import (
	"cloud-native/utils/service"
	"context"
	"time"
)

// Timeout wraps a service.Function and enforces a timeout on its execution.
// If the function does not complete within the specified duration, the context is canceled,
// and the function returns an error indicating a timeout.
func Timeout[IN, OUT any](fn service.Function[IN, OUT], duration time.Duration) service.Function[IN, OUT] {
	type result struct {
		out OUT
		err error
	}
	return func(ctx context.Context, in IN) (out OUT, err error) {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		// Create a child context with a timeout based on the specified duration.
		withTimeout, cancel := context.WithTimeout(ctx, duration)
		defer cancel() // Ensure the timeout context is properly cleaned up to avoid resource leaks.
		// Create a channel to capture the result of the function execution.
		resCh := make(chan result)
		// Run the wrapped function in a separate goroutine.
		// This allows us to listen for both the function's result and the timeout in parallel.
		go func() {
			defer close(resCh)
			out, err := fn(withTimeout, in)
			resCh <- result{out, err}
		}()
		// Use a select statement to wait for either the function's result or the timeout.
		select {
		case res := <-resCh:
			return res.out, res.err
		case <-withTimeout.Done():
			return out, withTimeout.Err()
		}
	}
}
