package resource

import "context"

const (
	ErrorResourceAlreadyExists = "Resource already exists"
	ErrorResourceNotFound      = "Resource not found"
)

// Access specifies the CRUD operations for a resource using generics.
// It supports context.Context for cancellation and timeouts.
type Access[K, V any] interface {
	Create(ctx context.Context, key K, value V) error
	Read(ctx context.Context, key K) (*V, error)
	ReadAll(ctx context.Context) ([]V, error)
	Update(ctx context.Context, key K, value V) error
	Delete(ctx context.Context, key K) error
}
