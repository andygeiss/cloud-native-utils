package stability

import (
	"context"
	"sync"
	"time"

	"github.com/andygeiss/cloud-native-utils/service"
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
		// If the caller's context is already canceled, bail out early
		if ctx.Err() != nil {
			return out, ctx.Err()
		}

		mutex.Lock()
		defer mutex.Unlock()

		// If we are still "inside" the last debounce window, cancel the old call
		if time.Now().Before(debounceAt) {
			if lastCancel != nil {
				lastCancel()
			}
		}

		// Create a new cancellable context for the current (latest) call
		lastCtx, lastCancel = context.WithCancel(ctx)

		// Extend the "no new calls" window from *now*
		debounceAt = time.Now().Add(duration)

		// Now actually call fn for the latest input
		out, err = fn(lastCtx, in)

		return out, err
	}
}
