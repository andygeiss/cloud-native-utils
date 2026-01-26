package efficiency_test

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/andygeiss/cloud-native-utils/efficiency"
	"github.com/andygeiss/cloud-native-utils/resource"
)

// BenchDoc for benchmarking similarity search.
type BenchDoc struct {
	Indices []int
	Values  []float64
	Norm    float64
}

// generateZipfDocument creates a synthetic TF-IDF document with Zipf-like term distribution.
// vocabSize: total vocabulary size.
// avgTerms: average number of non-zero terms per document.
func generateZipfDocument(rng *rand.Rand, vocabSize, avgTerms int) BenchDoc {
	// Variance around avgTerms.
	numTerms := avgTerms/2 + rng.Intn(avgTerms)
	numTerms = min(numTerms, vocabSize)
	numTerms = max(numTerms, 1)

	// Select terms using Zipf-like distribution (more common terms more likely).
	termSet := make(map[int]struct{})
	for len(termSet) < numTerms {
		// Zipf: P(rank r) ~ 1/r - use power law.
		r := rng.Float64()
		term := int(r * r * float64(vocabSize))
		if term >= vocabSize {
			term = vocabSize - 1
		}
		termSet[term] = struct{}{}
	}

	// Convert to sorted indices.
	indices := make([]int, 0, len(termSet))
	for term := range termSet {
		indices = append(indices, term)
	}
	sort.Ints(indices)

	// Generate TF-IDF values (random for benchmark purposes).
	values := make([]float64, len(indices))
	var normSquared float64
	for i := range values {
		values[i] = rng.Float64()*0.5 + 0.1 // Range [0.1, 0.6].
		normSquared += values[i] * values[i]
	}
	norm := math.Sqrt(normSquared)

	return BenchDoc{
		Indices: indices,
		Values:  values,
		Norm:    norm,
	}
}

func createBenchStore(corpusSize, vocabSize, avgTermsPerDoc int) *resource.ShardedSparseAccess[string, BenchDoc] {
	store := resource.NewShardedSparseAccessWithCapacity[string, BenchDoc](32, corpusSize/32+1)
	ctx := context.Background()
	rng := rand.New(rand.NewSource(42)) //nolint:gosec // G404: weak random acceptable for benchmarks.

	for i := range corpusSize {
		doc := generateZipfDocument(rng, vocabSize, avgTermsPerDoc)
		_ = store.Create(ctx, fmt.Sprintf("doc-%d", i), doc)
	}
	return store
}

// newBenchRng creates a seeded random generator for benchmarks.
//
//nolint:gosec // G404: weak random acceptable for benchmarks.
func newBenchRng() *rand.Rand {
	return rand.New(rand.NewSource(99))
}

// --- Cosine Similarity Benchmarks ---

func BenchmarkSearchSimilar_Cosine_10k_Docs(b *testing.B) {
	store := createBenchStore(10000, 10000, 300)
	query := generateZipfDocument(newBenchRng(), 10000, 300)
	ctx := context.Background()
	opts := resource.SearchOptions{TopK: 10, Threshold: 0.0}

	b.ResetTimer()
	for b.Loop() {
		_ = store.SearchSimilar(ctx, func(d BenchDoc) float64 {
			return efficiency.CosineSimilarity(
				query.Indices, d.Indices,
				query.Values, d.Values,
				query.Norm, d.Norm,
			)
		}, opts)
	}
}

func BenchmarkSearchSimilar_Cosine_50k_Docs(b *testing.B) {
	store := createBenchStore(50000, 50000, 300)
	query := generateZipfDocument(newBenchRng(), 50000, 300)
	ctx := context.Background()
	opts := resource.SearchOptions{TopK: 10, Threshold: 0.0}

	b.ResetTimer()
	for b.Loop() {
		_ = store.SearchSimilar(ctx, func(d BenchDoc) float64 {
			return efficiency.CosineSimilarity(
				query.Indices, d.Indices,
				query.Values, d.Values,
				query.Norm, d.Norm,
			)
		}, opts)
	}
}

func BenchmarkSearchSimilar_Cosine_VaryingSparsity(b *testing.B) {
	for _, avgTerms := range []int{100, 300, 500} {
		b.Run(fmt.Sprintf("terms_%d", avgTerms), func(b *testing.B) {
			store := createBenchStore(10000, 10000, avgTerms)
			query := generateZipfDocument(newBenchRng(), 10000, avgTerms)
			ctx := context.Background()
			opts := resource.SearchOptions{TopK: 10, Threshold: 0.0}

			b.ResetTimer()
			for b.Loop() {
				_ = store.SearchSimilar(ctx, func(d BenchDoc) float64 {
					return efficiency.CosineSimilarity(
						query.Indices, d.Indices,
						query.Values, d.Values,
						query.Norm, d.Norm,
					)
				}, opts)
			}
		})
	}
}

// --- Jaccard Similarity Benchmarks ---

func BenchmarkSearchSimilar_Jaccard_10k_Docs(b *testing.B) {
	store := createBenchStore(10000, 10000, 300)
	query := generateZipfDocument(newBenchRng(), 10000, 300)
	ctx := context.Background()
	opts := resource.SearchOptions{TopK: 10, Threshold: 0.0}

	b.ResetTimer()
	for b.Loop() {
		_ = store.SearchSimilar(ctx, func(d BenchDoc) float64 {
			return efficiency.JaccardSimilarity(query.Indices, d.Indices)
		}, opts)
	}
}

func BenchmarkSearchSimilar_Jaccard_50k_Docs(b *testing.B) {
	store := createBenchStore(50000, 50000, 300)
	query := generateZipfDocument(newBenchRng(), 50000, 300)
	ctx := context.Background()
	opts := resource.SearchOptions{TopK: 10, Threshold: 0.0}

	b.ResetTimer()
	for b.Loop() {
		_ = store.SearchSimilar(ctx, func(d BenchDoc) float64 {
			return efficiency.JaccardSimilarity(query.Indices, d.Indices)
		}, opts)
	}
}

func BenchmarkSearchSimilar_Jaccard_VaryingSparsity(b *testing.B) {
	for _, avgTerms := range []int{100, 300, 500} {
		b.Run(fmt.Sprintf("terms_%d", avgTerms), func(b *testing.B) {
			store := createBenchStore(10000, 10000, avgTerms)
			query := generateZipfDocument(newBenchRng(), 10000, avgTerms)
			ctx := context.Background()
			opts := resource.SearchOptions{TopK: 10, Threshold: 0.0}

			b.ResetTimer()
			for b.Loop() {
				_ = store.SearchSimilar(ctx, func(d BenchDoc) float64 {
					return efficiency.JaccardSimilarity(query.Indices, d.Indices)
				}, opts)
			}
		})
	}
}

// --- Comparison: Jaccard vs Cosine ---

func BenchmarkComparison_Jaccard_vs_Cosine(b *testing.B) {
	store := createBenchStore(10000, 10000, 300)
	query := generateZipfDocument(newBenchRng(), 10000, 300)
	ctx := context.Background()
	opts := resource.SearchOptions{TopK: 10, Threshold: 0.0}

	b.Run("Cosine", func(b *testing.B) {
		for b.Loop() {
			_ = store.SearchSimilar(ctx, func(d BenchDoc) float64 {
				return efficiency.CosineSimilarity(
					query.Indices, d.Indices,
					query.Values, d.Values,
					query.Norm, d.Norm,
				)
			}, opts)
		}
	})

	b.Run("Jaccard", func(b *testing.B) {
		for b.Loop() {
			_ = store.SearchSimilar(ctx, func(d BenchDoc) float64 {
				return efficiency.JaccardSimilarity(query.Indices, d.Indices)
			}, opts)
		}
	})
}

// --- Top-K Variation ---

func BenchmarkSearchSimilar_Cosine_TopK_Variation(b *testing.B) {
	store := createBenchStore(10000, 10000, 300)
	query := generateZipfDocument(newBenchRng(), 10000, 300)
	ctx := context.Background()

	for _, k := range []int{1, 10, 50, 100} {
		b.Run(fmt.Sprintf("topK_%d", k), func(b *testing.B) {
			opts := resource.SearchOptions{TopK: k, Threshold: 0.0}
			b.ResetTimer()
			for b.Loop() {
				_ = store.SearchSimilar(ctx, func(d BenchDoc) float64 {
					return efficiency.CosineSimilarity(
						query.Indices, d.Indices,
						query.Values, d.Values,
						query.Norm, d.Norm,
					)
				}, opts)
			}
		})
	}
}

func BenchmarkSearchSimilar_Jaccard_TopK_Variation(b *testing.B) {
	store := createBenchStore(10000, 10000, 300)
	query := generateZipfDocument(newBenchRng(), 10000, 300)
	ctx := context.Background()

	for _, k := range []int{1, 10, 50, 100} {
		b.Run(fmt.Sprintf("topK_%d", k), func(b *testing.B) {
			opts := resource.SearchOptions{TopK: k, Threshold: 0.0}
			b.ResetTimer()
			for b.Loop() {
				_ = store.SearchSimilar(ctx, func(d BenchDoc) float64 {
					return efficiency.JaccardSimilarity(query.Indices, d.Indices)
				}, opts)
			}
		})
	}
}
