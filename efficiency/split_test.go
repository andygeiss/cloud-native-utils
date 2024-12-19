package efficiency_test

import (
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func TestSplit_One_Consumer(t *testing.T) {
	in := []int{1, 2, 3}
	producer := efficiency.Generate[int](in...)
	consumer := efficiency.Split(producer, 1)
	sum := 0
	for range 3 {
		val := <-consumer[0]
		sum += val
	}
	assert.That(t, "sum must be correct", sum, 6)
}

func TestSplit_Two_Consumers(t *testing.T) {
	in := []int{1, 2, 3, 5}
	producer := efficiency.Generate[int](in...)
	consumer := efficiency.Split(producer, 2)
	sum := 0
	go func() {
		for val := range consumer[0] {
			sum += val
		}
	}()
	go func() {
		for val := range consumer[1] {
			sum += val
		}
	}()
	time.Sleep(100 * time.Millisecond)
	assert.That(t, "sum must be correct", sum, 11)
}
