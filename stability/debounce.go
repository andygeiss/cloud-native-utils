package stability

import (
	"context"
	"sync"
	"time"

	"github.com/andygeiss/cloud-native-utils/service"
)

// Debounce ensures that fn(...) is ultimately called only once if multiple
// calls arrive within the given duration. Each call waits until that one
// call finishes, returning the same result/error.
func Debounce[IN, OUT any](fn service.Function[IN, OUT], duration time.Duration) service.Function[IN, OUT] {
	var (
		mu sync.Mutex

		// Timer that triggers the actual fn call
		timer *time.Timer

		// We'll store the last call's context and input
		lastCtx context.Context
		lastIn  IN

		// Once fn has run, we store its output and error
		result OUT
		err    error

		// All waiting callers block on this list of channels.
		// When the timer fires and fn completes, we close them.
		waiters []chan struct{}
	)

	return func(ctx context.Context, in IN) (OUT, error) {
		mu.Lock()

		// Update our "latest" context & input
		lastCtx = ctx //nolint:fatcontext // intentional: debounce uses the latest context
		lastIn = in

		// Each caller waits on its own channel to learn when the
		// "debounced" call is done.
		done := make(chan struct{})
		waiters = append(waiters, done)

		// If there was already a timer waiting to call fn,
		// stop & reset it.
		if timer != nil {
			timer.Stop()
		}

		// Start (or reset) the timer to call fn after 'duration'
		timer = time.AfterFunc(duration, func() {
			// When the timer fires, we do the actual call:
			mu.Lock()

			// Actually call the underlying fn with the *last* inputs
			result, err = fn(lastCtx, lastIn)

			// Wake all callers waiting on this batch.
			for _, w := range waiters {
				close(w)
			}
			waiters = nil

			mu.Unlock()
		})

		mu.Unlock()

		// Now, each caller waits until the "done" channel is closed,
		// meaning the single fn call has completed.
		select {
		case <-done:
			// The single call is finished. We'll return the same result to everyone.
		case <-ctx.Done():
			// If the caller's own context is canceled first, return immediately.
			// You could also do further logic to remove this caller from waiters,
			// but that’s more complex.
			return *new(OUT), ctx.Err()
		}

		// Once here, the result/err are set (or we bailed early).
		mu.Lock()
		defer mu.Unlock()
		return result, err
	}
}
