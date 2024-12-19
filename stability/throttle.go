package stability

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/andygeiss/cloud-native-utils/service"
)

var (
	// ErrorThrottleTooManyCalls is returned when the rate limit is exceeded.
	ErrorThrottleTooManyCalls = errors.New("too many calls")
)

// Throttle adds rate-limiting behavior to the provided function (`fn`).
// The function can only be called up to `maxTokens` times initially,
// and then tokens are refilled by `refill` every `duration`. If the limit is exceeded,
// the function returns `ErrorThrottleTooManyCalls`.
func Throttle[IN any](fn service.Function[IN], maxTokens, refill uint, duration time.Duration) service.Function[IN] {
	var tokens = maxTokens
	var once sync.Once
	var mutex sync.Mutex

	return func(ctx context.Context, in IN) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Start the refill process only once.
		once.Do(func() {
			ticker := time.NewTicker(duration)
			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						mutex.Lock()
						if tokens < maxTokens {
							tokens += refill
							if tokens > maxTokens {
								tokens = maxTokens
							}
						}
						mutex.Unlock()
					}
				}
			}()
		})

		mutex.Lock()
		defer mutex.Unlock()

		// Check if a token is available.
		if tokens == 0 {
			return ErrorThrottleTooManyCalls
		}
		tokens--

		// Execute the function.
		return fn(ctx, in)
	}
}

// Throttle2 adds rate-limiting behavior to the provided function (`fn`).
// The function can only be called up to `maxTokens` times initially,
// and then tokens are refilled by `refill` every `duration`. If the limit is exceeded,
// the function returns `ErrorThrottleTooManyCalls`.
func Throttle2[IN, OUT any](fn service.Function2[IN, OUT], maxTokens, refill uint, duration time.Duration) service.Function2[IN, OUT] {
	var tokens = maxTokens
	var once sync.Once
	var mutex sync.Mutex

	return func(ctx context.Context, in IN) (out OUT, err error) {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}

		// Start the refill process only once.
		once.Do(func() {
			ticker := time.NewTicker(duration)
			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						mutex.Lock()
						if tokens < maxTokens {
							tokens += refill
							if tokens > maxTokens {
								tokens = maxTokens
							}
						}
						mutex.Unlock()
					}
				}
			}()
		})

		mutex.Lock()
		defer mutex.Unlock()

		// Check if a token is available.
		if tokens == 0 {
			return out, ErrorThrottleTooManyCalls
		}
		tokens--

		// Execute the function and return its result.
		return fn(ctx, in)
	}
}
