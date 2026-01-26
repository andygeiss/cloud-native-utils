package resource

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"sync"
)

// sparseShard holds one partition of the sharded store with sparse-dense storage.
// The bidirectional mapping (keyToIndex + indexToKey) enables O(1) swap-remove.
type sparseShard[K comparable, V any] struct {
	keyToIndex map[K]int // Sparse: key -> index in dense
	indexToKey []K       // Reverse mapping: index -> key (for O(1) swap-remove)
	dense      []V       // Dense: contiguous value storage
	size       int       // Number of active elements
	mu         sync.RWMutex
}

// ShardedSparseAccess is a high-performance in-memory Access implementation
// using sharding for reduced lock contention and sparse-dense storage for
// cache-friendly iteration.
type ShardedSparseAccess[K comparable, V any] struct {
	shards    []sparseShard[K, V]
	numShards int
}

// NewShardedSparseAccess creates a new sharded sparse access with the given
// number of shards. A typical value is 32 or runtime.NumCPU() * 4.
// If numShards <= 0, defaults to 32.
func NewShardedSparseAccess[K comparable, V any](numShards int) *ShardedSparseAccess[K, V] {
	if numShards <= 0 {
		numShards = 32
	}
	shards := make([]sparseShard[K, V], numShards)
	for i := range shards {
		shards[i] = sparseShard[K, V]{
			keyToIndex: make(map[K]int),
			indexToKey: make([]K, 0, 64),
			dense:      make([]V, 0, 64),
		}
	}
	return &ShardedSparseAccess[K, V]{
		shards:    shards,
		numShards: numShards,
	}
}

// NewShardedSparseAccessWithCapacity creates a new sharded sparse access with
// pre-allocated capacity per shard for better performance when size is known.
func NewShardedSparseAccessWithCapacity[K comparable, V any](numShards, capacityPerShard int) *ShardedSparseAccess[K, V] {
	if numShards <= 0 {
		numShards = 32
	}
	if capacityPerShard <= 0 {
		capacityPerShard = 64
	}
	shards := make([]sparseShard[K, V], numShards)
	for i := range shards {
		shards[i] = sparseShard[K, V]{
			keyToIndex: make(map[K]int, capacityPerShard),
			indexToKey: make([]K, 0, capacityPerShard),
			dense:      make([]V, 0, capacityPerShard),
		}
	}
	return &ShardedSparseAccess[K, V]{
		shards:    shards,
		numShards: numShards,
	}
}

// Clear removes all elements from all shards.
func (a *ShardedSparseAccess[K, V]) Clear() {
	for i := range a.shards {
		shard := &a.shards[i]
		shard.mu.Lock()
		shard.keyToIndex = make(map[K]int)
		shard.indexToKey = shard.indexToKey[:0]
		shard.dense = shard.dense[:0]
		shard.size = 0
		shard.mu.Unlock()
	}
}

// Create creates a new resource.
func (a *ShardedSparseAccess[K, V]) Create(ctx context.Context, key K, value V) error {
	// Skip if context is canceled or timed out.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	shard := a.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	// Check if key already exists.
	if _, exists := shard.keyToIndex[key]; exists {
		return errors.New(ErrorResourceAlreadyExists)
	}

	// Add to sparse-dense structure.
	idx := shard.size
	shard.keyToIndex[key] = idx
	shard.indexToKey = append(shard.indexToKey, key)
	shard.dense = append(shard.dense, value)
	shard.size++

	return nil
}

// Delete deletes a resource.
func (a *ShardedSparseAccess[K, V]) Delete(ctx context.Context, key K) error {
	// Skip if context is canceled or timed out.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	shard := a.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	idx, exists := shard.keyToIndex[key]
	if !exists {
		return errors.New(ErrorResourceNotFound)
	}

	// O(1) swap-remove.
	lastIdx := shard.size - 1
	if idx != lastIdx {
		// Move last element to the deleted position.
		lastKey := shard.indexToKey[lastIdx]
		lastValue := shard.dense[lastIdx]

		shard.dense[idx] = lastValue
		shard.indexToKey[idx] = lastKey
		shard.keyToIndex[lastKey] = idx
	}

	// Remove last element.
	delete(shard.keyToIndex, key)
	shard.indexToKey = shard.indexToKey[:lastIdx]
	shard.dense = shard.dense[:lastIdx]
	shard.size--

	return nil
}

// ForEach iterates over all elements. Stops if fn returns false.
// Iteration order is not guaranteed.
func (a *ShardedSparseAccess[K, V]) ForEach(fn func(K, V) bool) {
	for i := range a.shards {
		shard := &a.shards[i]
		shard.mu.RLock()
		for j := 0; j < shard.size; j++ {
			if !fn(shard.indexToKey[j], shard.dense[j]) {
				shard.mu.RUnlock()
				return
			}
		}
		shard.mu.RUnlock()
	}
}

// getShard returns the shard for a given key using FNV-1a hash.
func (a *ShardedSparseAccess[K, V]) getShard(key K) *sparseShard[K, V] {
	hash := fnv.New32a()
	_, _ = hash.Write(fmt.Appendf(nil, "%v", key))
	idx := int(hash.Sum32()) % a.numShards
	return &a.shards[idx]
}

// Len returns the total number of elements across all shards.
func (a *ShardedSparseAccess[K, V]) Len() int {
	total := 0
	for i := range a.shards {
		a.shards[i].mu.RLock()
		total += a.shards[i].size
		a.shards[i].mu.RUnlock()
	}
	return total
}

// Read reads a resource.
func (a *ShardedSparseAccess[K, V]) Read(ctx context.Context, key K) (*V, error) {
	// Skip if context is canceled or timed out.
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	shard := a.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	idx, exists := shard.keyToIndex[key]
	if !exists {
		return nil, errors.New(ErrorResourceNotFound)
	}

	return &shard.dense[idx], nil
}

// ReadAll reads all resources.
func (a *ShardedSparseAccess[K, V]) ReadAll(ctx context.Context) ([]V, error) {
	// Skip if context is canceled or timed out.
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Pre-calculate total size for efficient allocation.
	totalSize := 0
	for i := range a.shards {
		a.shards[i].mu.RLock()
		totalSize += a.shards[i].size
		a.shards[i].mu.RUnlock()
	}

	result := make([]V, 0, totalSize)

	// Collect from all shards - lock each shard briefly.
	for i := range a.shards {
		shard := &a.shards[i]
		shard.mu.RLock()
		// Cache-friendly: copy from contiguous dense array.
		result = append(result, shard.dense[:shard.size]...)
		shard.mu.RUnlock()

		// Check context between shards for responsiveness.
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return result, nil
}

// Update updates a resource.
func (a *ShardedSparseAccess[K, V]) Update(ctx context.Context, key K, value V) error {
	// Skip if context is canceled or timed out.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	shard := a.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	idx, exists := shard.keyToIndex[key]
	if !exists {
		return errors.New(ErrorResourceNotFound)
	}

	shard.dense[idx] = value
	return nil
}
