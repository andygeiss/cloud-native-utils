package efficiency

// SparseSet is a data structure that provides efficient storage and retrieval of elements.
type SparseSet[T comparable] struct {
	dense  []T
	sparse []int
	lastId int
	size   int
}

// NewSparseSet creates a new SparseSet with the given initial size.
func NewSparseSet[T comparable](initialSize int) *SparseSet[T] {
	return &SparseSet[T]{
		dense:  make([]T, initialSize),
		sparse: make([]int, initialSize),
		lastId: 0,
		size:   0,
	}
}

// Add adds an element to the SparseSet.
func (a *SparseSet[T]) Add(id int, item T) {
	if id < 0 {
		return
	}
	index := a.size
	a.dense[index] = item
	a.sparse[id] = index
	a.lastId = id
	a.size++
}

// Dense returns the dense representation of the SparseSet.
func (a *SparseSet[T]) Dense() []T {
	return a.dense[:a.size]
}

// Remove removes an element from the SparseSet.
func (a *SparseSet[T]) Remove(id int) {
	if id < 0 || id > a.lastId {
		return
	}
	index := a.sparse[id]
	a.dense[index] = a.dense[a.size-1]
	a.sparse[a.lastId] = index
	a.size--
}

// Size returns the size of the SparseSet.
func (a *SparseSet[T]) Size() int {
	return a.size
}
