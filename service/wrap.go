package service

import (
	"context"
)

// Wrap converts a simple function into one that respects a provided context.
// This enables the wrapped function to respond to context cancellation or timeout.
func Wrap[IN, OUT any](fn func(in IN) (OUT, error)) Function[IN, OUT] {
	// Define a type to hold the function's result and error.
	type response struct {
		result OUT
		err    error
	}

	// Return a function that incorporates timeout logic.
	return func(ctx context.Context, in IN) (OUT, error) {
		var zero OUT
		if ctx.Err() != nil {
			return zero, ctx.Err()
		}

		// Execute the function in a separate goroutine.
		ch := make(chan response, 1)
		go func() {
			res, err := fn(in)
			ch <- response{res, err}
		}()

		// Wait for either the function to complete or the context to be canceled.
		select {
		case res := <-ch:
			return res.result, res.err
		case <-ctx.Done():
			return zero, ctx.Err()
		}
	}
}
