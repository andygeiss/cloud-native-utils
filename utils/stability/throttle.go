package stability

import (
	"cloud-native/utils/service"
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrorThrottleTooManyCalls = errors.New("too many calls")
)

// Throttle adds rate-limiting behavior to the provided function (`fn`).
// The function can only be called up to `maxTokens` times initially,
// and then tokens are refilled by `refill` every `duration`. If the limit is exceeded,
// the function returns `ErrorThrottleTooManyCalls`.
func Throttle[IN, OUT any](fn service.Function[IN, OUT], maxTokens, refill uint, duration time.Duration) service.Function[IN, OUT] {
	var tokens = maxTokens
	var once sync.Once
	var mutex sync.Mutex
	return func(ctx context.Context, in IN) (out OUT, err error) {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		// Use `once` to ensure the refill logic runs exactly once, even with multiple callers.
		once.Do(func() {
			// Create a ticker to trigger token refills at the specified `duration`.
			ticker := time.NewTimer(duration)
			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						mutex.Lock()
						count := tokens + refill
						if count > maxTokens {
							count = maxTokens
						}
						tokens = count
						mutex.Unlock()
					}
				}
			}()
		})
		mutex.Lock()
		defer mutex.Unlock()
		if tokens <= 0 {
			return out, ErrorThrottleTooManyCalls
		}
		tokens--
		// Call the wrapped function and return its result.
		return fn(ctx, in)
	}
}
