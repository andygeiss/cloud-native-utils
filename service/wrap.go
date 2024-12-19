package service

import (
	"context"
)

// Wrap takes a simple function `fn` and converts it into a context-aware `Function`.
// It ensures the function respects context cancellation or timeout.
func Wrap[IN any](fn func(in IN) (err error)) Function[IN] {

	// Define a type to hold the function's error.
	type response struct {
		err error
	}

	// Return a context-aware function.
	return func(ctx context.Context, in IN) (err error) {
		if ctx.Err() != nil { // Check if the context is already canceled or expired.
			return ctx.Err()
		}
		ch := make(chan response, 1)

		// Run the function in a separate goroutine.
		go func() {
			err := fn(in)
			ch <- response{err}
		}()

		// Wait for the function to complete or the context to cancel.
		select {
		case res := <-ch:
			return res.err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Wrap2 converts a function `fn` into a context-aware `Function2`.
// It supports cancellation and timeout via the provided context.
func Wrap2[IN, OUT any](fn func(in IN) (out OUT, err error)) Function2[IN, OUT] {

	// Define a type to hold the function's result and error.
	type response struct {
		result OUT
		err    error
	}

	// Return a context-aware function.
	return func(ctx context.Context, in IN) (out OUT, err error) {
		if ctx.Err() != nil { // Check if the context is already canceled or expired.
			return out, ctx.Err()
		}
		ch := make(chan response, 1)

		// Run the function in a separate goroutine.
		go func() {
			res, err := fn(in)
			ch <- response{res, err}
		}()

		// Wait for the function to complete or the context to cancel.
		select {
		case res := <-ch:
			return res.result, res.err
		case <-ctx.Done():
			return out, ctx.Err()
		}
	}
}
