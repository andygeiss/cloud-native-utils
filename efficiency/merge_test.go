package efficiency_test

import (
	"cloud-native/efficiency"
	"testing"
)

func TestMerge_One_Producer(t *testing.T) {
	in := []int{1, 2, 3}
	ch := efficiency.Generate[int](in...)
	producer := efficiency.Split(ch, 1)
	consumer := efficiency.Merge(producer...)
	sum := 0
	for val := range consumer {
		sum += val
	}
	if sum != 6 {
		t.Fatalf("sum must be 6, but got %d", sum)
	}
}

func TestMerge_Three_Producer(t *testing.T) {
	in := []int{1, 2, 3}
	ch := efficiency.Generate[int](in...)
	producer := efficiency.Split(ch, 3)
	consumer := efficiency.Merge(producer...)
	sum := 0
	for val := range consumer {
		sum += val
	}
	if sum != 6 {
		t.Fatalf("sum must be 6, but got %d", sum)
	}
}
