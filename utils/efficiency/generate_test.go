package efficiency_test

import (
	"cloud-native/utils/assert"
	"cloud-native/utils/efficiency"
	"testing"
)

func TestGenerate_Empty(t *testing.T) {
	in := []int{}
	ch := efficiency.Generate[int](in...)
	out := make([]int, 0)
	for val := range ch {
		out = append(out, val)
	}
	assert.That(t, "out slice len must be correct", len(out), 0)
}

func TestGenerate_Three_Int(t *testing.T) {
	in := []int{1, 2, 3}
	ch := efficiency.Generate[int](in...)
	out := make([]int, 0)
	for val := range ch {
		out = append(out, val)
	}
	assert.That(t, "in and out slice must be equal", in, out)
}
