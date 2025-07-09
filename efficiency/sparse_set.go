package efficiency

// SparseSet is a data structure that provides efficient storage and retrieval of elements.
type SparseSet[T comparable] struct {
	dense  []T
	sparse []int
	size   int
}

// NewSparseSet creates a new SparseSet with the given initial size.
func NewSparseSet[T comparable](initialSize int) *SparseSet[T] {
	return &SparseSet[T]{
		dense:  make([]T, initialSize),
		sparse: make([]int, initialSize),
		size:   0,
	}
}

// Add adds an element to the SparseSet.
func (a *SparseSet[T]) Add(id int, item T) {
	if id < 0 || id > a.size {
		return
	}
	index := a.size
	a.dense[index] = item
	a.sparse[id] = index
	a.size++
}

// Dense returns the dense representation of the SparseSet.
func (a *SparseSet[T]) Dense() []T {
	return a.dense[:a.size]
}

// Remove removes an element from the SparseSet.
func (a *SparseSet[T]) Remove(id int) {
	if id < 0 || id >= len(a.sparse) {
		return
	}
	index := a.sparse[id]
	a.dense[index] = a.dense[a.size-1]
	// Correct the sparse index of the last element
	for i := range a.sparse {
		if a.sparse[i] == a.size-1 {
			a.sparse[i] = index
			break
		}
	}
	a.size--
}

// Size returns the size of the SparseSet.
func (a *SparseSet[T]) Size() int {
	return a.size
}
