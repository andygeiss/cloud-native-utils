package efficiency_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func Test_SparseSharding_With_Clear_Should_RemoveAllElements(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](4)
	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	// Act
	s.Clear()

	// Assert
	assert.That(t, "len must be 0", s.Len(), 0)
}

func Test_SparseSharding_With_ConcurrentAccess_Should_HandleSafely(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](32)
	numGoroutines := runtime.NumCPU() * 2
	numOpsPerGoroutine := 100

	// Act
	var wg sync.WaitGroup
	for i := range numGoroutines {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			for j := range numOpsPerGoroutine {
				key := fmt.Sprintf("key-%d-%d", gid, j)
				s.Put(key, gid*1000+j)
			}
		}(i)
	}
	wg.Wait()

	// Assert
	expectedLen := numGoroutines * numOpsPerGoroutine
	assert.That(t, "len must match expected", s.Len(), expectedLen)
}

func Test_SparseSharding_With_ConcurrentDelete_Should_HandleSafely(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](32)
	numKeys := 1000

	for i := range numKeys {
		s.Put(fmt.Sprintf("key-%d", i), i)
	}

	// Act
	var wg sync.WaitGroup
	for i := range numKeys {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			s.Delete(fmt.Sprintf("key-%d", idx))
		}(i)
	}
	wg.Wait()

	// Assert
	assert.That(t, "len must be 0", s.Len(), 0)
}

func Test_SparseSharding_With_DeleteExistingKey_Should_ReturnTrue(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](4)
	s.Put("a", 1)

	// Act
	deleted := s.Delete("a")

	// Assert
	assert.That(t, "delete must return true", deleted, true)
	assert.That(t, "len must be 0", s.Len(), 0)
}

func Test_SparseSharding_With_DeleteMissingKey_Should_ReturnFalse(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](4)

	// Act
	deleted := s.Delete("missing")

	// Assert
	assert.That(t, "delete must return false", deleted, false)
}

func Test_SparseSharding_With_ForEach_Should_VisitAllElements(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](4)
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

func Test_SparseSharding_With_ForEachEarlyStop_Should_StopIteration(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](1) // Single shard for determinism
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

func Test_SparseSharding_With_ForEachShard_Should_ProvideShardIndex(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](4)
	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	// Act
	shardIndices := make(map[int]bool)
	s.ForEachShard(func(shardIdx int, iterate func(fn func(string, int) bool)) {
		shardIndices[shardIdx] = true
		iterate(func(k string, v int) bool {
			return true
		})
	})

	// Assert
	assert.That(t, "should visit 4 shards", len(shardIndices), 4)
}

func Test_SparseSharding_With_GetExistingKey_Should_ReturnValue(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](4)
	s.Put("a", 42)

	// Act
	value := s.Get("a")

	// Assert
	assert.That(t, "value must not be nil", value != nil, true)
	assert.That(t, "value must be 42", *value, 42)
}

func Test_SparseSharding_With_GetMissingKey_Should_ReturnNil(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](4)

	// Act
	value := s.Get("missing")

	// Assert
	assert.That(t, "value must be nil", value == nil, true)
}

func Test_SparseSharding_With_Has_Should_ReturnCorrectResult(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](4)
	s.Put("a", 1)

	// Act & Assert
	assert.That(t, "has a must be true", s.Has("a"), true)
	assert.That(t, "has b must be false", s.Has("b"), false)
}

func Test_SparseSharding_With_PutDuplicateKey_Should_UpdateValue(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](4)
	s.Put("a", 1)

	// Act
	isNew := s.Put("a", 42)

	// Assert
	assert.That(t, "isNew must be false", isNew, false)
	assert.That(t, "value must be 42", *s.Get("a"), 42)
	assert.That(t, "len must be 1", s.Len(), 1)
}

func Test_SparseSharding_With_PutNewKey_Should_AddElement(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](4)

	// Act
	isNew := s.Put("a", 42)

	// Assert
	assert.That(t, "isNew must be true", isNew, true)
	assert.That(t, "value must be 42", *s.Get("a"), 42)
}

func Test_SparseSharding_With_Values_Should_ReturnAllValues(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseSharding[string, int](4)
	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	// Act
	values := s.Values()

	// Assert
	assert.That(t, "values len must be 3", len(values), 3)
}

func Test_SparseSharding_With_ZeroShards_Should_UseDefault(t *testing.T) {
	// Arrange & Act
	s := efficiency.NewSparseSharding[string, int](0)
	s.Put("a", 1)

	// Assert
	assert.That(t, "value must be 1", *s.Get("a"), 1)
}

func Test_SparseShardingWithCapacity_Should_Work(t *testing.T) {
	// Arrange
	s := efficiency.NewSparseShardingWithCapacity[string, int](4, 100)

	// Act
	s.Put("a", 42)

	// Assert
	assert.That(t, "value must be 42", *s.Get("a"), 42)
}
