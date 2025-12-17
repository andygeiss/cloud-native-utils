package assert_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_That_With_DifferentBools_Should_Fail(t *testing.T) {
	// Arrange
	mockT := &testing.T{}

	// Act
	assert.That(mockT, "bools should be equal", true, false)

	// Assert
	if !mockT.Failed() {
		t.Error("expected test to fail for different bools")
	}
}

func Test_That_With_DifferentFloats_Should_Fail(t *testing.T) {
	// Arrange
	mockT := &testing.T{}

	// Act
	assert.That(mockT, "floats should be equal", 1.5, 2.5)

	// Assert
	if !mockT.Failed() {
		t.Error("expected test to fail for different floats")
	}
}

func Test_That_With_DifferentIntegers_Should_Fail(t *testing.T) {
	// Arrange
	mockT := &testing.T{}

	// Act
	assert.That(mockT, "integers should be equal", 1, 2)

	// Assert
	if !mockT.Failed() {
		t.Error("expected test to fail")
	}
}

func Test_That_With_DifferentMaps_Should_Fail(t *testing.T) {
	// Arrange
	mockT := &testing.T{}
	m1 := map[string]int{"a": 1}
	m2 := map[string]int{"a": 2}

	// Act
	assert.That(mockT, "maps should be equal", m1, m2)

	// Assert
	if !mockT.Failed() {
		t.Error("expected test to fail for different maps")
	}
}

func Test_That_With_DifferentNestedSlices_Should_Fail(t *testing.T) {
	// Arrange
	mockT := &testing.T{}
	s1 := [][]int{{1, 2}, {3, 4}}
	s2 := [][]int{{1, 2}, {3, 5}}

	// Act
	assert.That(mockT, "nested slices should be equal", s1, s2)

	// Assert
	if !mockT.Failed() {
		t.Error("expected test to fail for different nested slices")
	}
}

func Test_That_With_DifferentSlices_Should_Fail(t *testing.T) {
	// Arrange
	mockT := &testing.T{}
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 4}

	// Act
	assert.That(mockT, "slices should be equal", s1, s2)

	// Assert
	if !mockT.Failed() {
		t.Error("expected test to fail for different slices")
	}
}

func Test_That_With_DifferentStrings_Should_Fail(t *testing.T) {
	// Arrange
	mockT := &testing.T{}

	// Act
	assert.That(mockT, "strings should be equal", "hello", "world")

	// Assert
	if !mockT.Failed() {
		t.Error("expected test to fail for different strings")
	}
}

func Test_That_With_DifferentStructs_Should_Fail(t *testing.T) {
	// Arrange
	mockT := &testing.T{}
	type Person struct {
		Name string
		Age  int
	}
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 30}

	// Act
	assert.That(mockT, "structs should be equal", p1, p2)

	// Assert
	if !mockT.Failed() {
		t.Error("expected test to fail for different structs")
	}
}

func Test_That_With_EqualBools_Should_Pass(t *testing.T) {
	// Arrange
	mockT := &testing.T{}

	// Act
	assert.That(mockT, "bools should be equal", true, true)

	// Assert
	if mockT.Failed() {
		t.Error("expected test to pass for equal bools")
	}
}

func Test_That_With_EqualFloats_Should_Pass(t *testing.T) {
	// Arrange
	mockT := &testing.T{}

	// Act
	assert.That(mockT, "floats should be equal", 1.5, 1.5)

	// Assert
	if mockT.Failed() {
		t.Error("expected test to pass for equal floats")
	}
}

func Test_That_With_EqualIntegers_Should_Pass(t *testing.T) {
	// Arrange
	mockT := &testing.T{}

	// Act
	assert.That(mockT, "integers should be equal", 1, 1)

	// Assert
	if mockT.Failed() {
		t.Error("expected test to pass")
	}
}

func Test_That_With_EqualMaps_Should_Pass(t *testing.T) {
	// Arrange
	mockT := &testing.T{}
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"a": 1, "b": 2}

	// Act
	assert.That(mockT, "maps should be equal", m1, m2)

	// Assert
	if mockT.Failed() {
		t.Error("expected test to pass for equal maps")
	}
}

func Test_That_With_EqualNestedSlices_Should_Pass(t *testing.T) {
	// Arrange
	mockT := &testing.T{}
	s1 := [][]int{{1, 2}, {3, 4}}
	s2 := [][]int{{1, 2}, {3, 4}}

	// Act
	assert.That(mockT, "nested slices should be equal", s1, s2)

	// Assert
	if mockT.Failed() {
		t.Error("expected test to pass for equal nested slices")
	}
}

func Test_That_With_EqualSlices_Should_Pass(t *testing.T) {
	// Arrange
	mockT := &testing.T{}
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}

	// Act
	assert.That(mockT, "slices should be equal", s1, s2)

	// Assert
	if mockT.Failed() {
		t.Error("expected test to pass for equal slices")
	}
}

func Test_That_With_EqualStrings_Should_Pass(t *testing.T) {
	// Arrange
	mockT := &testing.T{}

	// Act
	assert.That(mockT, "strings should be equal", "hello", "hello")

	// Assert
	if mockT.Failed() {
		t.Error("expected test to pass for equal strings")
	}
}

func Test_That_With_EqualStructs_Should_Pass(t *testing.T) {
	// Arrange
	mockT := &testing.T{}
	type Person struct {
		Name string
		Age  int
	}
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Alice", Age: 30}

	// Act
	assert.That(mockT, "structs should be equal", p1, p2)

	// Assert
	if mockT.Failed() {
		t.Error("expected test to pass for equal structs")
	}
}

func Test_That_With_NilValues_Should_Pass(t *testing.T) {
	// Arrange
	mockT := &testing.T{}
	var p1 *int = nil
	var p2 *int = nil

	// Act
	assert.That(mockT, "nil values should be equal", p1, p2)

	// Assert
	if mockT.Failed() {
		t.Error("expected test to pass for nil values")
	}
}
