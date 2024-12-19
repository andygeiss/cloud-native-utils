package service

import "context"

// Function2 defines a generic function type with two type parameters, IN and OUT.
// This function type takes a context (ctx) and an input of type IN, then returns
// an output of type OUT along with an error.
// The use of context ensures "cloud native" capabilities by supporting
// request-scoped values, deadlines, and cancellation signals, which are critical
// for resource-efficient and robust distributed systems.
type Function2[IN, OUT any] func(ctx context.Context, in IN) (out OUT, err error)

// Function defines a generic function type with a single type parameter, IN.
// This function type takes a context (ctx) and an input of type IN, then returns
// an error. It is useful for operations where only an error status needs to be returned,
// and the use of context facilitates proper resource management in asynchronous workflows.
type Function[IN any] func(ctx context.Context, in IN) (err error)
