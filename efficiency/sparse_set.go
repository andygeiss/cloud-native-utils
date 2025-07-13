package efficiency

import (
	"bytes"
	"encoding/gob"
)

// SparseSet is a data structure that provides efficient storage and retrieval of elements.
type SparseSet[T comparable] struct {
	Dense  []T
	Sparse []int
	Size   int
}

// NewSparseSet creates a new SparseSet with the given initial Size.
func NewSparseSet[T comparable](initialSize int) *SparseSet[T] {
	return &SparseSet[T]{
		Dense:  make([]T, initialSize),
		Sparse: make([]int, initialSize),
		Size:   0,
	}
}

// Add adds an element to the SparseSet.
func (a *SparseSet[T]) Add(id int, item T) {
	if id < 0 || id > a.Size {
		return
	}
	index := a.Size
	a.Dense[index] = item
	a.Sparse[id] = index
	a.Size++
}

// Densed returns the Dense representation of the SparseSet.
func (a *SparseSet[T]) Densed() []T {
	return a.Dense[:a.Size]
}

// Deserialize deserializes the SparseSet from a byte slice.
func (a *SparseSet[T]) Deserialize(in []byte) error {
	buf := bytes.NewBuffer(in)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(a)
}

// Remove removes an element from the SparseSet.
func (a *SparseSet[T]) Remove(id int) {
	if id < 0 || id >= len(a.Sparse) {
		return
	}
	index := a.Sparse[id]
	a.Dense[index] = a.Dense[a.Size-1]
	// Correct the Sparse index of the last element
	for i := range a.Sparse {
		if a.Sparse[i] == a.Size-1 {
			a.Sparse[i] = index
			break
		}
	}
	a.Size--
}

// Serialize serializes the SparseSet into a byte slice.
func (a *SparseSet[T]) Serialize() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	_ = encoder.Encode(a)
	return buf.Bytes()
}
