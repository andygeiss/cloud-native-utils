package efficiency_test

import (
	"cloud-native/utils/efficiency"
	"fmt"
	"testing"
)

func TestSharding(t *testing.T) {
	shards := efficiency.NewSharding[string, int](3)
	key, value := "0", 42
	shards.Put(key, value)
	value, exists := shards.Get(key)
	if !exists {
		t.Errorf("expected '%s' to be in the shards, but it was not found.", key)
	}
	if value != 42 {
		t.Errorf("value must be correct, but got %v", value)
	}
	shards.Delete(key)
	_, exists = shards.Get(key)
	if exists {
		t.Errorf("expected '%s' to be not in the shards, but it was found.", key)
	}
}

func TestSharding_Concurrency(t *testing.T) {
	shards := efficiency.NewSharding[string, string](3)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			key := fmt.Sprintf("key %d", i)
			value := fmt.Sprintf("value %d", i)
			shards.Put(key, value)
			result, exists := shards.Get(key)
			if !exists {
				t.Errorf("expected '%s' to be in the shards, but it was not found.", key)
			}
			if value != result {
				t.Errorf("value must be correct, but got %v", value)
			}
			shards.Delete(key)
			_, exists = shards.Get(key)
			if exists {
				t.Errorf("expected '%s' to be not in the shards, but it was found.", key)
			}
		}(i)
	}
}
