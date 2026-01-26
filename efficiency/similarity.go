package efficiency

import (
	"container/heap"
	"context"
	"sync/atomic"
)

// Result represents a similarity search result with the key, value, and similarity score.
type Result[K comparable, V any] struct {
	Key   K
	Value *V
	Score float64
}

// SparseVectorProvider extracts a sparse vector representation from a value.
// Used for cosine similarity on TF-IDF vectors.
type SparseVectorProvider interface {
	// SparseVector returns sorted term indices, corresponding TF-IDF values, and pre-computed L2 norm.
	// Indices MUST be sorted in ascending order for merge-loop efficiency.
	SparseVector() (indices []int, values []float64, norm float64)
}

// SparseSetProvider extracts a sparse set (binary term presence) from a value.
// Used for Jaccard similarity on token sets.
type SparseSetProvider interface {
	// SparseSet returns sorted term indices representing set membership.
	// Indices MUST be sorted in ascending order for merge-loop efficiency.
	SparseSet() []int
}

// SearchOptions configures similarity search behavior.
type SearchOptions struct {
	// TopK limits results to the K most similar items. Default: 10.
	TopK int

	// Threshold filters results below this similarity score. Default: 0.0.
	Threshold float64

	// EarlyStopCount stops after finding this many results above threshold.
	// Zero means no early stopping. Default: 0.
	EarlyStopCount int
}

// DefaultSearchOptions returns sensible defaults.
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		TopK:           10,
		Threshold:      0.0,
		EarlyStopCount: 0,
	}
}

// SearchContext creates a stopped flag that respects context cancellation.
// Use this to integrate with context-based cancellation since ForEach does not accept context.
func SearchContext(ctx context.Context) *atomic.Bool {
	stopped := &atomic.Bool{}
	go func() {
		<-ctx.Done()
		stopped.Store(true)
	}()
	return stopped
}

// resultHeap is a min-heap for tracking top-K results by score.
// Lower scores are at the top, allowing efficient removal of the smallest.
type resultHeap[K comparable, V any] []Result[K, V]

func (h resultHeap[K, V]) Len() int           { return len(h) }
func (h resultHeap[K, V]) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h resultHeap[K, V]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *resultHeap[K, V]) Push(x any) {
	*h = append(*h, x.(Result[K, V]))
}

func (h *resultHeap[K, V]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// cosineSimilarity computes cosine similarity between two sparse vectors using merge-loop.
// Both index slices MUST be sorted in ascending order.
// Formula: cosine(A, B) = dot(A, B) / (norm(A) * norm(B)).
func cosineSimilarity(indicesA, indicesB []int, valuesA, valuesB []float64, normA, normB float64) float64 {
	if normA == 0 || normB == 0 {
		return 0.0
	}

	var dot float64
	i, j := 0, 0

	for i < len(indicesA) && j < len(indicesB) {
		switch {
		case indicesA[i] == indicesB[j]:
			dot += valuesA[i] * valuesB[j]
			i++
			j++
		case indicesA[i] < indicesB[j]:
			i++
		default:
			j++
		}
	}

	return dot / (normA * normB)
}

// jaccardSimilarity computes Jaccard similarity between two sparse sets using merge-loop.
// Both index slices MUST be sorted in ascending order.
// Formula: jaccard(A, B) = |A ∩ B| / |A ∪ B|.
func jaccardSimilarity(setA, setB []int) float64 {
	if len(setA) == 0 && len(setB) == 0 {
		return 1.0 // Both empty sets are identical
	}
	if len(setA) == 0 || len(setB) == 0 {
		return 0.0
	}

	var intersection, union int
	i, j := 0, 0

	for i < len(setA) && j < len(setB) {
		switch {
		case setA[i] == setB[j]:
			intersection++
			union++
			i++
			j++
		case setA[i] < setB[j]:
			union++
			i++
		default:
			union++
			j++
		}
	}

	// Count remaining elements
	union += (len(setA) - i) + (len(setB) - j)

	return float64(intersection) / float64(union)
}

// addToHeap adds a result to the min-heap, maintaining at most topK elements.
func addToHeap[K comparable, V any](h *resultHeap[K, V], key K, value *V, score float64, topK int) {
	if h.Len() < topK {
		heap.Push(h, Result[K, V]{Key: key, Value: value, Score: score})
	} else if score > (*h)[0].Score {
		// Replace the smallest if this score is better
		heap.Pop(h)
		heap.Push(h, Result[K, V]{Key: key, Value: value, Score: score})
	}
}

// extractSortedResults extracts results from the heap sorted by descending score.
func extractSortedResults[K comparable, V any](h *resultHeap[K, V]) []Result[K, V] {
	results := make([]Result[K, V], h.Len())
	for i := len(results) - 1; i >= 0; i-- {
		results[i] = heap.Pop(h).(Result[K, V])
	}
	return results
}

// FindSimilarCosine searches for items similar to the query using cosine similarity.
// V must implement SparseVectorProvider to extract TF-IDF sparse vectors.
// Returns results sorted by descending similarity score.
//
// The function iterates over the store using ForEach, computing similarity
// for each item and maintaining a top-K heap for efficient result collection.
//
// Use stopped flag (from SearchContext) to enable cancellation since ForEach
// does not accept context. Pass nil if cancellation is not needed.
func FindSimilarCosine[K comparable, V SparseVectorProvider](
	store interface{ ForEach(fn func(K, V) bool) },
	query V,
	opts SearchOptions,
	stopped *atomic.Bool,
) []Result[K, V] {
	if opts.TopK <= 0 {
		opts.TopK = 10
	}

	// Extract query vector once
	queryIndices, queryValues, queryNorm := query.SparseVector()
	if queryNorm == 0 {
		return nil // Query has no content
	}

	// Min-heap for top-K tracking
	h := &resultHeap[K, V]{}
	heap.Init(h)

	earlyStopFound := 0

	store.ForEach(func(key K, value V) bool {
		// Check for cancellation
		if stopped != nil && stopped.Load() {
			return false
		}

		// Extract candidate vector
		candIndices, candValues, candNorm := value.SparseVector()
		if candNorm == 0 {
			return true // Skip empty vectors
		}

		// Compute cosine similarity using merge loop
		score := cosineSimilarity(queryIndices, candIndices, queryValues, candValues, queryNorm, candNorm)

		// Apply threshold filter
		if score < opts.Threshold {
			return true
		}

		// Add to heap
		valueCopy := value
		addToHeap(h, key, &valueCopy, score, opts.TopK)

		// Track for early stop
		if opts.EarlyStopCount > 0 {
			earlyStopFound++
			if earlyStopFound >= opts.EarlyStopCount {
				return false
			}
		}

		return true
	})

	return extractSortedResults(h)
}

// FindSimilarJaccard searches for items similar to the query using Jaccard similarity.
// V must implement SparseSetProvider to extract term sets.
// Returns results sorted by descending similarity score.
//
// The function iterates over the store using ForEach, computing similarity
// for each item and maintaining a top-K heap for efficient result collection.
//
// Use stopped flag (from SearchContext) to enable cancellation since ForEach
// does not accept context. Pass nil if cancellation is not needed.
func FindSimilarJaccard[K comparable, V SparseSetProvider](
	store interface{ ForEach(fn func(K, V) bool) },
	query V,
	opts SearchOptions,
	stopped *atomic.Bool,
) []Result[K, V] {
	if opts.TopK <= 0 {
		opts.TopK = 10
	}

	// Extract query set once
	querySet := query.SparseSet()
	if len(querySet) == 0 {
		return nil
	}

	// Min-heap for top-K tracking
	h := &resultHeap[K, V]{}
	heap.Init(h)

	earlyStopFound := 0

	store.ForEach(func(key K, value V) bool {
		// Check for cancellation
		if stopped != nil && stopped.Load() {
			return false
		}

		// Extract candidate set
		candSet := value.SparseSet()
		if len(candSet) == 0 {
			return true // Skip empty sets
		}

		// Compute Jaccard similarity using merge loop
		score := jaccardSimilarity(querySet, candSet)

		// Apply threshold filter
		if score < opts.Threshold {
			return true
		}

		// Add to heap
		valueCopy := value
		addToHeap(h, key, &valueCopy, score, opts.TopK)

		// Track for early stop
		if opts.EarlyStopCount > 0 {
			earlyStopFound++
			if earlyStopFound >= opts.EarlyStopCount {
				return false
			}
		}

		return true
	})

	return extractSortedResults(h)
}
