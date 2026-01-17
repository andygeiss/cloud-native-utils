package efficiency_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func Test_Process_With_ErrorReturned_Should_SendToErrorChannel(t *testing.T) {
	// Arrange
	in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	inCh := efficiency.Generate(in...)

	// Act
	_, errCh := efficiency.Process(inCh, func(_ context.Context, _ int) (int, error) {
		return 0, errors.New("error")
	})
	err := <-errCh

	// Assert
	assert.That(t, "err must be correct", err.Error(), "error")
}

func Test_Process_With_TenValues_Should_SumCorrectly(t *testing.T) {
	// Arrange
	in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	inCh := efficiency.Generate(in...)

	// Act
	outCh, _ := efficiency.Process(inCh, func(_ context.Context, in int) (int, error) {
		return in, nil
	})
	sum := 0
	for val := range outCh {
		sum += val
	}

	// Assert
	assert.That(t, "sum must be correct", sum, 55)
}

func Test_Process_With_ThreeValues_Should_SumCorrectly(t *testing.T) {
	// Arrange
	in := []int{1, 2, 3}
	inCh := efficiency.Generate(in...)

	// Act
	outCh, _ := efficiency.Process(inCh, func(_ context.Context, in int) (int, error) {
		return in, nil
	})
	sum := 0
	for val := range outCh {
		sum += val
	}

	// Assert
	assert.That(t, "sum must be correct", sum, 6)
}
