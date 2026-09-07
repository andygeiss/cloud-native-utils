package slices_test

import (
	"fmt"

	"github.com/andygeiss/cloud-native-utils/slices"
)

func ExampleMap() {
	lengths := slices.Map([]string{"go", "rust", "zig"}, func(s string) int { return len(s) })
	fmt.Println(lengths)
	// Output: [2 4 3]
}

func ExampleFilter() {
	even := slices.Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
	fmt.Println(even)
	// Output: [2 4 6]
}

func ExampleUnique() {
	fmt.Println(slices.Unique([]string{"a", "b", "a", "c", "b"}))
	// Output: [a b c]
}

func ExampleContains() {
	fmt.Println(slices.Contains([]int{1, 2, 3}, 2))
	fmt.Println(slices.Contains([]int{1, 2, 3}, 9))
	// Output:
	// true
	// false
}

func ExampleFirst() {
	value, ok := slices.First([]string{"alpha", "beta"})
	fmt.Println(value, ok)

	// An empty slice reports false rather than panicking.
	_, ok = slices.First([]string{})
	fmt.Println(ok)
	// Output:
	// alpha true
	// false
}
