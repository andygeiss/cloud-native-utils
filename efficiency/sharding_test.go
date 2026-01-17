package efficiency_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func Test_Sharding_With_ConcurrentAccess_Should_HandleSafely(t *testing.T) {
	// Arrange
	shards := efficiency.NewSharding[string, string](3)
	// Act
	for i := range 1000 {
		go func(i int) {
			key := fmt.Sprintf("key %d", i)
			value := fmt.Sprintf("value %d", i)
			shards.Put(key, value)
			result, exists := shards.Get(key)
			assert.That(t, "key found", exists, true)
			shards.Delete(key)
			_, exists = shards.Get(key)
			assert.That(t, "value must be correct", value, result)
			assert.That(t, "key not found", !exists, true)
		}(i)
	}
	// Assert - concurrent operations completed without panic
}

func Test_Sharding_With_DeletedKey_Should_NotExist(t *testing.T) {
	// Arrange
	shards := efficiency.NewSharding[string, int](3)
	key := "0"
	shards.Put(key, 42)
	// Act
	shards.Delete(key)
	_, exists := shards.Get(key)
	// Assert
	assert.That(t, "key not found", !exists, true)
}

func Test_Sharding_With_PutAndGet_Should_ReturnCorrectValue(t *testing.T) {
	// Arrange
	shards := efficiency.NewSharding[string, int](3)
	key, value := "0", 42
	// Act
	shards.Put(key, value)
	result, exists := shards.Get(key)
	// Assert
	assert.That(t, "key found", exists, true)
	assert.That(t, "value must be correct", result, 42)
}

// User represents a simple data structure for benchmarking purposes.
type User struct {
	ID   string
	Name string
}

// UserAccess struct represents a collection of users with a mutex for concurrency control.
// This will be used for the standard library's map.
type UserAccess struct {
	Users map[string]User
	mutex sync.Mutex
}

func BenchmarkMap_Delete(b *testing.B) {
	users := UserAccess{
		Users: make(map[string]User),
	}

	// Initialize the map.
	for i := range runtime.NumCPU() {
		for j := range b.N {
			key := fmt.Sprintf("key %d %d", i, j)
			users.Users[key] = User{ID: key, Name: fmt.Sprintf("value %d", i)}
		}
	}

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := range runtime.NumCPU() {
		wg.Add(1)

		// The goroutine uses a mutex to ensure thread safety when accessing the map.
		go func(i int) {
			defer wg.Done()
			for j := range b.N {
				key := fmt.Sprintf("key %d %d", i, j)
				users.mutex.Lock()
				delete(users.Users, key)
				users.mutex.Unlock()
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkMap_Get(b *testing.B) {
	users := UserAccess{
		Users: make(map[string]User),
	}

	// Initialize the map.
	for i := range runtime.NumCPU() {
		for j := range b.N {
			key := fmt.Sprintf("key %d %d", i, j)
			users.Users[key] = User{ID: key, Name: fmt.Sprintf("value %d", i)}
		}
	}

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := range runtime.NumCPU() {
		wg.Add(1)

		// The goroutine uses a mutex to ensure thread safety when accessing the map.
		go func(i int) {
			defer wg.Done()
			for j := range b.N {
				key := fmt.Sprintf("key %d %d", i, j)
				users.mutex.Lock()
				_ = users.Users[key]
				users.mutex.Unlock()
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkMap_Put(b *testing.B) {
	users := UserAccess{
		Users: make(map[string]User),
	}

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := range runtime.NumCPU() {
		wg.Add(1)

		// The goroutine uses a mutex to ensure thread safety when accessing the map.
		go func(i int) {
			defer wg.Done()
			for j := range b.N {
				key := fmt.Sprintf("key %d %d", i, j)
				users.mutex.Lock()
				users.Users[key] = User{ID: key, Name: fmt.Sprintf("value %d", i)}
				users.mutex.Unlock()
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkSharding_Delete(b *testing.B) {
	shards := efficiency.NewSharding[string, User](32)

	// Initialize the shards.
	for i := range runtime.NumCPU() {
		for j := range b.N {
			key := fmt.Sprintf("key %d %d", i, j)
			shards.Put(key, User{ID: key, Name: fmt.Sprintf("value %d", i)})
		}
	}

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := range runtime.NumCPU() {
		wg.Add(1)

		// Each goroutine operates on the shards.
		go func(i int) {
			defer wg.Done()
			for j := range b.N {
				key := fmt.Sprintf("key %d %d", i, j)
				shards.Delete(key)
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkSharding_Get(b *testing.B) {
	shards := efficiency.NewSharding[string, User](32)

	// Initialize the shards.
	for i := range runtime.NumCPU() {
		for j := range b.N {
			key := fmt.Sprintf("key %d %d", i, j)
			shards.Put(key, User{ID: key, Name: fmt.Sprintf("value %d", i)})
		}
	}

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := range runtime.NumCPU() {
		wg.Add(1)

		// Each goroutine operates on the shards.
		go func(i int) {
			defer wg.Done()
			for j := range b.N {
				key := fmt.Sprintf("key %d %d", i, j)
				_, _ = shards.Get(key)
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkSharding_Put(b *testing.B) {
	shards := efficiency.NewSharding[string, User](32)

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := range runtime.NumCPU() {
		wg.Add(1)

		// Each goroutine operates on the shards.
		go func(i int) {
			defer wg.Done()
			for j := range b.N {
				key := fmt.Sprintf("key %d %d", i, j)
				shards.Put(key, User{ID: key, Name: fmt.Sprintf("value %d", i)})
			}
		}(i)
	}
	wg.Wait()
}
