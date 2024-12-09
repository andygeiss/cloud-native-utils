package stability

import (
	"cloud-native/service"
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrorBreakerServiceUnavailable = errors.New("service unavailable")
)

// Breaker wraps a service function with a circuit breaker mechanism.
// It tracks failures and prevents further calls when a failure threshold is reached.
func Breaker[IN, OUT any](fn service.Function[IN, OUT], threshold int) service.Function[IN, OUT] {
	var failureCount int
	var lastCall = time.Now()
	var mutex sync.RWMutex
	return func(ctx context.Context, in IN) (out OUT, err error) {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		// Acquire a read lock to check the breaker state.
		mutex.RLock()
		if diff := failureCount - threshold; diff >= 0 {
			// Calculate the next allowable retry time using exponential backoff.
			retryAt := lastCall.Add((2 << 1) * time.Second)
			if !time.Now().After(retryAt) {
				mutex.RUnlock()
				return out, ErrorBreakerServiceUnavailable
			}
		}
		mutex.RUnlock()
		// Call the underlying service function.
		res, err := fn(ctx, in)
		// Acquire a write lock to update shared state.
		mutex.Lock()
		defer mutex.Unlock()
		lastCall = time.Now()
		if err != nil {
			failureCount++
			return out, err
		}
		// Reset the failure count on a successful call.
		failureCount = 0
		return res, nil
	}
}
