package resource

import (
	"errors"

	"github.com/andygeiss/cloud-native-utils/efficiency"
)

const (
	ErrorResourceAlreadyExists = "resource already exists"
	ErrorResourceNotFound      = "resource not found"
)

type inMemoryAccess[K comparable, V any] struct {
	sharding efficiency.Sharding[K, V]
}

func NewInMemoryAccess[K comparable, V any](shards int) *inMemoryAccess[K, V] {
	return &inMemoryAccess[K, V]{
		sharding: efficiency.NewSharding[K, V](shards),
	}
}

func (a *inMemoryAccess[K, V]) Create(key K, value V) error {
	// Check if resource already exists.
	if _, alreadyExists := a.sharding.Get(key); alreadyExists {
		return errors.New(ErrorResourceAlreadyExists)
	}

	// Add resource if not exists.
	a.sharding.Put(key, value)
	return nil
}

func (a *inMemoryAccess[K, V]) Read(key K) (*V, error) {
	// Check if resource already exists.
	if val, exists := a.sharding.Get(key); exists {
		return &val, nil
	}

	return nil, errors.New(ErrorResourceNotFound)
}

func (a *inMemoryAccess[K, V]) Update(key K, value V) error {
	// Check if resource exists.
	if _, exists := a.sharding.Get(key); exists {
		a.sharding.Put(key, value)
		return nil
	}

	return errors.New(ErrorResourceNotFound)
}

func (a *inMemoryAccess[K, V]) Delete(key K) error {
	// Check if resource exists.
	if _, exists := a.sharding.Get(key); exists {
		a.sharding.Delete(key)
		return nil
	}

	return errors.New(ErrorResourceNotFound)
}
