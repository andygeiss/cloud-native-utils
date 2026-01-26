package resource

import (
	"container/heap"
	"context"
	"errors"
	"sort"

	"github.com/andygeiss/cloud-native-utils/efficiency"
)

// ShardedSparseAccess is a high-performance in-memory Access implementation
// using sharding for reduced lock contention and sparse-dense storage for
// cache-friendly iteration. It wraps efficiency.SparseSharding with context
// handling and CRUD error semantics.
type ShardedSparseAccess[K comparable, V any] struct {
	shards *efficiency.SparseSharding[K, V]
}

// NewShardedSparseAccess creates a new sharded sparse access with the given
// number of shards. A typical value is 32 or runtime.NumCPU() * 4.
// If numShards <= 0, defaults to 32.
func NewShardedSparseAccess[K comparable, V any](numShards int) *ShardedSparseAccess[K, V] {
	return &ShardedSparseAccess[K, V]{
		shards: efficiency.NewSparseSharding[K, V](numShards),
	}
}

// NewShardedSparseAccessWithCapacity creates a new sharded sparse access with
// pre-allocated capacity per shard for better performance when size is known.
func NewShardedSparseAccessWithCapacity[K comparable, V any](numShards, capacityPerShard int) *ShardedSparseAccess[K, V] {
	return &ShardedSparseAccess[K, V]{
		shards: efficiency.NewSparseShardingWithCapacity[K, V](numShards, capacityPerShard),
	}
}

// Clear removes all elements from all shards.
func (a *ShardedSparseAccess[K, V]) Clear() {
	a.shards.Clear()
}

// Create creates a new resource.
func (a *ShardedSparseAccess[K, V]) Create(ctx context.Context, key K, value V) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if a.shards.Has(key) {
		return errors.New(ErrorResourceAlreadyExists)
	}

	a.shards.Put(key, value)
	return nil
}

// Delete deletes a resource.
func (a *ShardedSparseAccess[K, V]) Delete(ctx context.Context, key K) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if !a.shards.Delete(key) {
		return errors.New(ErrorResourceNotFound)
	}
	return nil
}

// ForEach iterates over all elements. Stops if fn returns false.
// Iteration order is not guaranteed.
func (a *ShardedSparseAccess[K, V]) ForEach(fn func(K, V) bool) {
	a.shards.ForEach(fn)
}

// Len returns the total number of elements across all shards.
func (a *ShardedSparseAccess[K, V]) Len() int {
	return a.shards.Len()
}

// Read reads a resource.
func (a *ShardedSparseAccess[K, V]) Read(ctx context.Context, key K) (*V, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	value := a.shards.Get(key)
	if value == nil {
		return nil, errors.New(ErrorResourceNotFound)
	}
	return value, nil
}

// ReadAll reads all resources.
func (a *ShardedSparseAccess[K, V]) ReadAll(ctx context.Context) ([]V, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return a.shards.Values(), nil
}

// SearchOptions configures similarity search behavior.
type SearchOptions struct {
	// TopK limits results to the K most similar items. Default: 10.
	TopK int

	// Threshold filters results below this similarity score. Default: 0.0.
	Threshold float64
}

// SearchResult represents a similarity search result.
type SearchResult[K comparable, V any] struct {
	Key   K
	Value V
	Score float64
}

// searchHeap is a min-heap for tracking top-K results by score.
type searchHeap[K comparable, V any] []SearchResult[K, V]

func (h searchHeap[K, V]) Len() int           { return len(h) }
func (h searchHeap[K, V]) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h searchHeap[K, V]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *searchHeap[K, V]) Push(x any) {
	*h = append(*h, x.(SearchResult[K, V]))
}

func (h *searchHeap[K, V]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// SearchSimilar finds the top-K most similar values using a custom scoring function.
// The scorer function is called for each value and should return a similarity score.
// Higher scores indicate more similarity. Results are sorted by descending score.
//
// Example usage with cosine similarity:
//
//	results := store.SearchSimilar(ctx, func(doc Document) float64 {
//	    return efficiency.CosineSimilarity(
//	        query.Indices, doc.Indices,
//	        query.Values, doc.Values,
//	        query.Norm, doc.Norm,
//	    )
//	}, resource.SearchOptions{TopK: 10, Threshold: 0.5})
func (a *ShardedSparseAccess[K, V]) SearchSimilar(
	ctx context.Context,
	scorer func(V) float64,
	opts SearchOptions,
) []SearchResult[K, V] {
	if opts.TopK <= 0 {
		opts.TopK = 10
	}

	// Min-heap for top-K tracking.
	h := &searchHeap[K, V]{}
	heap.Init(h)

	// Iterate with cancellation check between shards.
	a.shards.ForEachShard(func(shardIdx int, iterate func(fn func(K, V) bool)) {
		// Check context between shards for responsiveness.
		if ctx.Err() != nil {
			return
		}

		iterate(func(key K, value V) bool {
			score := scorer(value)

			// Apply threshold filter.
			if score < opts.Threshold {
				return true
			}

			// Add to heap.
			if h.Len() < opts.TopK {
				heap.Push(h, SearchResult[K, V]{Key: key, Value: value, Score: score})
			} else if score > (*h)[0].Score {
				heap.Pop(h)
				heap.Push(h, SearchResult[K, V]{Key: key, Value: value, Score: score})
			}
			return true
		})
	})

	// Extract results sorted by descending score.
	results := make([]SearchResult[K, V], h.Len())
	for i := len(results) - 1; i >= 0; i-- {
		results[i] = heap.Pop(h).(SearchResult[K, V])
	}

	// Ensure stable descending order.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// Update updates a resource.
func (a *ShardedSparseAccess[K, V]) Update(ctx context.Context, key K, value V) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if !a.shards.Has(key) {
		return errors.New(ErrorResourceNotFound)
	}

	a.shards.Put(key, value)
	return nil
}
