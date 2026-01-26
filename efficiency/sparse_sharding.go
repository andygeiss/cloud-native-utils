package efficiency

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// SparseShard is a single partition using KeyedSparseSet for O(1) operations.
type SparseShard[K comparable, V any] struct {
	set *KeyedSparseSet[K, V]
	mu  sync.RWMutex
}

// SparseSharding distributes key-value pairs across multiple KeyedSparseSets.
// It provides per-shard locking for high-concurrency workloads and O(1) delete.
type SparseSharding[K comparable, V any] struct {
	shards    []SparseShard[K, V]
	numShards int
}

// NewSparseSharding creates a new SparseSharding with the given number of shards.
// A typical value is 32 or runtime.NumCPU() * 4.
func NewSparseSharding[K comparable, V any](numShards int) *SparseSharding[K, V] {
	if numShards <= 0 {
		numShards = 32
	}
	shards := make([]SparseShard[K, V], numShards)
	for i := range shards {
		shards[i] = SparseShard[K, V]{
			set: NewKeyedSparseSet[K, V](64),
		}
	}
	return &SparseSharding[K, V]{
		shards:    shards,
		numShards: numShards,
	}
}

// NewSparseShardingWithCapacity creates a new SparseSharding with pre-allocated capacity per shard.
func NewSparseShardingWithCapacity[K comparable, V any](numShards, capacityPerShard int) *SparseSharding[K, V] {
	if numShards <= 0 {
		numShards = 32
	}
	if capacityPerShard <= 0 {
		capacityPerShard = 64
	}
	shards := make([]SparseShard[K, V], numShards)
	for i := range shards {
		shards[i] = SparseShard[K, V]{
			set: NewKeyedSparseSet[K, V](capacityPerShard),
		}
	}
	return &SparseSharding[K, V]{
		shards:    shards,
		numShards: numShards,
	}
}

// Clear removes all elements from all shards.
func (s *SparseSharding[K, V]) Clear() {
	for i := range s.shards {
		shard := &s.shards[i]
		shard.mu.Lock()
		shard.set.Clear()
		shard.mu.Unlock()
	}
}

// Delete removes a key-value pair. Returns true if the key existed.
func (s *SparseSharding[K, V]) Delete(key K) bool {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	return shard.set.Delete(key)
}

// ForEach iterates over all elements across all shards. Stops if fn returns false.
// Iteration order is not guaranteed.
func (s *SparseSharding[K, V]) ForEach(fn func(K, V) bool) {
	for i := range s.shards {
		shard := &s.shards[i]
		shard.mu.RLock()
		stopped := false
		shard.set.ForEach(func(k K, v V) bool {
			if !fn(k, v) {
				stopped = true
				return false
			}
			return true
		})
		shard.mu.RUnlock()
		if stopped {
			return
		}
	}
}

// ForEachShard iterates over each shard with the shard index.
// The callback receives the shard index and a function to iterate within that shard.
// Useful for checking cancellation between shards.
func (s *SparseSharding[K, V]) ForEachShard(fn func(shardIdx int, iterate func(fn func(K, V) bool))) {
	for i := range s.shards {
		shard := &s.shards[i]
		fn(i, func(innerFn func(K, V) bool) {
			shard.mu.RLock()
			defer shard.mu.RUnlock()
			shard.set.ForEach(innerFn)
		})
	}
}

// Get retrieves a value by key. Returns nil if not found.
func (s *SparseSharding[K, V]) Get(key K) *V {
	shard := s.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	return shard.set.Get(key)
}

// Has returns true if the key exists.
func (s *SparseSharding[K, V]) Has(key K) bool {
	shard := s.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	return shard.set.Has(key)
}

// Len returns the total number of elements across all shards.
func (s *SparseSharding[K, V]) Len() int {
	total := 0
	for i := range s.shards {
		s.shards[i].mu.RLock()
		total += s.shards[i].set.Len()
		s.shards[i].mu.RUnlock()
	}
	return total
}

// Put adds or updates a key-value pair. Returns true if the key was new.
func (s *SparseSharding[K, V]) Put(key K, value V) bool {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	return shard.set.Put(key, value)
}

// Values returns all values across all shards.
func (s *SparseSharding[K, V]) Values() []V {
	// Pre-calculate total size for efficient allocation.
	totalSize := 0
	for i := range s.shards {
		s.shards[i].mu.RLock()
		totalSize += s.shards[i].set.Len()
		s.shards[i].mu.RUnlock()
	}

	result := make([]V, 0, totalSize)
	for i := range s.shards {
		shard := &s.shards[i]
		shard.mu.RLock()
		result = append(result, shard.set.Values()...)
		shard.mu.RUnlock()
	}
	return result
}

// getShard returns the shard for a given key using FNV-1a hash.
func (s *SparseSharding[K, V]) getShard(key K) *SparseShard[K, V] {
	hash := fnv.New32a()
	_, _ = hash.Write(fmt.Appendf(nil, "%v", key))
	idx := int(hash.Sum32()) % s.numShards
	return &s.shards[idx]
}
