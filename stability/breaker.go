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
func Breaker[T any](fn service.Function[T], threshold int) service.Function[T] {
	var failureCount int
	var lastCall = time.Now()
	var mutex sync.RWMutex
	return func(ctx context.Context) (*T, error) {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		// Acquire a read lock to check the breaker state.
		mutex.RLock()
		if diff := failureCount - threshold; diff >= 0 {
			// Calculate the next allowable retry time using exponential backoff.
			retryAt := lastCall.Add((2 << 1) * time.Second)
			if !time.Now().After(retryAt) {
				mutex.RUnlock()
				return nil, ErrorBreakerServiceUnavailable
			}
		}
		mutex.RUnlock()
		// Call the underlying service function.
		res, err := fn(ctx)
		// Acquire a write lock to update shared state.
		mutex.Lock()
		defer mutex.Unlock()
		lastCall = time.Now()
		if err != nil {
			failureCount++
			return nil, err
		}
		// Reset the failure count on a successful call.
		failureCount = 0
		return res, nil
	}
}
