package service

import (
	"context"
)

// Wrap converts a simple function into one that respects a provided context.
// This enables the wrapped function to respond to context cancellation or timeout.
func Wrap[T any](fn func() (*T, error)) Function[T] {
	// Define a type to hold the function's result and error.
	type response struct {
		result *T
		err    error
	}
	// Return a function that incorporates timeout logic.
	return func(ctx context.Context) (*T, error) {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		ch := make(chan response, 1)
		// Execute the function in a separate goroutine.
		go func() {
			res, err := fn()
			ch <- response{res, err}
		}()
		// Wait for either the function to complete or the context to be canceled.
		select {
		case res := <-ch:
			return res.result, res.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
