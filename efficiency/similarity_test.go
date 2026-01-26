package efficiency_test

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
	"github.com/andygeiss/cloud-native-utils/resource"
)

// TestDocument implements both SparseVectorProvider and SparseSetProvider.
type TestDocument struct {
	ID      string
	Indices []int
	Values  []float64
	Norm    float64
}

func (d TestDocument) SparseVector() ([]int, []float64, float64) {
	return d.Indices, d.Values, d.Norm
}

func (d TestDocument) SparseSet() []int {
	return d.Indices
}

// computeNorm calculates L2 norm for test setup.
func computeNorm(values []float64) float64 {
	var sum float64
	for _, v := range values {
		sum += v * v
	}
	return math.Sqrt(sum)
}

// --- Cosine Similarity Tests ---

func Test_FindSimilarCosine_With_EmptyQuery_Should_ReturnNoResults(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	_ = store.Create(ctx, "doc1", TestDocument{Indices: []int{0, 1}, Values: []float64{0.8, 0.6}, Norm: 1.0})

	query := TestDocument{Indices: []int{}, Values: []float64{}, Norm: 0.0}
	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarCosine(store, query, opts, nil)

	// Assert
	assert.That(t, "results must be empty", len(results), 0)
}

func Test_FindSimilarCosine_With_IdenticalVectors_Should_ReturnScoreOne(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	values := []float64{0.6, 0.8}
	norm := computeNorm(values)
	doc := TestDocument{ID: "doc1", Indices: []int{0, 2}, Values: values, Norm: norm}
	_ = store.Create(ctx, "doc1", doc)

	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarCosine(store, doc, opts, nil)

	// Assert
	assert.That(t, "results len must be 1", len(results), 1)
	assert.That(t, "score must be approximately 1.0", results[0].Score > 0.999, true)
}

func Test_FindSimilarCosine_With_OrthogonalVectors_Should_ReturnScoreZero(t *testing.T) {
	// Arrange: vectors with no overlapping indices have dot product 0
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	doc1 := TestDocument{ID: "doc1", Indices: []int{0, 1}, Values: []float64{1.0, 0.0}, Norm: 1.0}
	doc2 := TestDocument{ID: "doc2", Indices: []int{2, 3}, Values: []float64{1.0, 0.0}, Norm: 1.0}
	_ = store.Create(ctx, "doc1", doc1)
	_ = store.Create(ctx, "doc2", doc2)

	query := TestDocument{Indices: []int{0, 1}, Values: []float64{1.0, 0.0}, Norm: 1.0}
	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarCosine(store, query, opts, nil)

	// Assert - doc1 should score 1.0, doc2 should score 0.0
	assert.That(t, "results len must be 2", len(results), 2)
	assert.That(t, "first result score must be 1.0", results[0].Score, 1.0)
	assert.That(t, "second result score must be 0.0", results[1].Score, 0.0)
}

func Test_FindSimilarCosine_With_PartialOverlap_Should_ReturnCorrectScore(t *testing.T) {
	// Arrange: A=[1,0], B=[0.6,0.8] at indices [0,1]
	// dot = 1*0.6 + 0*0.8 = 0.6, normA=1, normB=1, cosine = 0.6
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	doc := TestDocument{Indices: []int{0, 1}, Values: []float64{0.6, 0.8}, Norm: 1.0}
	_ = store.Create(ctx, "doc1", doc)

	query := TestDocument{Indices: []int{0, 1}, Values: []float64{1.0, 0.0}, Norm: 1.0}
	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarCosine(store, query, opts, nil)

	// Assert
	assert.That(t, "results len must be 1", len(results), 1)
	assert.That(t, "score must be 0.6", math.Abs(results[0].Score-0.6) < 0.001, true)
}

func Test_FindSimilarCosine_With_StoppedFlag_Should_StopEarly(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	for i := range 100 {
		_ = store.Create(ctx, fmt.Sprintf("doc%d", i), TestDocument{
			Indices: []int{i % 10},
			Values:  []float64{1.0},
			Norm:    1.0,
		})
	}

	query := TestDocument{Indices: []int{0}, Values: []float64{1.0}, Norm: 1.0}
	opts := efficiency.DefaultSearchOptions()

	stopped := &atomic.Bool{}
	stopped.Store(true) // Pre-cancel

	// Act
	results := efficiency.FindSimilarCosine(store, query, opts, stopped)

	// Assert - should return partial or no results
	assert.That(t, "results should be fewer than total docs", len(results) < 100, true)
}

func Test_FindSimilarCosine_With_Threshold_Should_FilterLowScores(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	_ = store.Create(ctx, "high", TestDocument{Indices: []int{0, 1}, Values: []float64{0.8, 0.6}, Norm: 1.0})
	_ = store.Create(ctx, "low", TestDocument{Indices: []int{5, 6}, Values: []float64{0.8, 0.6}, Norm: 1.0})

	query := TestDocument{Indices: []int{0, 1}, Values: []float64{0.8, 0.6}, Norm: 1.0}
	opts := efficiency.SearchOptions{TopK: 10, Threshold: 0.5}

	// Act
	results := efficiency.FindSimilarCosine(store, query, opts, nil)

	// Assert
	assert.That(t, "only high-scoring result should be returned", len(results), 1)
	assert.That(t, "result key must be high", results[0].Key, "high")
}

func Test_FindSimilarCosine_With_TopK_Should_LimitResults(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	for i := range 20 {
		val := float64(i+1) / 20.0
		_ = store.Create(ctx, fmt.Sprintf("doc%d", i), TestDocument{
			Indices: []int{0},
			Values:  []float64{val},
			Norm:    val,
		})
	}

	query := TestDocument{Indices: []int{0}, Values: []float64{1.0}, Norm: 1.0}
	opts := efficiency.SearchOptions{TopK: 5, Threshold: 0.0}

	// Act
	results := efficiency.FindSimilarCosine(store, query, opts, nil)

	// Assert
	assert.That(t, "results len must be 5", len(results), 5)
	// Results should be sorted descending
	for i := range len(results) - 1 {
		assert.That(t, "results must be descending", results[i].Score >= results[i+1].Score, true)
	}
}

// --- Jaccard Similarity Tests ---

func Test_FindSimilarJaccard_With_DisjointSets_Should_ReturnScoreZero(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	_ = store.Create(ctx, "doc1", TestDocument{Indices: []int{0, 1, 2}})

	query := TestDocument{Indices: []int{3, 4, 5}}
	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarJaccard(store, query, opts, nil)

	// Assert
	assert.That(t, "results len must be 1", len(results), 1)
	assert.That(t, "score must be 0.0", results[0].Score, 0.0)
}

func Test_FindSimilarJaccard_With_EmptyQuery_Should_ReturnNoResults(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	_ = store.Create(ctx, "doc1", TestDocument{Indices: []int{1, 2, 3}})

	query := TestDocument{Indices: []int{}}
	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarJaccard(store, query, opts, nil)

	// Assert
	assert.That(t, "results must be empty", len(results), 0)
}

func Test_FindSimilarJaccard_With_IdenticalSets_Should_ReturnScoreOne(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	doc := TestDocument{Indices: []int{1, 3, 5, 7}}
	_ = store.Create(ctx, "doc1", doc)

	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarJaccard(store, doc, opts, nil)

	// Assert
	assert.That(t, "results len must be 1", len(results), 1)
	assert.That(t, "score must be 1.0", results[0].Score, 1.0)
}

func Test_FindSimilarJaccard_With_PartialOverlap_Should_ReturnCorrectScore(t *testing.T) {
	// Arrange: A={1,2,3}, B={2,3,4} -> intersection=2, union=4 -> 0.5
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	_ = store.Create(ctx, "doc1", TestDocument{Indices: []int{2, 3, 4}})

	query := TestDocument{Indices: []int{1, 2, 3}}
	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarJaccard(store, query, opts, nil)

	// Assert
	assert.That(t, "results len must be 1", len(results), 1)
	assert.That(t, "score must be 0.5", results[0].Score, 0.5)
}

func Test_FindSimilarJaccard_With_StoppedFlag_Should_StopEarly(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	for i := range 100 {
		_ = store.Create(ctx, fmt.Sprintf("doc%d", i), TestDocument{
			Indices: []int{i, i + 1, i + 2},
		})
	}

	query := TestDocument{Indices: []int{0, 1, 2}}
	opts := efficiency.DefaultSearchOptions()

	stopped := &atomic.Bool{}
	stopped.Store(true) // Pre-cancel

	// Act
	results := efficiency.FindSimilarJaccard(store, query, opts, stopped)

	// Assert - should return partial or no results
	assert.That(t, "results should be fewer than total docs", len(results) < 100, true)
}

func Test_FindSimilarJaccard_With_Threshold_Should_FilterLowScores(t *testing.T) {
	// Arrange: high overlap (3/3=1.0), low overlap (1/5=0.2)
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	_ = store.Create(ctx, "high", TestDocument{Indices: []int{1, 2, 3}})
	_ = store.Create(ctx, "low", TestDocument{Indices: []int{1, 10, 11}})

	query := TestDocument{Indices: []int{1, 2, 3}}
	opts := efficiency.SearchOptions{TopK: 10, Threshold: 0.5}

	// Act
	results := efficiency.FindSimilarJaccard(store, query, opts, nil)

	// Assert
	assert.That(t, "only high-scoring result should be returned", len(results), 1)
	assert.That(t, "result key must be high", results[0].Key, "high")
}

func Test_FindSimilarJaccard_With_TopK_Should_LimitResults(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	// Create documents with varying overlap with query {0,1,2,3,4}
	for i := range 20 {
		indices := make([]int, 5)
		for j := range 5 {
			indices[j] = i + j // Shift pattern
		}
		_ = store.Create(ctx, fmt.Sprintf("doc%d", i), TestDocument{Indices: indices})
	}

	query := TestDocument{Indices: []int{0, 1, 2, 3, 4}}
	opts := efficiency.SearchOptions{TopK: 5, Threshold: 0.0}

	// Act
	results := efficiency.FindSimilarJaccard(store, query, opts, nil)

	// Assert
	assert.That(t, "results len must be 5", len(results), 5)
	// Results should be sorted descending
	for i := range len(results) - 1 {
		assert.That(t, "results must be descending", results[i].Score >= results[i+1].Score, true)
	}
}

// --- SearchContext Tests ---

func Test_SearchContext_With_CancelledContext_Should_SetStoppedFlag(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())

	// Act
	stopped := efficiency.SearchContext(ctx)
	assert.That(t, "stopped must be false initially", stopped.Load(), false)

	cancel()
	// Give goroutine time to process
	for i := 0; i < 100 && !stopped.Load(); i++ {
		time.Sleep(time.Millisecond)
	}

	// Assert
	assert.That(t, "stopped must be true after cancel", stopped.Load(), true)
}

// --- Edge Cases ---

func Test_FindSimilarCosine_With_EmptyStore_Should_ReturnNoResults(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	query := TestDocument{Indices: []int{0}, Values: []float64{1.0}, Norm: 1.0}
	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarCosine(store, query, opts, nil)

	// Assert
	assert.That(t, "results must be empty", len(results), 0)
}

func Test_FindSimilarJaccard_With_EmptyStore_Should_ReturnNoResults(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	query := TestDocument{Indices: []int{0, 1, 2}}
	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarJaccard(store, query, opts, nil)

	// Assert
	assert.That(t, "results must be empty", len(results), 0)
}

func Test_FindSimilarCosine_With_ZeroNormCandidate_Should_SkipCandidate(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	_ = store.Create(ctx, "zero", TestDocument{Indices: []int{0}, Values: []float64{0.0}, Norm: 0.0})
	_ = store.Create(ctx, "valid", TestDocument{Indices: []int{0}, Values: []float64{1.0}, Norm: 1.0})

	query := TestDocument{Indices: []int{0}, Values: []float64{1.0}, Norm: 1.0}
	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarCosine(store, query, opts, nil)

	// Assert - zero-norm candidate should be skipped
	assert.That(t, "results len must be 1", len(results), 1)
	assert.That(t, "result key must be valid", results[0].Key, "valid")
}

func Test_FindSimilarJaccard_With_EmptyCandidate_Should_SkipCandidate(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	_ = store.Create(ctx, "empty", TestDocument{Indices: []int{}})
	_ = store.Create(ctx, "valid", TestDocument{Indices: []int{0, 1, 2}})

	query := TestDocument{Indices: []int{0, 1, 2}}
	opts := efficiency.DefaultSearchOptions()

	// Act
	results := efficiency.FindSimilarJaccard(store, query, opts, nil)

	// Assert - empty candidate should be skipped
	assert.That(t, "results len must be 1", len(results), 1)
	assert.That(t, "result key must be valid", results[0].Key, "valid")
}
