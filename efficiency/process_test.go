package efficiency_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

func TestProcess_Three_Int_Values(t *testing.T) {
	in := []int{1, 2, 3}
	inCh := efficiency.Generate[int](in...)
	outCh, _ := efficiency.Process(inCh, func(ctx context.Context, in int) (out int, err error) {
		// Forward the input to the output channel.
		return in, nil
	})
	sum := 0
	for val := range outCh {
		sum += val
	}
	assert.That(t, "sum must be correct", sum, 6)
}

func TestProcess_Ten_Values(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	inCh := efficiency.Generate[int](in...)
	outCh, _ := efficiency.Process(inCh, func(ctx context.Context, in int) (out int, err error) {
		// Forward the input to the output channel.
		return in, nil
	})
	sum := 0
	for val := range outCh {
		sum += val
	}
	assert.That(t, "sum must be correct", sum, 55)
}

func TestProcess_Error_Handling(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	inCh := efficiency.Generate[int](in...)
	_, errCh := efficiency.Process(inCh, func(ctx context.Context, in int) (out int, err error) {
		// Forward an error to the error channel.
		return 0, errors.New("error")
	})
	select {
	case err := <-errCh:
		assert.That(t, "err must be correct", err.Error(), "error")
	}
}
