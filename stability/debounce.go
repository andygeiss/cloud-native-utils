package stability

import (
	"cloud-native/service"
	"context"
	"sync"
	"time"
)

// Debounce wraps a given function (`fn`) to ensure it is not executed more often
// than the specified `duration`. If a new call occurs within the debounce duration,
// the previous call is canceled, and only the latest call proceeds.
func Debounce[T any](fn service.Function[T], duration time.Duration) service.Function[T] {
	var debounceAt time.Time
	var err error
	var lastCancel context.CancelFunc
	var lastCtx context.Context
	var res *T
	var mutex sync.Mutex
	return func(ctx context.Context) (*T, error) {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		mutex.Lock()
		defer mutex.Unlock()
		// If the current time is within the debounce duration (`debounceAt`),
		// cancel the previous execution and return its result.
		if time.Now().Before(debounceAt) {
			if lastCancel != nil {
				lastCancel()
			}
			return res, err
		}
		// Create a new cancellable context for this execution.
		lastCtx, lastCancel = context.WithCancel(ctx)
		debounceAt = time.Now().Add(duration)
		// Execute the provided function `fn` with the new context and store its result.
		res, err = fn(lastCtx)
		return res, err
	}
}
