package efficiency_test

import (
	"fmt"

	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func ExampleGenerate() {
	for value := range efficiency.Generate(1, 2, 3) {
		fmt.Println(value)
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleJaccardSimilarity() {
	// Both sets must be sorted ascending: the comparison is a merge walk.
	score := efficiency.JaccardSimilarity([]int{1, 2, 3}, []int{2, 3, 4})

	// Two terms shared out of four distinct terms.
	fmt.Println(score)
	// Output: 0.5
}

func ExampleCosineSimilarity() {
	// A sparse vector is its sorted indices, its values, and its precomputed
	// norm. Caching the norm is what keeps the comparison cheap.
	indices := []int{0, 1}
	values := []float64{3, 4}
	norm := 5.0

	fmt.Println(efficiency.CosineSimilarity(indices, indices, values, values, norm, norm))
	// Output: 1
}

func ExampleNewKeyedSparseSet() {
	set := efficiency.NewKeyedSparseSet[string, int](8)

	fmt.Println(set.Put("a", 1))
	fmt.Println(set.Put("a", 2))
	fmt.Println(*set.Get("a"))
	fmt.Println(set.Delete("a"), set.Len())
	// Output:
	// true
	// false
	// 2
	// true 0
}
