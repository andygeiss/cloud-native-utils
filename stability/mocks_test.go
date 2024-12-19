package stability_test

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/andygeiss/cloud-native-utils/service"
)

// Mocks for Function (single input, error output).
func mockAlwaysFails() service.Function[int] {
	return func(ctx context.Context, in int) error {
		return errors.New("error")
	}
}

func mockAlwaysSucceeds() service.Function[int] {
	return func(ctx context.Context, in int) error {
		return nil
	}
}

func mockCancel() service.Function[int] {
	return func(ctx context.Context, in int) error {
		return ctx.Err()
	}
}

func mockSlow(duration time.Duration) service.Function[int] {
	return func(ctx context.Context, in int) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(duration):
			return nil
		}
	}
}

func mockFailsTimes(n int) service.Function[int] {
	var count int
	var mutex sync.Mutex
	return func(ctx context.Context, in int) error {
		mutex.Lock()
		defer mutex.Unlock()
		if count >= n {
			return nil
		}
		count++
		return errors.New("error")
	}
}

func mockSucceedsTimes(n int) service.Function[int] {
	var count int
	var mutex sync.Mutex
	return func(ctx context.Context, in int) error {
		mutex.Lock()
		defer mutex.Unlock()
		if count >= n {
			return errors.New("error")
		}
		count++
		return nil
	}
}

// Mocks for Function2 (input, output, and error).
func mockAlwaysFails2() service.Function2[int, int] {
	return func(ctx context.Context, in int) (int, error) {
		return 0, errors.New("error")
	}
}

func mockAlwaysSucceeds2() service.Function2[int, int] {
	return func(ctx context.Context, in int) (int, error) {
		return 42, nil
	}
}

func mockCancel2() service.Function2[int, int] {
	return func(ctx context.Context, in int) (int, error) {
		return in, ctx.Err()
	}
}

func mockSlow2(duration time.Duration) service.Function2[int, int] {
	return func(ctx context.Context, in int) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(duration):
			return in * 2, nil
		}
	}
}

func mockFailsTimes2(n int) service.Function2[int, int] {
	var count int
	var mutex sync.Mutex
	return func(ctx context.Context, in int) (int, error) {
		mutex.Lock()
		defer mutex.Unlock()
		if count >= n {
			return 42, nil
		}
		count++
		return 0, errors.New("error")
	}
}

func mockSucceedsTimes2(n int) service.Function2[int, int] {
	var count int
	var mutex sync.Mutex
	return func(ctx context.Context, in int) (int, error) {
		mutex.Lock()
		defer mutex.Unlock()
		if count >= n {
			return 0, errors.New("error")
		}
		count++
		return 42, nil
	}
}
