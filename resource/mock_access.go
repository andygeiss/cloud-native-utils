package resource

import "context"

// mockAccess is a mock implementation of Access[K, V].
// Each method has a corresponding mock function pointer (...Ptr).
// Use the builder-pattern to set the mock function pointers for each method.
// This allows for more flexible and readable test cases.
type mockAccess[K, V any] struct {
	createFn  func(ctx context.Context, key K, value V) error
	readFn    func(ctx context.Context, key K) (*V, error)
	readAllFn func(ctx context.Context) ([]V, error)
	updateFn  func(ctx context.Context, key K, value V) error
	deleteFn  func(ctx context.Context, key K) error
}

// NewMockAccess creates a new instance of MockAccess[K, V].
func NewMockAccess[K, V any]() *mockAccess[K, V] {
	return &mockAccess[K, V]{}
}

// Create creates a new resource with the given key and value.
func (a *mockAccess[K, V]) Create(ctx context.Context, key K, value V) error {
	return a.createFn(ctx, key, value)
}

// Read reads a resource with the given key.
func (a *mockAccess[K, V]) Read(ctx context.Context, key K) (value *V, err error) {
	return a.readFn(ctx, key)
}

// ReadAll reads all resources.
func (a *mockAccess[K, V]) ReadAll(ctx context.Context) (values []V, err error) {
	return a.readAllFn(ctx)
}

// Update updates a resource with the given key and value.
func (a *mockAccess[K, V]) Update(ctx context.Context, key K, value V) error {
	return a.updateFn(ctx, key, value)
}

// Delete deletes a resource with the given key.
func (a *mockAccess[K, V]) Delete(ctx context.Context, key K) (err error) {
	return a.deleteFn(ctx, key)
}

// WithCreateFn sets the mock function pointer for creating a resource.
func (a *mockAccess[K, V]) WithCreateFn(fn func(ctx context.Context, key K, value V) error) *mockAccess[K, V] {
	a.createFn = fn
	return a
}

// WithReadFn sets the mock function pointer for reading a resource.
func (a *mockAccess[K, V]) WithReadFn(fn func(ctx context.Context, key K) (*V, error)) *mockAccess[K, V] {
	a.readFn = fn
	return a
}

// WithReadAllFn sets the mock function pointer for reading all resources.
func (a *mockAccess[K, V]) WithReadAllFn(fn func(ctx context.Context) ([]V, error)) *mockAccess[K, V] {
	a.readAllFn = fn
	return a
}

// WithUpdateFn sets the mock function pointer for updating a resource.
func (a *mockAccess[K, V]) WithUpdateFn(fn func(ctx context.Context, key K, value V) error) *mockAccess[K, V] {
	a.updateFn = fn
	return a
}

// WithDeleteFn sets the mock function pointer for deleting a resource.
func (a *mockAccess[K, V]) WithDeleteFn(fn func(ctx context.Context, key K) error) *mockAccess[K, V] {
	a.deleteFn = fn
	return a
}
