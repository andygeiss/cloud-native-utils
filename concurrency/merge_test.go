package concurrency_test

import (
	"cloud-native/concurrency"
	"testing"
)

func TestMerge_One_Producer(t *testing.T) {
	in := []int{1, 2, 3}
	ch := concurrency.Generate[int](in...)
	producer := concurrency.Split(ch, 1)
	consumer := concurrency.Merge(producer...)
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
	ch := concurrency.Generate[int](in...)
	producer := concurrency.Split(ch, 3)
	consumer := concurrency.Merge(producer...)
	sum := 0
	for val := range consumer {
		sum += val
	}
	if sum != 6 {
		t.Fatalf("sum must be 6, but got %d", sum)
	}
}
