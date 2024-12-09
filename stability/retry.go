package stability

import (
	"cloud-native/service"
	"context"
	"time"
)

// Retry wraps a given function (`fn`) to retry its execution upon failure.
// The function will be retried up to `maxRetries` times with a delay of `delay` between retries.
// If the context is canceled during retries, it stops immediately and returns the context error.
func Retry[IN, OUT any](fn service.Function[IN, OUT], maxRetries int, delay time.Duration) service.Function[IN, OUT] {
	return func(ctx context.Context, in IN) (out OUT, err error) {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		for retries := 0; ; retries++ {
			// Call the provided function and capture its result and error.
			res, err := fn(ctx, in)
			// If the function succeeds (err == nil), or the maximum number of retries has been reached,
			// return the result and error (if any).
			if err == nil || retries >= maxRetries {
				return res, err
			}
			select {
			// Wait for the delay duration before retrying.
			case <-time.After(delay):
			// If the context is canceled during the wait, stop retrying.
			case <-ctx.Done():
				return out, ctx.Err()
			}
		}
	}
}
