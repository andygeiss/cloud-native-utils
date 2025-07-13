package efficiency_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func TestSparseSet_Add(t *testing.T) {
	s := efficiency.NewSparseSet[int](10)
	s.Add(0, 1)
	s.Add(1, 2)
	s.Add(2, 3)
	assert.That(t, "size must be 3", s.Size, 3)
	assert.That(t, "dense must be [1, 2, 3]", s.Densed(), []int{1, 2, 3})
}

func TestSparseSet_Remove(t *testing.T) {
	s := efficiency.NewSparseSet[int](10)
	s.Add(0, 1)
	s.Add(1, 2)
	s.Add(2, 3)
	s.Add(3, 4)
	s.Add(4, 5)
	s.Remove(1)
	s.Remove(3)
	assert.That(t, "size must be 3", s.Size, 3)
	assert.That(t, "dense must be [1, 3]", s.Densed(), []int{1, 5, 3})
}

func TestSparseSet_Serialize(t *testing.T) {
	s := efficiency.NewSparseSet[int](1)
	s.Add(0, 1)
	bytes := s.Serialize()
	assert.That(t, "bytes length must be 92", len(bytes), 92)
}
