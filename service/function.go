package service

import "context"

// Function gathers together things that change for the same reasons.
// A context must be handled to be "cloud native" because it allows
// propagation of deadlines, cancellation signals, and other request-scoped values
// across API boundaries and between processes. This ensures efficient resource
// utilization and proper handling of asynchronous workflows in distributed systems.
type Function[T any] func(ctx context.Context) (*T, error)
