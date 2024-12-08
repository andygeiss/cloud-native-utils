package stability

import (
	"cloud-native/service"
	"context"
	"time"
)

// Retry wraps a given function (`fn`) to retry its execution upon failure.
// The function will be retried up to `maxRetries` times with a delay of `delay` between retries.
// If the context is canceled during retries, it stops immediately and returns the context error.
func Retry[T any](fn service.Function[T], maxRetries int, delay time.Duration) service.Function[T] {
	return func(ctx context.Context) (*T, error) {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		for retries := 0; ; retries++ {
			// Call the provided function and capture its result and error.
			res, err := fn(ctx)
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
				return nil, ctx.Err()
			}
		}
	}
}
