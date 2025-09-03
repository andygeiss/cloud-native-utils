package efficiency_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func TestSharding(t *testing.T) {
	shards := efficiency.NewSharding[string, int](3)
	key, value := "0", 42
	shards.Put(key, value)
	value, exists := shards.Get(key)
	assert.That(t, "key found", exists, true)
	shards.Delete(key)
	_, exists = shards.Get(key)
	assert.That(t, "value must be correct", value, 42)
	assert.That(t, "key not found", !exists, true)
}

func TestSharding_Concurrency(t *testing.T) {
	shards := efficiency.NewSharding[string, string](3)
	for i := 0; i < 1000; i++ {
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
	for i := 0; i < runtime.NumCPU(); i++ {
		for j := 0; j < b.N; j++ {
			key := fmt.Sprintf("key %d %d", i, j)
			users.Users[key] = User{ID: key, Name: fmt.Sprintf("value %d", i)}
		}
	}

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		// The goroutine uses a mutex to ensure thread safety when accessing the map.
		go func(i int) {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
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
	for i := 0; i < runtime.NumCPU(); i++ {
		for j := 0; j < b.N; j++ {
			key := fmt.Sprintf("key %d %d", i, j)
			users.Users[key] = User{ID: key, Name: fmt.Sprintf("value %d", i)}
		}
	}

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		// The goroutine uses a mutex to ensure thread safety when accessing the map.
		go func(i int) {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
				key := fmt.Sprintf("key %d %d", i, j)
				users.mutex.Lock()
				_, ok := users.Users[key]
				users.mutex.Unlock()
				if !ok {
					b.Errorf("key %s not found", key)
				}
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkMap_Put(b *testing.B) {
	access := UserAccess{
		Users: make(map[string]User),
	}

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		// The goroutine uses a mutex to ensure thread safety when accessing the map.
		go func(i int) {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
				key := fmt.Sprintf("key %d %d", i, j)
				value := fmt.Sprintf("value %d %d", i, j)
				access.mutex.Lock()
				access.Users[key] = User{ID: key, Name: value}
				access.mutex.Unlock()
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkSharding_Delete(b *testing.B) {
	// The sharding implementation handles concurrency internally.
	// We use one shard per CPU core to distribute the load evenly and maximize throughput.
	shards := efficiency.NewSharding[string, User](runtime.NumCPU())

	// Initialize the map.
	for i := 0; i < runtime.NumCPU(); i++ {
		for j := 0; j < b.N; j++ {
			key := fmt.Sprintf("key %d %d", i, j)
			shards.Put(key, User{ID: key, Name: fmt.Sprintf("value %d", i)})
		}
	}

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		// The goroutine uses the sharding implementation to delete values concurrently.
		go func(i int) {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
				id := fmt.Sprintf("%d %d", i, j)
				shards.Delete(id)
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkSharding_Get(b *testing.B) {
	// The sharding implementation handles concurrency internally.
	// We use one shard per CPU core to distribute the load evenly and maximize throughput.
	shards := efficiency.NewSharding[string, User](runtime.NumCPU())

	// Initialize the map.
	for i := 0; i < runtime.NumCPU(); i++ {
		for j := 0; j < b.N; j++ {
			key := fmt.Sprintf("key %d %d", i, j)
			shards.Put(key, User{ID: key, Name: fmt.Sprintf("value %d", i)})
		}
	}

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		// The goroutine uses the sharding implementation to get values concurrently.
		go func(i int) {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
				id := fmt.Sprintf("%d %d", i, j)
				shards.Get(id)
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkSharding_Put(b *testing.B) {
	// The sharding implementation handles concurrency internally.
	// We use one shard per CPU core to distribute the load evenly and maximize throughput.
	shards := efficiency.NewSharding[string, User](runtime.NumCPU())

	b.ResetTimer()

	// Spawn a goroutine for each CPU core to perform concurrent operations.
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		// The goroutine uses the sharding implementation to put values concurrently.
		go func(i int) {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
				id := fmt.Sprintf("%d %d", i, j)
				name := fmt.Sprintf("user %d %d", i, j)
				shards.Put(
					id,
					User{ID: id, Name: name},
				)
			}
		}(i)
	}
	wg.Wait()
}
