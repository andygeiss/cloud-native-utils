package resource

const (
	ErrorResourceAlreadyExists = "resource already exists"
	ErrorResourceNotFound      = "resource not found"
)

// Access specifies the CRUD operations for a resource using generics.
type Access[K, V any] interface {
	Create(key K, value V) error
	Read(key K) (*V, error)
	Update(key K, value V) error
	Delete(key K) error
}
