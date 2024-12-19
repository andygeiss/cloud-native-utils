package stability

import (
	"context"
	"sync"
	"time"

	"github.com/andygeiss/cloud-native-utils/service"
)

// Debounce delays the invocation of the provided function until a specified duration has passed since its last invocation.
func Debounce[IN any](fn service.Function[IN], duration time.Duration) service.Function[IN] {
	var mutex sync.RWMutex
	var timer *time.Timer
	var lastErr error

	return func(ctx context.Context, in IN) (err error) {

		// Early return if a timer is active
		mutex.RLock()
		if timer != nil {
			defer mutex.RUnlock()
			return lastErr
		}
		mutex.RUnlock()

		mutex.Lock()
		defer mutex.Unlock()

		// Reset the timer after the duration
		timer = time.AfterFunc(duration, func() {
			timer = nil
		})

		// Execute the function and store the result
		lastErr = fn(ctx, in)

		return lastErr
	}
}

// Debounce2 delays the invocation of a function with input and output until a specified duration has passed since its last invocation.
func Debounce2[IN, OUT any](fn service.Function2[IN, OUT], duration time.Duration) service.Function2[IN, OUT] {
	var lastErr error
	var lastOut OUT
	var timer *time.Timer
	var mutex sync.RWMutex

	return func(ctx context.Context, in IN) (out OUT, err error) {
		// Early return if a timer is active
		mutex.RLock()
		if timer != nil {
			defer mutex.RUnlock()
			return lastOut, lastErr
		}
		mutex.RUnlock()

		mutex.Lock()
		defer mutex.Unlock()

		// Reset the timer after the duration
		timer = time.AfterFunc(duration, func() {
			timer = nil
		})

		// Execute the function and store the results
		lastOut, lastErr = fn(ctx, in)

		return lastOut, lastErr
	}
}
