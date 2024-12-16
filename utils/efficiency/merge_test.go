package efficiency_test

import (
	"cloud-native/utils/assert"
	"cloud-native/utils/efficiency"
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
	assert.That(t, "sum must be correct", sum, 6)
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
	assert.That(t, "sum must be correct", sum, 6)
}
