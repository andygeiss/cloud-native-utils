package resource

// Access specifies the CRUD operations for a resource using generics.
type Access[K, V any] interface {
	Create(key K, model V) error
	Read(key K) (*V, error)
	Update(key K, model V) error
	Delete(key K) error
}
