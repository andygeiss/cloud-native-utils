package resource

import (
	"context"
	"sync"
)

// IndexFunc extracts an index key from a value.
// Returns an empty string if the value should not be indexed.
type IndexFunc[V any] func(V) string

// IndexedAccess wraps a resource.Access and maintains secondary indexes.
// It supports both unique and non-unique indexes (stored as lists of keys).
type IndexedAccess[K comparable, V any] struct {
	access     Access[K, V]
	indexes    map[string]map[string][]K
	indexFuncs map[string]IndexFunc[V]
	mu         sync.RWMutex
}

// NewIndexedAccess creates a new indexed access wrapper.
func NewIndexedAccess[K comparable, V any](access Access[K, V]) *IndexedAccess[K, V] {
	return &IndexedAccess[K, V]{
		access:     access,
		indexes:    make(map[string]map[string][]K),
		indexFuncs: make(map[string]IndexFunc[V]),
	}
}

// AddIndex adds a new secondary index.
// name is the unique name of the index.
// fn is the function to extract the index key from the value.
func (a *IndexedAccess[K, V]) AddIndex(name string, fn IndexFunc[V]) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.indexFuncs[name] = fn
	a.indexes[name] = make(map[string][]K)
}

// Create stores a new value and updates secondary indexes.
func (a *IndexedAccess[K, V]) Create(ctx context.Context, key K, value V) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := a.access.Create(ctx, key, value); err != nil {
		return err
	}

	// Update secondary indexes
	for name, fn := range a.indexFuncs {
		idxKey := fn(value)
		if idxKey != "" {
			a.indexes[name][idxKey] = append(a.indexes[name][idxKey], key)
		}
	}

	return nil
}

// Read retrieves a value by its primary key.
func (a *IndexedAccess[K, V]) Read(ctx context.Context, key K) (*V, error) {
	return a.access.Read(ctx, key)
}

// ReadAll retrieves all values.
func (a *IndexedAccess[K, V]) ReadAll(ctx context.Context) ([]V, error) {
	return a.access.ReadAll(ctx)
}

// Update updates an existing value and its secondary indexes.
func (a *IndexedAccess[K, V]) Update(ctx context.Context, key K, value V) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Remove old index entries
	oldValue, err := a.access.Read(ctx, key)
	if err == nil && oldValue != nil {
		for name, fn := range a.indexFuncs {
			idxKey := fn(*oldValue)
			if idxKey != "" {
				// Remove key from slice
				keys := a.indexes[name][idxKey]
				for i, k := range keys {
					if k == key {
						a.indexes[name][idxKey] = append(keys[:i], keys[i+1:]...)
						break
					}
				}
			}
		}
	}

	if err := a.access.Update(ctx, key, value); err != nil {
		return err
	}

	// Add new index entries
	for name, fn := range a.indexFuncs {
		idxKey := fn(value)
		if idxKey != "" {
			a.indexes[name][idxKey] = append(a.indexes[name][idxKey], key)
		}
	}

	return nil
}

// Delete removes a value and its secondary index entries.
func (a *IndexedAccess[K, V]) Delete(ctx context.Context, key K) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Remove old index entries
	oldValue, err := a.access.Read(ctx, key)
	if err == nil && oldValue != nil {
		for name, fn := range a.indexFuncs {
			idxKey := fn(*oldValue)
			if idxKey != "" {
				// Remove key from slice
				keys := a.indexes[name][idxKey]
				for i, k := range keys {
					if k == key {
						a.indexes[name][idxKey] = append(keys[:i], keys[i+1:]...)
						break
					}
				}
			}
		}
	}

	return a.access.Delete(ctx, key)
}

// FindByIndex retrieves values by a secondary index key.
func (a *IndexedAccess[K, V]) FindByIndex(ctx context.Context, indexName string, indexKey string) ([]V, error) {
	a.mu.RLock()
	idx, ok := a.indexes[indexName]
	if !ok {
		a.mu.RUnlock()
		return nil, nil
	}
	keys, found := idx[indexKey]
	a.mu.RUnlock()

	if !found || len(keys) == 0 {
		return nil, nil
	}

	results := make([]V, 0, len(keys))
	for _, key := range keys {
		val, err := a.access.Read(ctx, key)
		if err == nil && val != nil {
			results = append(results, *val)
		}
	}
	return results, nil
}

// FindOneByIndex retrieves a single value by a secondary index key.
// Returns nil, false if not found.
func (a *IndexedAccess[K, V]) FindOneByIndex(ctx context.Context, indexName string, indexKey string) (*V, bool) {
	results, err := a.FindByIndex(ctx, indexName, indexKey)
	if err != nil || len(results) == 0 {
		return nil, false
	}
	return &results[0], true
}
