package resource

import "context"

// MockAccess is a mock implementation of Access[K, V].
// Each method has a corresponding mock function pointer (...Ptr).
// Use the builder-pattern to set the mock function pointers for each method.
// This allows for more flexible and readable test cases.
type MockAccess[K, V any] struct {
	createFn  func(ctx context.Context, key K, value V) error
	deleteFn  func(ctx context.Context, key K) error
	readAllFn func(ctx context.Context) ([]V, error)
	readFn    func(ctx context.Context, key K) (*V, error)
	updateFn  func(ctx context.Context, key K, value V) error
}

// NewMockAccess creates a new instance of MockAccess[K, V].
func NewMockAccess[K, V any]() *MockAccess[K, V] {
	return &MockAccess[K, V]{}
}

// Create creates a new resource with the given key and value.
func (a *MockAccess[K, V]) Create(ctx context.Context, key K, value V) error {
	return a.createFn(ctx, key, value)
}

// Delete deletes a resource with the given key.
func (a *MockAccess[K, V]) Delete(ctx context.Context, key K) error {
	return a.deleteFn(ctx, key)
}

// Read reads a resource with the given key.
func (a *MockAccess[K, V]) Read(ctx context.Context, key K) (*V, error) {
	return a.readFn(ctx, key)
}

// ReadAll reads all resources.
func (a *MockAccess[K, V]) ReadAll(ctx context.Context) ([]V, error) {
	return a.readAllFn(ctx)
}

// Update updates a resource with the given key and value.
func (a *MockAccess[K, V]) Update(ctx context.Context, key K, value V) error {
	return a.updateFn(ctx, key, value)
}

// WithCreateFn sets the mock function pointer for creating a resource.
func (a *MockAccess[K, V]) WithCreateFn(fn func(ctx context.Context, key K, value V) error) *MockAccess[K, V] {
	a.createFn = fn
	return a
}

// WithDeleteFn sets the mock function pointer for deleting a resource.
func (a *MockAccess[K, V]) WithDeleteFn(fn func(ctx context.Context, key K) error) *MockAccess[K, V] {
	a.deleteFn = fn
	return a
}

// WithReadAllFn sets the mock function pointer for reading all resources.
func (a *MockAccess[K, V]) WithReadAllFn(fn func(ctx context.Context) ([]V, error)) *MockAccess[K, V] {
	a.readAllFn = fn
	return a
}

// WithReadFn sets the mock function pointer for reading a resource.
func (a *MockAccess[K, V]) WithReadFn(fn func(ctx context.Context, key K) (*V, error)) *MockAccess[K, V] {
	a.readFn = fn
	return a
}

// WithUpdateFn sets the mock function pointer for updating a resource.
func (a *MockAccess[K, V]) WithUpdateFn(fn func(ctx context.Context, key K, value V) error) *MockAccess[K, V] {
	a.updateFn = fn
	return a
}
