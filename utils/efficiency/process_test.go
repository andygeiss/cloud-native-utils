package efficiency_test

import (
	"cloud-native/utils/efficiency"
	"context"
	"errors"
	"testing"
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
	if sum != 6 {
		t.Fatalf("sum must be 6, but got %d", sum)
	}
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
	if sum != 55 {
		t.Fatalf("sum must be 55, but got %d", sum)
	}
}

func TestProcess_Error_Handling(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	inCh := efficiency.Generate[int](in...)
	outCh, _ := efficiency.Process(inCh, func(ctx context.Context, in int) (out int, err error) {
		// Forward an error to the error channel.
		return 0, errors.New("error")
	})
	sum := 0
	for val := range outCh {
		sum += val
	}
	if sum != 55 {
		t.Fatalf("sum must be 55, but got %d", sum)
	}
}
