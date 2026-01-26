package efficiency_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func Test_KeyedSparseSet_With_Clear_Should_RemoveAllElements(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)
	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	// Act
	s.Clear()

	// Assert
	assert.That(t, "len must be 0", s.Len(), 0)
	assert.That(t, "get a must be nil", s.Get("a") == nil, true)
}

func Test_KeyedSparseSet_With_DeleteFirstElement_Should_MaintainOtherElements(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)
	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	// Act
	deleted := s.Delete("a")

	// Assert
	assert.That(t, "delete must return true", deleted, true)
	assert.That(t, "len must be 2", s.Len(), 2)
	assert.That(t, "get a must be nil", s.Get("a") == nil, true)
	assert.That(t, "get b must be 2", *s.Get("b"), 2)
	assert.That(t, "get c must be 3", *s.Get("c"), 3)
}

func Test_KeyedSparseSet_With_DeleteLastElement_Should_ShrinkCorrectly(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)
	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	// Act
	deleted := s.Delete("c")

	// Assert
	assert.That(t, "delete must return true", deleted, true)
	assert.That(t, "len must be 2", s.Len(), 2)
	assert.That(t, "get a must be 1", *s.Get("a"), 1)
	assert.That(t, "get b must be 2", *s.Get("b"), 2)
	assert.That(t, "get c must be nil", s.Get("c") == nil, true)
}

func Test_KeyedSparseSet_With_DeleteMiddleElement_Should_SwapRemove(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)
	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	// Act
	deleted := s.Delete("b")

	// Assert
	assert.That(t, "delete must return true", deleted, true)
	assert.That(t, "len must be 2", s.Len(), 2)
	assert.That(t, "get a must be 1", *s.Get("a"), 1)
	assert.That(t, "get b must be nil", s.Get("b") == nil, true)
	assert.That(t, "get c must be 3", *s.Get("c"), 3)
}

func Test_KeyedSparseSet_With_DeleteMissingKey_Should_ReturnFalse(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)
	s.Put("a", 1)

	// Act
	deleted := s.Delete("missing")

	// Assert
	assert.That(t, "delete must return false", deleted, false)
	assert.That(t, "len must be 1", s.Len(), 1)
}

func Test_KeyedSparseSet_With_ForEach_Should_VisitAllElements(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)
	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	// Act
	visited := make(map[string]int)
	s.ForEach(func(k string, v int) bool {
		visited[k] = v
		return true
	})

	// Assert
	assert.That(t, "visited len must be 3", len(visited), 3)
	assert.That(t, "visited a must be 1", visited["a"], 1)
	assert.That(t, "visited b must be 2", visited["b"], 2)
	assert.That(t, "visited c must be 3", visited["c"], 3)
}

func Test_KeyedSparseSet_With_ForEachEarlyStop_Should_StopIteration(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)
	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	// Act
	count := 0
	s.ForEach(func(k string, v int) bool {
		count++
		return count < 2
	})

	// Assert
	assert.That(t, "count must be 2", count, 2)
}

func Test_KeyedSparseSet_With_GetMissingKey_Should_ReturnNil(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)

	// Act
	value := s.Get("missing")

	// Assert
	assert.That(t, "value must be nil", value == nil, true)
}

func Test_KeyedSparseSet_With_Has_Should_ReturnCorrectResult(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)
	s.Put("a", 1)

	// Act & Assert
	assert.That(t, "has a must be true", s.Has("a"), true)
	assert.That(t, "has b must be false", s.Has("b"), false)
}

func Test_KeyedSparseSet_With_PutDuplicateKey_Should_UpdateValue(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)
	s.Put("a", 1)

	// Act
	isNew := s.Put("a", 42)

	// Assert
	assert.That(t, "isNew must be false", isNew, false)
	assert.That(t, "value must be 42", *s.Get("a"), 42)
	assert.That(t, "len must be 1", s.Len(), 1)
}

func Test_KeyedSparseSet_With_PutNewKey_Should_AddElement(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)

	// Act
	isNew := s.Put("a", 42)

	// Assert
	assert.That(t, "isNew must be true", isNew, true)
	assert.That(t, "value must be 42", *s.Get("a"), 42)
	assert.That(t, "len must be 1", s.Len(), 1)
}

func Test_KeyedSparseSet_With_Values_Should_ReturnDenseSlice(t *testing.T) {
	// Arrange
	s := efficiency.NewKeyedSparseSet[string, int](10)
	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	// Act
	values := s.Values()

	// Assert
	assert.That(t, "values len must be 3", len(values), 3)
}

func Test_KeyedSparseSet_With_ZeroCapacity_Should_UseDefault(t *testing.T) {
	// Arrange & Act
	s := efficiency.NewKeyedSparseSet[string, int](0)
	s.Put("a", 1)

	// Assert
	assert.That(t, "value must be 1", *s.Get("a"), 1)
}
