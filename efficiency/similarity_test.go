package efficiency_test

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
	"github.com/andygeiss/cloud-native-utils/resource"
)

// TestDocument for integration tests with SearchSimilar.
type TestDocument struct {
	ID      string
	Indices []int
	Values  []float64
	Norm    float64
}

// computeNorm calculates L2 norm for test setup.
func computeNorm(values []float64) float64 {
	var sum float64
	for _, v := range values {
		sum += v * v
	}
	return math.Sqrt(sum)
}

// --- CosineSimilarity Unit Tests ---

func Test_CosineSimilarity_With_EmptyVectors_Should_ReturnZero(t *testing.T) {
	// Arrange
	indicesA := []int{}
	valuesA := []float64{}
	normA := 0.0
	indicesB := []int{0, 1}
	valuesB := []float64{1.0, 0.0}
	normB := 1.0

	// Act
	score := efficiency.CosineSimilarity(indicesA, indicesB, valuesA, valuesB, normA, normB)

	// Assert
	assert.That(t, "score must be 0.0", score, 0.0)
}

func Test_CosineSimilarity_With_IdenticalVectors_Should_ReturnOne(t *testing.T) {
	// Arrange
	indices := []int{0, 2}
	values := []float64{0.6, 0.8}
	norm := computeNorm(values)

	// Act
	score := efficiency.CosineSimilarity(indices, indices, values, values, norm, norm)

	// Assert
	assert.That(t, "score must be approximately 1.0", score > 0.999, true)
}

func Test_CosineSimilarity_With_OrthogonalVectors_Should_ReturnZero(t *testing.T) {
	// Arrange: vectors with no overlapping indices have dot product 0.
	indicesA := []int{0, 1}
	valuesA := []float64{1.0, 0.0}
	normA := 1.0
	indicesB := []int{2, 3}
	valuesB := []float64{1.0, 0.0}
	normB := 1.0

	// Act
	score := efficiency.CosineSimilarity(indicesA, indicesB, valuesA, valuesB, normA, normB)

	// Assert
	assert.That(t, "score must be 0.0", score, 0.0)
}

func Test_CosineSimilarity_With_PartialOverlap_Should_ReturnCorrectScore(t *testing.T) {
	// Arrange: A=[1,0], B=[0.6,0.8] at indices [0,1].
	// dot = 1*0.6 + 0*0.8 = 0.6, normA=1, normB=1, cosine = 0.6.
	indicesA := []int{0, 1}
	valuesA := []float64{1.0, 0.0}
	normA := 1.0
	indicesB := []int{0, 1}
	valuesB := []float64{0.6, 0.8}
	normB := 1.0

	// Act
	score := efficiency.CosineSimilarity(indicesA, indicesB, valuesA, valuesB, normA, normB)

	// Assert
	assert.That(t, "score must be 0.6", math.Abs(score-0.6) < 0.001, true)
}

func Test_CosineSimilarity_With_ZeroNorm_Should_ReturnZero(t *testing.T) {
	// Arrange
	indices := []int{0}
	values := []float64{0.0}

	// Act
	score := efficiency.CosineSimilarity(indices, indices, values, values, 0.0, 1.0)

	// Assert
	assert.That(t, "score must be 0.0", score, 0.0)
}

// --- JaccardSimilarity Unit Tests ---

func Test_JaccardSimilarity_With_BothEmpty_Should_ReturnOne(t *testing.T) {
	// Arrange
	setA := []int{}
	setB := []int{}

	// Act
	score := efficiency.JaccardSimilarity(setA, setB)

	// Assert
	assert.That(t, "score must be 1.0", score, 1.0)
}

func Test_JaccardSimilarity_With_DisjointSets_Should_ReturnZero(t *testing.T) {
	// Arrange
	setA := []int{0, 1, 2}
	setB := []int{3, 4, 5}

	// Act
	score := efficiency.JaccardSimilarity(setA, setB)

	// Assert
	assert.That(t, "score must be 0.0", score, 0.0)
}

func Test_JaccardSimilarity_With_IdenticalSets_Should_ReturnOne(t *testing.T) {
	// Arrange
	set := []int{1, 3, 5, 7}

	// Act
	score := efficiency.JaccardSimilarity(set, set)

	// Assert
	assert.That(t, "score must be 1.0", score, 1.0)
}

func Test_JaccardSimilarity_With_OneEmpty_Should_ReturnZero(t *testing.T) {
	// Arrange
	setA := []int{0, 1, 2}
	setB := []int{}

	// Act
	score := efficiency.JaccardSimilarity(setA, setB)

	// Assert
	assert.That(t, "score must be 0.0", score, 0.0)
}

func Test_JaccardSimilarity_With_PartialOverlap_Should_ReturnCorrectScore(t *testing.T) {
	// Arrange: A={1,2,3}, B={2,3,4} -> intersection=2, union=4 -> 0.5.
	setA := []int{1, 2, 3}
	setB := []int{2, 3, 4}

	// Act
	score := efficiency.JaccardSimilarity(setA, setB)

	// Assert
	assert.That(t, "score must be 0.5", score, 0.5)
}

// --- Integration Tests with SearchSimilar ---

func Test_SearchSimilar_With_CosineSimilarity_Should_FindSimilarDocs(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()

	values := []float64{0.6, 0.8}
	norm := computeNorm(values)
	doc := TestDocument{ID: "doc1", Indices: []int{0, 2}, Values: values, Norm: norm}
	_ = store.Create(ctx, "doc1", doc)

	query := doc

	// Act
	results := store.SearchSimilar(ctx, func(d TestDocument) float64 {
		return efficiency.CosineSimilarity(
			query.Indices, d.Indices,
			query.Values, d.Values,
			query.Norm, d.Norm,
		)
	}, resource.SearchOptions{TopK: 10})

	// Assert
	assert.That(t, "results len must be 1", len(results), 1)
	assert.That(t, "score must be approximately 1.0", results[0].Score > 0.999, true)
}

func Test_SearchSimilar_With_CosineSimilarity_Threshold_Should_Filter(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	_ = store.Create(ctx, "high", TestDocument{Indices: []int{0, 1}, Values: []float64{0.8, 0.6}, Norm: 1.0})
	_ = store.Create(ctx, "low", TestDocument{Indices: []int{5, 6}, Values: []float64{0.8, 0.6}, Norm: 1.0})

	query := TestDocument{Indices: []int{0, 1}, Values: []float64{0.8, 0.6}, Norm: 1.0}

	// Act
	results := store.SearchSimilar(ctx, func(d TestDocument) float64 {
		return efficiency.CosineSimilarity(
			query.Indices, d.Indices,
			query.Values, d.Values,
			query.Norm, d.Norm,
		)
	}, resource.SearchOptions{TopK: 10, Threshold: 0.5})

	// Assert
	assert.That(t, "only high-scoring result should be returned", len(results), 1)
	assert.That(t, "result key must be high", results[0].Key, "high")
}

func Test_SearchSimilar_With_CosineSimilarity_TopK_Should_Limit(t *testing.T) {
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

	// Act
	results := store.SearchSimilar(ctx, func(d TestDocument) float64 {
		return efficiency.CosineSimilarity(
			query.Indices, d.Indices,
			query.Values, d.Values,
			query.Norm, d.Norm,
		)
	}, resource.SearchOptions{TopK: 5})

	// Assert
	assert.That(t, "results len must be 5", len(results), 5)
	for i := range len(results) - 1 {
		assert.That(t, "results must be descending", results[i].Score >= results[i+1].Score, true)
	}
}

func Test_SearchSimilar_With_JaccardSimilarity_Should_FindSimilarDocs(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	doc := TestDocument{Indices: []int{1, 3, 5, 7}}
	_ = store.Create(ctx, "doc1", doc)

	query := doc

	// Act
	results := store.SearchSimilar(ctx, func(d TestDocument) float64 {
		return efficiency.JaccardSimilarity(query.Indices, d.Indices)
	}, resource.SearchOptions{TopK: 10})

	// Assert
	assert.That(t, "results len must be 1", len(results), 1)
	assert.That(t, "score must be 1.0", results[0].Score, 1.0)
}

func Test_SearchSimilar_With_JaccardSimilarity_PartialOverlap_Should_ReturnCorrectScore(t *testing.T) {
	// Arrange: A={1,2,3}, B={2,3,4} -> intersection=2, union=4 -> 0.5.
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	_ = store.Create(ctx, "doc1", TestDocument{Indices: []int{2, 3, 4}})

	query := TestDocument{Indices: []int{1, 2, 3}}

	// Act
	results := store.SearchSimilar(ctx, func(d TestDocument) float64 {
		return efficiency.JaccardSimilarity(query.Indices, d.Indices)
	}, resource.SearchOptions{TopK: 10})

	// Assert
	assert.That(t, "results len must be 1", len(results), 1)
	assert.That(t, "score must be 0.5", results[0].Score, 0.5)
}

func Test_SearchSimilar_With_EmptyStore_Should_ReturnNoResults(t *testing.T) {
	// Arrange
	store := resource.NewShardedSparseAccess[string, TestDocument](4)
	ctx := context.Background()
	query := TestDocument{Indices: []int{0}, Values: []float64{1.0}, Norm: 1.0}

	// Act
	results := store.SearchSimilar(ctx, func(d TestDocument) float64 {
		return efficiency.CosineSimilarity(
			query.Indices, d.Indices,
			query.Values, d.Values,
			query.Norm, d.Norm,
		)
	}, resource.SearchOptions{TopK: 10})

	// Assert
	assert.That(t, "results must be empty", len(results), 0)
}
