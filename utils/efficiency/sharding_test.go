package efficiency_test

import (
	"cloud-native/utils/efficiency"
	"fmt"
	"testing"
)

func TestSharding(t *testing.T) {
	shards := efficiency.NewSharding[int, int](3)
	shards.Add(0, 42)
	if !shards.Contains(0) {
		t.Error("key must be found")
	}
	shards.Add(1, 21)
	if !shards.Contains(1) {
		t.Error("key must be found")
	}
	shards.Delete(1)
	if shards.Contains(1) {
		t.Error("key must not be found")
	}
	if !shards.Contains(0) {
		t.Error("key must still be found")
	}
}

func TestShardingConcurrency(t *testing.T) {
	shards := efficiency.NewSharding[string, string](3)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			key := fmt.Sprintf("key %d", i)
			value := fmt.Sprintf("value %d", i)
			shards.Add(key, value)
			if !shards.Contains(key) {
				t.Errorf("expected '%s' to be in the shards, but it was not found.", key)
			}
			shards.Delete(key)
			if shards.Contains(key) {
				t.Errorf("expected '%s' to be deleted from the shards, but it was found.", key)
			}
		}(i)
	}
}
