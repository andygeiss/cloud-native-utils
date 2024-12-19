package efficiency_test

import (
	"fmt"
	"testing"

	"github.com/andygeiss/cloud-native/utils/assert"
	"github.com/andygeiss/cloud-native/utils/efficiency"
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
