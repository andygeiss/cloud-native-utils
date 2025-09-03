package resource

import (
	"errors"
	"sync"
)

type inMemoryAccess[K comparable, V any] struct {
	kv    map[K]V
	mutex sync.RWMutex
}

// NewInMemoryAccess creates a new in-memory access.
func NewInMemoryAccess[K comparable, V any]() *inMemoryAccess[K, V] {
	return &inMemoryAccess[K, V]{
		kv: make(map[K]V),
	}
}

func (a *inMemoryAccess[K, V]) Create(key K, value V) error {
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

func (a *inMemoryAccess[K, V]) Read(key K) (*V, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// Check if resource already exists.
	if val, exists := a.kv[key]; exists {
		return &val, nil
	}

	return nil, errors.New(ErrorResourceNotFound)
}

func (a *inMemoryAccess[K, V]) ReadAll() ([]V, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	var values []V
	for _, value := range a.kv {
		values = append(values, value)
	}

	return values, nil
}

func (a *inMemoryAccess[K, V]) Update(key K, value V) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Check if resource exists.
	if _, exists := a.kv[key]; exists {
		a.kv[key] = value
		return nil
	}

	return errors.New(ErrorResourceNotFound)
}

func (a *inMemoryAccess[K, V]) Delete(key K) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Check if resource exists.
	if _, exists := a.kv[key]; exists {
		delete(a.kv, key)
		return nil
	}

	return errors.New(ErrorResourceNotFound)
}
