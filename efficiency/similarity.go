package efficiency

// CosineSimilarity computes cosine similarity between two sparse vectors using merge-loop.
// Both index slices MUST be sorted in ascending order.
// Formula: cosine(A, B) = dot(A, B) / (norm(A) * norm(B)).
//
// Example usage with ShardedSparseAccess.SearchSimilar:
//
//	results := store.SearchSimilar(ctx, func(doc Document) float64 {
//	    return efficiency.CosineSimilarity(
//	        query.Indices, doc.Indices,
//	        query.Values, doc.Values,
//	        query.Norm, doc.Norm,
//	    )
//	}, resource.SearchOptions{TopK: 10})
func CosineSimilarity(indicesA, indicesB []int, valuesA, valuesB []float64, normA, normB float64) float64 {
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

// JaccardSimilarity computes Jaccard similarity between two sparse sets using merge-loop.
// Both index slices MUST be sorted in ascending order.
// Formula: jaccard(A, B) = |A ∩ B| / |A ∪ B|.
//
// Example usage with ShardedSparseAccess.SearchSimilar:
//
//	results := store.SearchSimilar(ctx, func(article Article) float64 {
//	    return efficiency.JaccardSimilarity(queryTags, article.Tags)
//	}, resource.SearchOptions{TopK: 10})
func JaccardSimilarity(setA, setB []int) float64 {
	if len(setA) == 0 && len(setB) == 0 {
		return 1.0 // Both empty sets are identical.
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

	// Count remaining elements.
	union += (len(setA) - i) + (len(setB) - j)

	return float64(intersection) / float64(union)
}
