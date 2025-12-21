package slices

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
)

func Test_Contains_With_ExistingElement_Should_ReturnTrue(t *testing.T) {
	// Arrange, Act & Assert
	slice := []string{"a", "b", "c"}
	assert.That(t, "should contain 'a'", Contains(slice, "a"), true)
	assert.That(t, "should contain 'b'", Contains(slice, "b"), true)
	assert.That(t, "should contain 'c'", Contains(slice, "c"), true)
}

func Test_Contains_With_MissingElement_Should_ReturnFalse(t *testing.T) {
	// Arrange, Act & Assert
	slice := []string{"a", "b", "c"}
	assert.That(t, "should not contain 'd'", Contains(slice, "d"), false)
}

func Test_Contains_With_EmptySlice_Should_ReturnFalse(t *testing.T) {
	// Arrange, Act & Assert
	slice := []string{}
	assert.That(t, "empty slice should not contain anything", Contains(slice, "a"), false)
}

func Test_Contains_With_Ints_Should_Work(t *testing.T) {
	// Arrange, Act & Assert
	slice := []int{1, 2, 3, 4, 5}
	assert.That(t, "should contain 3", Contains(slice, 3), true)
	assert.That(t, "should not contain 6", Contains(slice, 6), false)
}

func Test_ContainsAny_With_MatchingElements_Should_ReturnTrue(t *testing.T) {
	// Arrange, Act & Assert
	slice := []string{"a", "b", "c"}
	assert.That(t, "should find 'b' or 'd'", ContainsAny(slice, []string{"b", "d"}), true)
	assert.That(t, "should find 'c'", ContainsAny(slice, []string{"c"}), true)
}

func Test_ContainsAny_With_NoMatchingElements_Should_ReturnFalse(t *testing.T) {
	// Arrange, Act & Assert
	slice := []string{"a", "b", "c"}
	assert.That(t, "should not find 'd' or 'e'", ContainsAny(slice, []string{"d", "e"}), false)
}

func Test_ContainsAll_With_AllElements_Should_ReturnTrue(t *testing.T) {
	// Arrange, Act & Assert
	slice := []string{"a", "b", "c", "d"}
	assert.That(t, "should contain all of a,b,c", ContainsAll(slice, []string{"a", "b", "c"}), true)
}

func Test_ContainsAll_With_MissingElement_Should_ReturnFalse(t *testing.T) {
	// Arrange, Act & Assert
	slice := []string{"a", "b", "c"}
	assert.That(t, "should not contain all of a,b,d", ContainsAll(slice, []string{"a", "b", "d"}), false)
}

func Test_IndexOf_With_ExistingElement_Should_ReturnIndex(t *testing.T) {
	// Arrange, Act & Assert
	slice := []string{"a", "b", "c"}
	assert.That(t, "index of 'a' should be 0", IndexOf(slice, "a"), 0)
	assert.That(t, "index of 'b' should be 1", IndexOf(slice, "b"), 1)
	assert.That(t, "index of 'c' should be 2", IndexOf(slice, "c"), 2)
}

func Test_IndexOf_With_MissingElement_Should_ReturnNegativeOne(t *testing.T) {
	// Arrange, Act & Assert
	slice := []string{"a", "b", "c"}
	assert.That(t, "index of 'd' should be -1", IndexOf(slice, "d"), -1)
}

func Test_Filter_With_Predicate_Should_ReturnMatchingElements(t *testing.T) {
	// Arrange
	slice := []int{1, 2, 3, 4, 5, 6}
	// Act
	result := Filter(slice, func(n int) bool { return n%2 == 0 })
	// Assert
	assert.That(t, "should have 3 even numbers", len(result), 3)
	assert.That(t, "first even should be 2", result[0], 2)
	assert.That(t, "second even should be 4", result[1], 4)
	assert.That(t, "third even should be 6", result[2], 6)
}

func Test_Map_With_Mapper_Should_TransformElements(t *testing.T) {
	// Arrange
	slice := []int{1, 2, 3}
	// Act
	result := Map(slice, func(n int) int { return n * 2 })
	// Assert
	assert.That(t, "should have 3 elements", len(result), 3)
	assert.That(t, "first should be 2", result[0], 2)
	assert.That(t, "second should be 4", result[1], 4)
	assert.That(t, "third should be 6", result[2], 6)
}

func Test_Unique_With_Duplicates_Should_RemoveDuplicates(t *testing.T) {
	// Arrange
	slice := []string{"a", "b", "a", "c", "b", "d"}
	// Act
	result := Unique(slice)
	// Assert
	assert.That(t, "should have 4 unique elements", len(result), 4)
	assert.That(t, "first should be 'a'", result[0], "a")
}

func Test_First_With_NonEmptySlice_Should_ReturnFirst(t *testing.T) {
	// Arrange
	slice := []string{"a", "b", "c"}
	// Act
	first, ok := First(slice)
	// Assert
	assert.That(t, "ok should be true", ok, true)
	assert.That(t, "first should be 'a'", first, "a")
}

func Test_First_With_EmptySlice_Should_ReturnFalse(t *testing.T) {
	// Arrange
	slice := []string{}
	// Act
	_, ok := First(slice)
	// Assert
	assert.That(t, "ok should be false", ok, false)
}

func Test_Last_With_NonEmptySlice_Should_ReturnLast(t *testing.T) {
	// Arrange
	slice := []string{"a", "b", "c"}
	// Act
	last, ok := Last(slice)
	// Assert
	assert.That(t, "ok should be true", ok, true)
	assert.That(t, "last should be 'c'", last, "c")
}

func Test_Last_With_EmptySlice_Should_ReturnFalse(t *testing.T) {
	// Arrange
	slice := []string{}
	// Act
	_, ok := Last(slice)
	// Assert
	assert.That(t, "ok should be false", ok, false)
}

func Test_Copy_With_Slice_Should_CreateIndependentCopy(t *testing.T) {
	// Arrange
	original := []int{1, 2, 3}
	// Act
	copied := Copy(original)
	copied[0] = 100
	// Assert
	assert.That(t, "original should be unchanged", original[0], 1)
	assert.That(t, "copy should be modified", copied[0], 100)
}

func Test_Copy_With_NilSlice_Should_ReturnNil(t *testing.T) {
	// Arrange
	var original []int = nil
	// Act
	copied := Copy(original)
	// Assert
	assert.That(t, "copy of nil should be nil", copied == nil, true)
}
