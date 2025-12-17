package efficiency_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func Test_SparseSet_With_AddAndGet_Should_ReturnCorrectValues(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSet[int](10)

	// Act
	s.Add(0, 1)
	s.Add(1, 2)
	s.Add(2, 3)

	// Assert
	assert.That(t, "element at index 0 must be 1", *s.Get(0), 1)
	assert.That(t, "element at index 1 must be 2", *s.Get(1), 2)
	assert.That(t, "element at index 2 must be 3", *s.Get(2), 3)
	assert.That(t, "element at index 3 must be nil", s.Get(3) == nil, true)
}

func Test_SparseSet_With_AddNegativeId_Should_NotAddElement(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSet[int](10)

	// Act
	s.Add(-1, 1)

	// Assert
	assert.That(t, "size should remain 0", s.Size, 0)
}

func Test_SparseSet_With_Empty_Should_HaveZeroSize(t *testing.T) {
	// Arrange & Act
	s := efficiency.NewSparseSet[int](10)

	// Assert
	assert.That(t, "size should be 0", s.Size, 0)
	assert.That(t, "densed should be empty", len(s.Densed()), 0)
}

func Test_SparseSet_With_GetNegativeId_Should_ReturnNil(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSet[int](10)
	s.Add(0, 1)

	// Act
	result := s.Get(-1)

	// Assert
	assert.That(t, "get with negative id should return nil", result == nil, true)
}

func Test_SparseSet_With_ModifyStruct_Should_UpdateValue(t *testing.T) {
	// Arrange
	type Person struct {
		name string
	}
	s := efficiency.NewSparseSet[Person](1)
	p := Person{name: "John"}
	s.Add(0, p)

	// Act
	p2 := s.Get(0)
	p2.name = "Jane"

	// Assert
	assert.That(t, "element at index 0 must be Jane", p2.name, "Jane")
}

func Test_SparseSet_With_MultipleAdds_Should_ReturnCorrectDense(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSet[int](10)

	// Act
	s.Add(0, 1)
	s.Add(1, 2)
	s.Add(2, 3)

	// Assert
	assert.That(t, "size must be 3", s.Size, 3)
	assert.That(t, "dense must be [1, 2, 3]", s.Densed(), []int{1, 2, 3})
}

func Test_SparseSet_With_RemoveElements_Should_UpdateDense(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSet[int](10)
	s.Add(0, 1)
	s.Add(1, 2)
	s.Add(2, 3)
	s.Add(3, 4)
	s.Add(4, 5)

	// Act
	s.Remove(1)
	s.Remove(3)

	// Assert
	assert.That(t, "size must be 3", s.Size, 3)
	assert.That(t, "dense must be [1, 3]", s.Densed(), []int{1, 5, 3})
}

func Test_SparseSet_With_RemoveNegativeId_Should_NotChangeSize(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSet[int](10)
	s.Add(0, 1)

	// Act
	s.Remove(-1)

	// Assert
	assert.That(t, "size should remain 1", s.Size, 1)
}

func Test_SparseSet_With_RemoveOutOfBounds_Should_NotChangeSize(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSet[int](10)
	s.Add(0, 1)

	// Act
	s.Remove(100)

	// Assert
	assert.That(t, "size should remain 1", s.Size, 1)
}

func Test_SparseSet_With_StringType_Should_StoreCorrectly(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSet[string](5)

	// Act
	s.Add(0, "hello")
	s.Add(1, "world")

	// Assert
	assert.That(t, "size must be 2", s.Size, 2)
	assert.That(t, "element at 0 must be hello", *s.Get(0), "hello")
	assert.That(t, "element at 1 must be world", *s.Get(1), "world")
}
