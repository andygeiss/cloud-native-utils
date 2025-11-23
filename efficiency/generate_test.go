package efficiency_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func TestGenerate_OK(t *testing.T) {
	in := []int{1, 2, 3}
	ch := efficiency.Generate[int](in...)
	out := make([]int, 0)
	for val := range ch {
		out = append(out, val)
	}
	assert.That(t, "in and out slice must be equal", in, out)
}
