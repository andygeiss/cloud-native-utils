package efficiency_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func Test_Merge_With_EmptyChannels_Should_ReturnZeroElements(t *testing.T) {
	// Arrange
	ch1 := efficiency.Generate[int]()
	ch2 := efficiency.Generate[int]()

	// Act
	consumer := efficiency.Merge(ch1, ch2)
	count := 0
	for range consumer {
		count++
	}

	// Assert
	assert.That(t, "count should be 0 for empty channels", count, 0)
}

func Test_Merge_With_NoChannels_Should_ReturnZeroElements(t *testing.T) {
	// Arrange & Act
	consumer := efficiency.Merge[int]()
	count := 0
	for range consumer {
		count++
	}

	// Assert
	assert.That(t, "count should be 0 for no channels", count, 0)
}

func Test_Merge_With_OneProducer_Should_ReturnCorrectSum(t *testing.T) {
	// Arrange
	in := []int{1, 2, 3}
	ch := efficiency.Generate[int](in...)
	producer := efficiency.Split(ch, 1)

	// Act
	consumer := efficiency.Merge(producer...)
	sum := 0
	for val := range consumer {
		sum += val
	}

	// Assert
	assert.That(t, "sum must be correct", sum, 6)
}

func Test_Merge_With_SingleChannel_Should_ReturnCorrectSum(t *testing.T) {
	// Arrange
	ch := efficiency.Generate[int](1, 2, 3)

	// Act
	consumer := efficiency.Merge(ch)
	sum := 0
	for val := range consumer {
		sum += val
	}

	// Assert
	assert.That(t, "sum must be correct", sum, 6)
}

func Test_Merge_With_ThreeProducers_Should_ReturnCorrectSum(t *testing.T) {
	// Arrange
	in := []int{1, 2, 3}
	ch := efficiency.Generate[int](in...)
	producer := efficiency.Split(ch, 3)

	// Act
	consumer := efficiency.Merge(producer...)
	sum := 0
	for val := range consumer {
		sum += val
	}

	// Assert
	assert.That(t, "sum must be correct", sum, 6)
}
