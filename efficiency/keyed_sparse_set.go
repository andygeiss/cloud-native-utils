package efficiency

// KeyedSparseSet is a generic key-value sparse set with O(1) operations.
// It uses bidirectional mapping (key->index, index->key) for O(1) swap-remove.
type KeyedSparseSet[K comparable, V any] struct {
	keyToIndex map[K]int // Sparse: key -> index in dense
	indexToKey []K       // Reverse mapping: index -> key (for O(1) swap-remove)
	dense      []V       // Dense: contiguous value storage
	size       int       // Number of active elements
}

// NewKeyedSparseSet creates a new KeyedSparseSet with the given initial capacity.
func NewKeyedSparseSet[K comparable, V any](capacity int) *KeyedSparseSet[K, V] {
	if capacity <= 0 {
		capacity = 64
	}
	return &KeyedSparseSet[K, V]{
		keyToIndex: make(map[K]int, capacity),
		indexToKey: make([]K, 0, capacity),
		dense:      make([]V, 0, capacity),
		size:       0,
	}
}

// Clear removes all elements.
func (s *KeyedSparseSet[K, V]) Clear() {
	s.keyToIndex = make(map[K]int)
	s.indexToKey = s.indexToKey[:0]
	s.dense = s.dense[:0]
	s.size = 0
}

// Delete removes a key-value pair. Returns true if the key existed.
// Uses O(1) swap-remove: swaps the deleted element with the last element.
func (s *KeyedSparseSet[K, V]) Delete(key K) bool {
	idx, exists := s.keyToIndex[key]
	if !exists {
		return false
	}

	lastIdx := s.size - 1
	if idx != lastIdx {
		// Move last element to the deleted position.
		lastKey := s.indexToKey[lastIdx]
		lastValue := s.dense[lastIdx]

		s.dense[idx] = lastValue
		s.indexToKey[idx] = lastKey
		s.keyToIndex[lastKey] = idx
	}

	// Remove last element.
	delete(s.keyToIndex, key)
	s.indexToKey = s.indexToKey[:lastIdx]
	s.dense = s.dense[:lastIdx]
	s.size--

	return true
}

// ForEach iterates over all elements. Stops if fn returns false.
func (s *KeyedSparseSet[K, V]) ForEach(fn func(K, V) bool) {
	for i := range s.size {
		if !fn(s.indexToKey[i], s.dense[i]) {
			return
		}
	}
}

// Get returns the value for a key, or nil if not found.
func (s *KeyedSparseSet[K, V]) Get(key K) *V {
	idx, exists := s.keyToIndex[key]
	if !exists {
		return nil
	}
	return &s.dense[idx]
}

// Has returns true if the key exists.
func (s *KeyedSparseSet[K, V]) Has(key K) bool {
	_, exists := s.keyToIndex[key]
	return exists
}

// Len returns the number of elements.
func (s *KeyedSparseSet[K, V]) Len() int {
	return s.size
}

// Put adds or updates a key-value pair. Returns true if the key was new.
func (s *KeyedSparseSet[K, V]) Put(key K, value V) bool {
	if idx, exists := s.keyToIndex[key]; exists {
		s.dense[idx] = value
		return false
	}

	idx := s.size
	s.keyToIndex[key] = idx
	s.indexToKey = append(s.indexToKey, key)
	s.dense = append(s.dense, value)
	s.size++
	return true
}

// Values returns a slice of all values (dense array).
func (s *KeyedSparseSet[K, V]) Values() []V {
	return s.dense[:s.size]
}
