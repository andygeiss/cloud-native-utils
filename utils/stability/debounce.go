package stability

import (
	"cloud-native/utils/service"
	"context"
	"sync"
	"time"
)

// Debounce wraps a given function (`fn`) to ensure it is not executed more often
// than the specified `duration`. If a new call occurs within the debounce duration,
// the previous call is canceled, and only the latest call proceeds.
func Debounce[IN, OUT any](fn service.Function[IN, OUT], duration time.Duration) service.Function[IN, OUT] {
	var debounceAt time.Time
	var err error
	var lastCancel context.CancelFunc
	var lastCtx context.Context
	var out OUT
	var mutex sync.Mutex
	return func(ctx context.Context, in IN) (OUT, error) {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		mutex.Lock()
		defer mutex.Unlock()
		// If the current time is within the debounce duration (`debounceAt`),
		// cancel the previous execution and return its result.
		if time.Now().Before(debounceAt) {
			if lastCancel != nil {
				lastCancel()
			}
			return out, err
		}
		// Create a new cancellable context for this execution.
		lastCtx, lastCancel = context.WithCancel(ctx)
		debounceAt = time.Now().Add(duration)
		// Execute the provided function `fn` with the new context and store its result.
		out, err = fn(lastCtx, in)
		return out, err
	}
}
