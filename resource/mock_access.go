package resource

// MockAccess is a mock implementation of Access[K, V].
// Each method has a corresponding mock function pointer (...Ptr).
// Use the builder-pattern to set the mock function pointers for each method.
// This allows for more flexible and readable test cases.
type MockAccess[K, V any] struct {
	CreatePtr  func(key K, value V) error
	ReadPtr    func(key K) (*V, error)
	ReadAllPtr func() ([]V, error)
	UpdatePtr  func(key K, value V) error
	DeletePtr  func(key K) error
}

// NewMockAccess creates a new instance of MockAccess[K, V].
func NewMockAccess[K, V any]() *MockAccess[K, V] {
	return &MockAccess[K, V]{}
}

// Create creates a new resource with the given key and value.
func (a *MockAccess[K, V]) Create(key K, value V) error {
	return a.CreatePtr(key, value)
}

// Read reads a resource with the given key.
func (a *MockAccess[K, V]) Read(key K) (value *V, err error) {
	return a.ReadPtr(key)
}

// ReadAll reads all resources.
func (a *MockAccess[K, V]) ReadAll() (values []V, err error) {
	return a.ReadAllPtr()
}

// Update updates a resource with the given key and value.
func (a *MockAccess[K, V]) Update(key K, value V) (err error) {
	return a.UpdatePtr(key, value)
}

// Delete deletes a resource with the given key.
func (a *MockAccess[K, V]) Delete(key K) (err error) {
	return a.DeletePtr(key)
}

// WithCreatePtr sets the mock function pointer for creating a resource.
func (a *MockAccess[K, V]) WithCreatePtr(fn func(key K, value V) error) {
	a.CreatePtr = fn
}

// WithReadPtr sets the mock function pointer for reading a resource.
func (a *MockAccess[K, V]) WithReadPtr(fn func(key K) (*V, error)) {
	a.ReadPtr = fn
}

// WithReadAllPtr sets the mock function pointer for reading all resources.
func (a *MockAccess[K, V]) WithReadAllPtr(fn func() ([]V, error)) {
	a.ReadAllPtr = fn
}

// WithUpdatePtr sets the mock function pointer for updating a resource.
func (a *MockAccess[K, V]) WithUpdatePtr(fn func(key K, value V) error) {
	a.UpdatePtr = fn
}

// WithDeletePtr sets the mock function pointer for deleting a resource.
func (a *MockAccess[K, V]) WithDeletePtr(fn func(key K) error) {
	a.DeletePtr = fn
}
