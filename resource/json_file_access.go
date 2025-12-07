package resource

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// jsonFileAccess is a json file access.
type jsonFileAccess[K comparable, V any] struct {
	path  string
	mutex sync.RWMutex
}

// NewJsonFileAccess creates a new json file access.
func NewJsonFileAccess[K comparable, V any](path string) *jsonFileAccess[K, V] {
	return &jsonFileAccess[K, V]{path: path}
}

// Create creates a new resource.
func (a *jsonFileAccess[K, V]) Create(ctx context.Context, key K, value V) error {
	// Skip if context is canceled or timed out.
	if err := ctx.Err(); err != nil {
		return err
	}

	// Ensure that only one goroutine can write to the map at a time.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Read data from file.
	data, err := fromJsonFile[K, V](a.path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// Ensure data is not nil.
	if data == nil {
		data = make(map[K]V)
	}

	// Check if resource exists.
	if _, alreadyExists := data[key]; alreadyExists {
		return errors.New(ErrorResourceAlreadyExists)
	}

	// Set resource if not exists.
	data[key] = value

	// Write data to file.
	if err := intoJsonFile[K, V](a.path, data); err != nil {
		return err
	}

	return nil
}

// Read reads a resource.
func (a *jsonFileAccess[K, V]) Read(ctx context.Context, key K) (*V, error) {
	// Skip if context is canceled or timed out.
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Ensure that read only access is allowed.
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// Read data from file.
	data, err := fromJsonFile[K, V](a.path)
	if err != nil {
		return nil, err
	}

	// Check if resource exists.
	value, exists := data[key]
	if !exists {
		return nil, errors.New(ErrorResourceNotFound)
	}

	return &value, nil
}

// ReadAll reads all resources.
func (a *jsonFileAccess[K, V]) ReadAll(ctx context.Context) ([]V, error) {
	// Skip if context is canceled or timed out.
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Ensure that read only access is allowed.
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// Read data from file.
	data, err := fromJsonFile[K, V](a.path)
	if err != nil {
		return nil, err
	}

	// Ensure data is not nil.
	if data == nil {
		data = make(map[K]V)
	}

	// Convert data to values.
	var values []V
	for _, value := range data {
		values = append(values, value)
	}

	return values, nil
}

// Update updates a resource.
func (a *jsonFileAccess[K, V]) Update(ctx context.Context, key K, value V) error {
	// Skip if context is canceled or timed out.
	if err := ctx.Err(); err != nil {
		return err
	}

	// Ensure that only one goroutine can write to the file.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Read data from file.
	data, err := fromJsonFile[K, V](a.path)
	if err != nil {
		return err
	}

	// Ensure data is not nil.
	if data == nil {
		data = make(map[K]V)
	}

	// Update resource if exists.
	if _, exists := data[key]; exists {
		data[key] = value
	} else {
		return errors.New(ErrorResourceNotFound)
	}

	// Write data to file.
	if err := intoJsonFile[K, V](a.path, data); err != nil {
		return err
	}

	return nil
}

// Delete deletes a resource.
func (a *jsonFileAccess[K, V]) Delete(ctx context.Context, key K) error {
	// Skip if context is canceled or timed out.
	if err := ctx.Err(); err != nil {
		return err
	}

	// Ensure that only one goroutine can write to the file.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Read data from file.
	data, err := fromJsonFile[K, V](a.path)
	if err != nil {
		return err
	}

	// Check if resource exists.
	_, exists := data[key]
	if !exists {
		return errors.New(ErrorResourceNotFound)
	}

	delete(data, key)

	// Write data to file.
	if err := intoJsonFile[K, V](a.path, data); err != nil {
		return err
	}

	return nil
}

func fromJsonFile[K comparable, V any](path string) (map[K]V, error) {
	var values map[K]V
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, err
	}
	return values, nil
}

func intoJsonFile[K comparable, V any](path string, values map[K]V) error {
	data, err := json.Marshal(values)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}
