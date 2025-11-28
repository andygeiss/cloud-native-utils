package resource

// mockAccess is a mock implementation of Access[K, V].
// Each method has a corresponding mock function pointer (...Ptr).
// Use the builder-pattern to set the mock function pointers for each method.
// This allows for more flexible and readable test cases.
type mockAccess[K, V any] struct {
	createFn  func(key K, value V) error
	readFn    func(key K) (*V, error)
	readAllFn func() ([]V, error)
	updateFn  func(key K, value V) error
	deleteFn  func(key K) error
}

// NewMockAccess creates a new instance of MockAccess[K, V].
func NewMockAccess[K, V any]() *mockAccess[K, V] {
	return &mockAccess[K, V]{}
}

// Create creates a new resource with the given key and value.
func (a *mockAccess[K, V]) Create(key K, value V) error {
	return a.createFn(key, value)
}

// Read reads a resource with the given key.
func (a *mockAccess[K, V]) Read(key K) (value *V, err error) {
	return a.readFn(key)
}

// ReadAll reads all resources.
func (a *mockAccess[K, V]) ReadAll() (values []V, err error) {
	return a.readAllFn()
}

// Update updates a resource with the given key and value.
func (a *mockAccess[K, V]) Update(key K, value V) (err error) {
	return a.updateFn(key, value)
}

// Delete deletes a resource with the given key.
func (a *mockAccess[K, V]) Delete(key K) (err error) {
	return a.deleteFn(key)
}

// WithCreateFn sets the mock function pointer for creating a resource.
func (a *mockAccess[K, V]) WithCreateFn(fn func(key K, value V) error) *mockAccess[K, V] {
	a.createFn = fn
	return a
}

// WithReadFn sets the mock function pointer for reading a resource.
func (a *mockAccess[K, V]) WithReadFn(fn func(key K) (*V, error)) *mockAccess[K, V] {
	a.readFn = fn
	return a
}

// WithReadAllFn sets the mock function pointer for reading all resources.
func (a *mockAccess[K, V]) WithReadAllFn(fn func() ([]V, error)) *mockAccess[K, V] {
	a.readAllFn = fn
	return a
}

// WithUpdateFn sets the mock function pointer for updating a resource.
func (a *mockAccess[K, V]) WithUpdateFn(fn func(key K, value V) error) *mockAccess[K, V] {
	a.updateFn = fn
	return a
}

// WithDeleteFn sets the mock function pointer for deleting a resource.
func (a *mockAccess[K, V]) WithDeleteFn(fn func(key K) error) *mockAccess[K, V] {
	a.deleteFn = fn
	return a
}
