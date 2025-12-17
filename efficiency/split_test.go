package efficiency_test

import (
	"sync"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func Test_Split_With_EmptyInput_Should_ReturnEmptyConsumerChannels(t *testing.T) {
	// Arrange
	producer := efficiency.Generate[int]()

	// Act
	consumer := efficiency.Split(producer, 2)

	// Assert
	assert.That(t, "should have 2 consumer channels", len(consumer), 2)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for range consumer[0] {
		}
	}()
	go func() {
		defer wg.Done()
		for range consumer[1] {
		}
	}()
	wg.Wait()
}

func Test_Split_With_OneConsumer_Should_ReturnCorrectSum(t *testing.T) {
	// Arrange
	in := []int{1, 2, 3}
	producer := efficiency.Generate[int](in...)

	// Act
	consumer := efficiency.Split(producer, 1)
	sum := 0
	for range 3 {
		val := <-consumer[0]
		sum += val
	}

	// Assert
	assert.That(t, "sum must be correct", sum, 6)
}

func Test_Split_With_TwoConsumers_Should_ReturnCorrectSum(t *testing.T) {
	// Arrange
	in := []int{1, 2, 3, 5}
	producer := efficiency.Generate[int](in...)

	// Act
	consumer := efficiency.Split(producer, 2)

	var mu sync.Mutex
	sum := 0
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for val := range consumer[0] {
			mu.Lock()
			sum += val
			mu.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		for val := range consumer[1] {
			mu.Lock()
			sum += val
			mu.Unlock()
		}
	}()

	wg.Wait()

	// Assert
	assert.That(t, "sum must be correct", sum, 11)
}

func Test_Split_With_ZeroConsumers_Should_ReturnEmptySlice(t *testing.T) {
	// Arrange
	producer := efficiency.Generate[int](1, 2, 3)

	// Act
	consumer := efficiency.Split(producer, 0)

	// Assert
	assert.That(t, "should return empty slice", len(consumer), 0)
}
