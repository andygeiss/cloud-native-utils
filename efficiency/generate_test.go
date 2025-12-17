package efficiency_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func Test_Generate_With_EmptyInput_Should_ReturnEmptyChannel(t *testing.T) {
	// Arrange & Act
	ch := efficiency.Generate[int]()
	out := make([]int, 0)
	for val := range ch {
		out = append(out, val)
	}

	// Assert
	assert.That(t, "output should be empty", len(out), 0)
}

func Test_Generate_With_MultipleValues_Should_CloseChannel(t *testing.T) {
	// Arrange & Act
	ch := efficiency.Generate[int](1, 2, 3)
	for range ch {
	}
	_, ok := <-ch

	// Assert
	assert.That(t, "channel should be closed", ok, false)
}

func Test_Generate_With_SingleElement_Should_ReturnOneElement(t *testing.T) {
	// Arrange & Act
	ch := efficiency.Generate[int](42)
	out := make([]int, 0)
	for val := range ch {
		out = append(out, val)
	}

	// Assert
	assert.That(t, "output should have one element", len(out), 1)
	assert.That(t, "element should be 42", out[0], 42)
}

func Test_Generate_With_StringSlice_Should_ReturnAllStrings(t *testing.T) {
	// Arrange
	in := []string{"a", "b", "c"}

	// Act
	ch := efficiency.Generate[string](in...)
	out := make([]string, 0)
	for val := range ch {
		out = append(out, val)
	}

	// Assert
	assert.That(t, "string slices must be equal", in, out)
}

func Test_Generate_With_ValidIntSlice_Should_ReturnAllInts(t *testing.T) {
	// Arrange
	in := []int{1, 2, 3}

	// Act
	ch := efficiency.Generate[int](in...)
	out := make([]int, 0)
	for val := range ch {
		out = append(out, val)
	}

	// Assert
	assert.That(t, "in and out slice must be equal", in, out)
}
