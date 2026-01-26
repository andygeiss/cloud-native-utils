package resource_test

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
)

func Test_ShardedSparseAccess_With_Clear_Should_RemoveAllElements(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()
	_ = a.Create(ctx, "key1", 1)
	_ = a.Create(ctx, "key2", 2)
	_ = a.Create(ctx, "key3", 3)

	// Act
	a.Clear()

	// Assert
	assert.That(t, "len must be 0", a.Len(), 0)
	values, _ := a.ReadAll(ctx)
	assert.That(t, "ReadAll must return empty slice", len(values), 0)
}

func Test_ShardedSparseAccess_With_ConcurrentCreate_Should_HandleSafely(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](32)
	ctx := context.Background()
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
				_ = a.Create(ctx, key, gid*1000+j)
			}
		}(i)
	}
	wg.Wait()

	// Assert
	expectedLen := numGoroutines * numOpsPerGoroutine
	assert.That(t, "len must match expected", a.Len(), expectedLen)
}

func Test_ShardedSparseAccess_With_ConcurrentDelete_Should_HandleSafely(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](32)
	ctx := context.Background()
	numKeys := 1000

	// Pre-populate data
	for i := range numKeys {
		_ = a.Create(ctx, fmt.Sprintf("key-%d", i), i)
	}

	// Act - concurrent deletes
	var wg sync.WaitGroup
	for i := range numKeys {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_ = a.Delete(ctx, fmt.Sprintf("key-%d", idx))
		}(i)
	}
	wg.Wait()

	// Assert
	assert.That(t, "len must be 0", a.Len(), 0)
}

func Test_ShardedSparseAccess_With_ConcurrentReadWrite_Should_HandleSafely(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](32)
	ctx := context.Background()
	numGoroutines := runtime.NumCPU() * 2

	// Pre-populate some data
	for i := range 100 {
		_ = a.Create(ctx, fmt.Sprintf("key-%d", i), i)
	}

	// Act - concurrent reads and writes
	var wg sync.WaitGroup
	for i := range numGoroutines {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			for j := range 100 {
				key := fmt.Sprintf("key-%d", j)
				if gid%2 == 0 {
					// Reader
					_, _ = a.Read(ctx, key)
				} else {
					// Writer
					_ = a.Update(ctx, key, gid*1000+j)
				}
			}
		}(i)
	}
	wg.Wait()

	// Assert - no panic, data still accessible
	assert.That(t, "len must be 100", a.Len(), 100)
}

func Test_ShardedSparseAccess_With_CreateContextCancelled_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	err := a.Create(ctx, "key", 42)

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be context canceled", err.Error(), "context canceled")
}

func Test_ShardedSparseAccess_With_CreateDuplicateKey_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)

	// Act
	err := a.Create(ctx, "key", 21)

	// Assert
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceAlreadyExists)
}

func Test_ShardedSparseAccess_With_CreateValidKey_Should_Succeed(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)

	// Assert
	assert.That(t, "err must be nil", err, nil)
}

func Test_ShardedSparseAccess_With_DeleteAllElements_Should_BeEmpty(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()
	_ = a.Create(ctx, "key1", 1)
	_ = a.Create(ctx, "key2", 2)
	_ = a.Create(ctx, "key3", 3)

	// Act
	_ = a.Delete(ctx, "key1")
	_ = a.Delete(ctx, "key2")
	_ = a.Delete(ctx, "key3")

	// Assert
	assert.That(t, "len must be 0", a.Len(), 0)
	values, _ := a.ReadAll(ctx)
	assert.That(t, "ReadAll must return empty slice", len(values), 0)
}

func Test_ShardedSparseAccess_With_DeleteContextCancelled_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	err := a.Delete(ctx, "key")

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be context canceled", err.Error(), "context canceled")
}

func Test_ShardedSparseAccess_With_DeleteFirstElement_Should_MaintainOtherElements(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](1)
	ctx := context.Background()
	_ = a.Create(ctx, "key1", 1)
	_ = a.Create(ctx, "key2", 2)
	_ = a.Create(ctx, "key3", 3)

	// Act
	err := a.Delete(ctx, "key1")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	v2, _ := a.Read(ctx, "key2")
	v3, _ := a.Read(ctx, "key3")
	assert.That(t, "key2 must still exist with value 2", *v2, 2)
	assert.That(t, "key3 must still exist with value 3", *v3, 3)
	assert.That(t, "len must be 2", a.Len(), 2)
}

func Test_ShardedSparseAccess_With_DeleteLastElement_Should_ShrinkCorrectly(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](1)
	ctx := context.Background()
	_ = a.Create(ctx, "key1", 1)
	_ = a.Create(ctx, "key2", 2)
	_ = a.Create(ctx, "key3", 3)

	// Act
	err := a.Delete(ctx, "key3")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	v1, _ := a.Read(ctx, "key1")
	v2, _ := a.Read(ctx, "key2")
	assert.That(t, "key1 must still exist with value 1", *v1, 1)
	assert.That(t, "key2 must still exist with value 2", *v2, 2)
	assert.That(t, "len must be 2", a.Len(), 2)
}

func Test_ShardedSparseAccess_With_DeleteMiddleElement_Should_MaintainOtherElements(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](1) // Single shard to ensure deterministic behavior
	ctx := context.Background()
	_ = a.Create(ctx, "key1", 1)
	_ = a.Create(ctx, "key2", 2)
	_ = a.Create(ctx, "key3", 3)

	// Act
	err := a.Delete(ctx, "key2")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	v1, _ := a.Read(ctx, "key1")
	v3, _ := a.Read(ctx, "key3")
	assert.That(t, "key1 must still exist with value 1", *v1, 1)
	assert.That(t, "key3 must still exist with value 3", *v3, 3)
	assert.That(t, "len must be 2", a.Len(), 2)
}

func Test_ShardedSparseAccess_With_DeleteMissingKey_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()

	// Act
	err := a.Delete(ctx, "key")

	// Assert
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}

func Test_ShardedSparseAccess_With_DeleteValidKey_Should_RemoveResource(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)

	// Act
	err := a.Delete(ctx, "key")
	value, _ := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be nil", value == nil, true)
}

func Test_ShardedSparseAccess_With_ForEach_Should_VisitAllElements(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()
	_ = a.Create(ctx, "key1", 1)
	_ = a.Create(ctx, "key2", 2)
	_ = a.Create(ctx, "key3", 3)

	// Act
	visited := make(map[string]int)
	a.ForEach(func(k string, v int) bool {
		visited[k] = v
		return true
	})

	// Assert
	assert.That(t, "visited len must be 3", len(visited), 3)
	assert.That(t, "key1 must be visited", visited["key1"], 1)
	assert.That(t, "key2 must be visited", visited["key2"], 2)
	assert.That(t, "key3 must be visited", visited["key3"], 3)
}

func Test_ShardedSparseAccess_With_ForEachEarlyStop_Should_StopIteration(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](1) // Single shard for deterministic order
	ctx := context.Background()
	_ = a.Create(ctx, "key1", 1)
	_ = a.Create(ctx, "key2", 2)
	_ = a.Create(ctx, "key3", 3)

	// Act
	count := 0
	a.ForEach(func(k string, v int) bool {
		count++
		return count < 2 // Stop after 2 elements
	})

	// Assert
	assert.That(t, "count must be 2", count, 2)
}

func Test_ShardedSparseAccess_With_LenEmptyStore_Should_ReturnZero(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)

	// Act
	length := a.Len()

	// Assert
	assert.That(t, "len must be 0", length, 0)
}

func Test_ShardedSparseAccess_With_LenMultipleKeys_Should_ReturnCorrectCount(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()
	_ = a.Create(ctx, "key1", 1)
	_ = a.Create(ctx, "key2", 2)
	_ = a.Create(ctx, "key3", 3)

	// Act
	length := a.Len()

	// Assert
	assert.That(t, "len must be 3", length, 3)
}

func Test_ShardedSparseAccess_With_NegativeShards_Should_UseDefaultShardCount(t *testing.T) {
	// Arrange & Act
	a := resource.NewShardedSparseAccess[string, int](-5)
	ctx := context.Background()

	// Assert - should work correctly with default shards
	err := a.Create(ctx, "key", 42)
	assert.That(t, "err must be nil", err, nil)
	v, _ := a.Read(ctx, "key")
	assert.That(t, "value must be 42", *v, 42)
}

func Test_ShardedSparseAccess_With_ReadAllContextCancelled_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	_, err := a.ReadAll(ctx)

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be context canceled", err.Error(), "context canceled")
}

func Test_ShardedSparseAccess_With_ReadAllEmptyStore_Should_ReturnEmptySlice(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()

	// Act
	values, err := a.ReadAll(ctx)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "values must be empty", len(values), 0)
}

func Test_ShardedSparseAccess_With_ReadAllMultipleKeys_Should_ReturnAllValues(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()
	_ = a.Create(ctx, "key1", 42)
	_ = a.Create(ctx, "key2", 21)

	// Act
	values, _ := a.ReadAll(ctx)

	// Assert
	assert.That(t, "values len must be 2", len(values), 2)
}

func Test_ShardedSparseAccess_With_ReadContextCancelled_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	_, err := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be context canceled", err.Error(), "context canceled")
}

func Test_ShardedSparseAccess_With_ReadMissingKey_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()

	// Act
	_, err := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}

func Test_ShardedSparseAccess_With_ReadValidKey_Should_ReturnValue(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)

	// Act
	v, err := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "v must be 42", *v, 42)
}

func Test_ShardedSparseAccess_With_UpdateContextCancelled_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	err := a.Update(ctx, "key", 42)

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be context canceled", err.Error(), "context canceled")
}

func Test_ShardedSparseAccess_With_UpdateMissingKey_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()

	// Act
	err := a.Update(ctx, "key", 21)

	// Assert
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}

func Test_ShardedSparseAccess_With_UpdateValidKey_Should_UpdateValue(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccess[string, int](4)
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)

	// Act
	err := a.Update(ctx, "key", 21)
	value, _ := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be 21", *value, 21)
}

func Test_ShardedSparseAccess_With_ZeroShards_Should_UseDefaultShardCount(t *testing.T) {
	// Arrange & Act
	a := resource.NewShardedSparseAccess[string, int](0)
	ctx := context.Background()

	// Assert - should work correctly with default shards
	err := a.Create(ctx, "key", 42)
	assert.That(t, "err must be nil", err, nil)
	v, _ := a.Read(ctx, "key")
	assert.That(t, "value must be 42", *v, 42)
}

func Test_ShardedSparseAccessWithCapacity_Should_Work(t *testing.T) {
	// Arrange
	a := resource.NewShardedSparseAccessWithCapacity[string, int](4, 100)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	v, _ := a.Read(ctx, "key")
	assert.That(t, "value must be 42", *v, 42)
}

// --- Benchmarks ---

func BenchmarkComparison_InMemoryAccess_Concurrent_Create(b *testing.B) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d-%d", b.N, i)
			_ = a.Create(ctx, key, i)
			i++
		}
	})
}

func BenchmarkComparison_InMemoryAccess_Concurrent_Read(b *testing.B) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()

	// Pre-populate
	for i := range 10000 {
		_ = a.Create(ctx, fmt.Sprintf("key-%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i%10000)
			_, _ = a.Read(ctx, key)
			i++
		}
	})
}

func BenchmarkComparison_InMemoryAccess_Create(b *testing.B) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("key-%d", i)
		_ = a.Create(ctx, key, i)
	}
}

func BenchmarkShardedSparseAccess_Concurrent_Create(b *testing.B) {
	a := resource.NewShardedSparseAccess[string, int](32)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d-%d", b.N, i)
			_ = a.Create(ctx, key, i)
			i++
		}
	})
}

func BenchmarkShardedSparseAccess_Concurrent_Mixed(b *testing.B) {
	a := resource.NewShardedSparseAccess[string, int](32)
	ctx := context.Background()

	// Pre-populate
	for i := range 10000 {
		_ = a.Create(ctx, fmt.Sprintf("key-%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i%10000)
			switch i % 4 {
			case 0:
				_, _ = a.Read(ctx, key)
			case 1:
				_ = a.Update(ctx, key, i)
			case 2:
				_, _ = a.Read(ctx, key)
			case 3:
				_, _ = a.Read(ctx, key)
			}
			i++
		}
	})
}

func BenchmarkShardedSparseAccess_Concurrent_Read(b *testing.B) {
	a := resource.NewShardedSparseAccess[string, int](32)
	ctx := context.Background()

	// Pre-populate
	for i := range 10000 {
		_ = a.Create(ctx, fmt.Sprintf("key-%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i%10000)
			_, _ = a.Read(ctx, key)
			i++
		}
	})
}

func BenchmarkShardedSparseAccess_Create(b *testing.B) {
	a := resource.NewShardedSparseAccess[string, int](32)
	ctx := context.Background()

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("key-%d", i)
		_ = a.Create(ctx, key, i)
	}
}

func BenchmarkShardedSparseAccess_Delete(b *testing.B) {
	a := resource.NewShardedSparseAccess[string, int](32)
	ctx := context.Background()

	// Pre-populate
	for i := range b.N {
		_ = a.Create(ctx, fmt.Sprintf("key-%d", i), i)
	}

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("key-%d", i)
		_ = a.Delete(ctx, key)
	}
}

func BenchmarkShardedSparseAccess_Read(b *testing.B) {
	a := resource.NewShardedSparseAccess[string, int](32)
	ctx := context.Background()

	// Pre-populate
	for i := range 10000 {
		_ = a.Create(ctx, fmt.Sprintf("key-%d", i), i)
	}

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("key-%d", i%10000)
		_, _ = a.Read(ctx, key)
	}
}

func BenchmarkShardedSparseAccess_ReadAll(b *testing.B) {
	a := resource.NewShardedSparseAccess[string, int](32)
	ctx := context.Background()

	// Pre-populate
	for i := range 10000 {
		_ = a.Create(ctx, fmt.Sprintf("key-%d", i), i)
	}

	b.ResetTimer()
	for range b.N {
		_, _ = a.ReadAll(ctx)
	}
}

func BenchmarkShardedSparseAccess_Update(b *testing.B) {
	a := resource.NewShardedSparseAccess[string, int](32)
	ctx := context.Background()

	// Pre-populate
	for i := range 10000 {
		_ = a.Create(ctx, fmt.Sprintf("key-%d", i), i)
	}

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("key-%d", i%10000)
		_ = a.Update(ctx, key, i)
	}
}
