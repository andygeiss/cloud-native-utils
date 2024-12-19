package stability

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/andygeiss/cloud-native-utils/service"
)

var (
	// ErrorBreakerServiceUnavailable indicates the circuit breaker is open.
	ErrorBreakerServiceUnavailable = errors.New("service unavailable")
)

// Breaker wraps a service function with a circuit breaker mechanism.
// Prevents calls to `fn` when failures exceed the threshold.
func Breaker[IN any](fn service.Function[IN], threshold int) service.Function[IN] {
	var failureCount int
	var lastCall = time.Now()
	var mutex sync.RWMutex

	return func(ctx context.Context, in IN) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Check circuit breaker state.
		mutex.RLock()
		if failureCount >= threshold {
			num := failureCount - threshold
			retryAt := lastCall.Add(time.Second * time.Duration(2<<num))
			if !time.Now().After(retryAt) {
				mutex.RUnlock()
				return ErrorBreakerServiceUnavailable
			}
		}
		mutex.RUnlock()

		// Call the wrapped function.
		err = fn(ctx, in)

		// Update state based on the result.
		mutex.Lock()
		defer mutex.Unlock()
		lastCall = time.Now()
		if err != nil {
			failureCount++
			return err
		}

		// Reset failure count on success.
		failureCount = 0
		return nil
	}
}

// Breaker2 wraps a Function2 with a circuit breaker mechanism.
// Prevents calls to `fn` when failures exceed the threshold.
func Breaker2[IN, OUT any](fn service.Function2[IN, OUT], threshold int) service.Function2[IN, OUT] {
	var failureCount int
	var lastCall = time.Now()
	var mutex sync.RWMutex

	return func(ctx context.Context, in IN) (out OUT, err error) {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}

		// Check circuit breaker state.
		mutex.RLock()
		if failureCount >= threshold {
			num := failureCount - threshold
			retryAt := lastCall.Add(time.Second * time.Duration(2<<num))
			if !time.Now().After(retryAt) {
				mutex.RUnlock()
				return out, ErrorBreakerServiceUnavailable
			}
		}
		mutex.RUnlock()

		// Call the wrapped function.
		res, err := fn(ctx, in)

		// Update state based on the result.
		mutex.Lock()
		defer mutex.Unlock()
		lastCall = time.Now()
		if err != nil {
			failureCount++
			return out, err
		}

		// Reset failure count on success.
		failureCount = 0
		return res, nil
	}
}
