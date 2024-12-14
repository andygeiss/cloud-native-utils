package efficiency_test

import (
	"cloud-native/utils/efficiency"
	"reflect"
	"testing"
)

func TestGenerate_Empty(t *testing.T) {
	in := []int{}
	ch := efficiency.Generate[int](in...)
	out := make([]int, 0)
	for val := range ch {
		out = append(out, val)
	}
	if len(out) != 0 {
		t.Fatalf("out slice len must be 0, but got %d", len(out))
	}
}

func TestGenerate_Three_Int(t *testing.T) {
	in := []int{1, 2, 3}
	ch := efficiency.Generate[int](in...)
	out := make([]int, 0)
	for val := range ch {
		out = append(out, val)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatal("in and out slice must be equal")
	}
}