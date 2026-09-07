package efficiency_test

import (
	"math"
	"sort"
	"testing"

	"github.com/andygeiss/cloud-native-utils/efficiency"
)

// sortedSet turns fuzz bytes into the ascending index slice the similarity
// functions document as their input.
func sortedSet(raw []byte) []int {
	set := make([]int, 0, len(raw))
	for _, b := range raw {
		set = append(set, int(b))
	}
	sort.Ints(set)
	return set
}

func FuzzJaccardSimilarity(f *testing.F) {
	f.Add([]byte{1, 2, 3}, []byte{2, 3, 4})
	f.Add([]byte{}, []byte{})
	f.Add([]byte{1}, []byte{})
	f.Add([]byte{7, 7, 7}, []byte{7}) // repeats, which a set should absorb

	f.Fuzz(func(t *testing.T, rawA, rawB []byte) {
		score := efficiency.JaccardSimilarity(sortedSet(rawA), sortedSet(rawB))

		if math.IsNaN(score) || score < 0 || score > 1 {
			t.Errorf("similarity must be a fraction, got %v", score)
		}
	})
}

func FuzzCosineSimilarity(f *testing.F) {
	f.Add([]byte{1, 2, 3}, []byte{2, 3, 4})
	f.Add([]byte{}, []byte{})
	f.Add([]byte{0}, []byte{0})

	f.Fuzz(func(t *testing.T, rawA, rawB []byte) {
		indicesA, indicesB := sortedSet(rawA), sortedSet(rawB)

		// Values run parallel to the indices, which the function requires.
		valuesA := make([]float64, len(indicesA))
		for i := range valuesA {
			valuesA[i] = 1
		}
		valuesB := make([]float64, len(indicesB))
		for i := range valuesB {
			valuesB[i] = 1
		}
		normA := math.Sqrt(float64(len(indicesA)))
		normB := math.Sqrt(float64(len(indicesB)))

		score := efficiency.CosineSimilarity(indicesA, indicesB, valuesA, valuesB, normA, normB)

		if math.IsNaN(score) || math.IsInf(score, 0) {
			t.Errorf("similarity must be finite, got %v", score)
		}
	})
}
