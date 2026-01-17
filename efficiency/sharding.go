package efficiency

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// Shard represents a single partition of a sharded key-value store.
// It holds a map of keys (of type K) to values (of type V), and a mutex for thread-safety.
type Shard[K comparable, V any] struct {
	items map[K]V
	mutex sync.RWMutex
}

// Sharding is a collection of Shard objects, each representing a separate partition
// of the key-value store. This allows for distributing keys across multiple shards.
type Sharding[K comparable, V any] []Shard[K, V]

// NewSharding creates an array of Shard objects with the given number of shards.
func NewSharding[K comparable, V any](num int) Sharding[K, V] {
	shards := make([]Shard[K, V], num)
	for i := range shards {
		shards[i] = Shard[K, V]{items: make(map[K]V)}
	}
	return shards
}

// Delete removes a key-value pair from the appropriate shard.
func (a Sharding[K, V]) Delete(key K) {
	shard := a.getShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	delete(shard.items, key)
}

// Get retrieves a value from the appropriate shard.
func (a Sharding[K, V]) Get(key K) (V, bool) {
	shard := a.getShard(key)
	shard.mutex.RLock()
	defer shard.mutex.RUnlock()
	value, exists := shard.items[key]
	return value, exists
}

// Put adds a key-value pair to the appropriate shard.
func (a Sharding[K, V]) Put(key K, value V) {
	shard := a.getShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	shard.items[key] = value
}

func (a Sharding[K, V]) getIndex(key K) int {
	hash := fnv.New32a()
	_, _ = hash.Write(fmt.Appendf(nil, "%v", key))
	sum := int(hash.Sum32())
	return sum % len(a)
}

func (a Sharding[K, V]) getShard(key K) *Shard[K, V] {
	index := a.getIndex(key)
	return &a[index]
}
