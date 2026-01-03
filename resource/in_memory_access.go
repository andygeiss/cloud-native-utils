package resource

import (
	"context"
	"errors"
	"sync"
)

// InMemoryAccess is a generic access implementation backed by a mock, in-memory and JSON file.
type InMemoryAccess[K comparable, V any] struct {
	kv    map[K]V
	mutex sync.RWMutex
}

// NewInMemoryAccess creates a new in-memory access.
func NewInMemoryAccess[K comparable, V any]() *InMemoryAccess[K, V] {
	return &InMemoryAccess[K, V]{
		kv: make(map[K]V),
	}
}

// Create creates a new resource.
func (a *InMemoryAccess[K, V]) Create(ctx context.Context, key K, value V) (err error) {
	// Skip if context is canceled or timed out.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Ensure that only one goroutine can write to the map at a time.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Check if resource already exists.
	if _, alreadyExists := a.kv[key]; alreadyExists {
		return errors.New(ErrorResourceAlreadyExists)
	}

	// Add resource if not exists.
	a.kv[key] = value
	return nil
}

// Delete deletes a resource.
func (a *InMemoryAccess[K, V]) Delete(ctx context.Context, key K) (err error) {
	// Skip if context is canceled or timed out.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Ensure that only one goroutine can write to the map at a time.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Check if resource exists.
	if _, exists := a.kv[key]; exists {
		delete(a.kv, key)
		return nil
	}

	return errors.New(ErrorResourceNotFound)
}

// Read reads a resource.
func (a *InMemoryAccess[K, V]) Read(ctx context.Context, key K) (ptr *V, err error) {
	// Skip if context is canceled or timed out.
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Ensure that read only access to the map is allowed.
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// Check if resource already exists.
	if val, exists := a.kv[key]; exists {
		return &val, nil
	}

	return nil, errors.New(ErrorResourceNotFound)
}

// ReadAll reads all resources.
func (a *InMemoryAccess[K, V]) ReadAll(ctx context.Context) (values []V, err error) {
	// Skip if context is canceled or timed out.
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Ensure that read only access to the map is allowed.
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// Create a slice to hold the values.
	for _, value := range a.kv {
		values = append(values, value)
	}
	return values, nil
}

// Update updates a resource.
func (a *InMemoryAccess[K, V]) Update(ctx context.Context, key K, value V) (err error) {
	// Skip if context is canceled or timed out.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Ensure that only one goroutine can write to the map at a time.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Check if resource exists.
	if _, exists := a.kv[key]; exists {
		a.kv[key] = value
		return nil
	}

	return errors.New(ErrorResourceNotFound)
}
