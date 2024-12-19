package stability_test

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/andygeiss/cloud-native-utils/service"
)

func mockAlwaysFails() service.Function[int, int] {
	return func() service.Function[int, int] {
		return func(ctx context.Context, in int) (out int, err error) {
			return out, errors.New("error")
		}
	}()
}

func mockAlwaysSucceeds() service.Function[int, int] {
	return func() service.Function[int, int] {
		return func(ctx context.Context, in int) (int, error) {
			return 42, nil
		}
	}()
}

func mockCancel() service.Function[int, int] {
	return func(ctx context.Context, in int) (int, error) {
		return in, ctx.Err()
	}
}

func mockSlow(duration time.Duration) service.Function[int, int] {
	return func(ctx context.Context, in int) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(duration):
			return in * 2, nil
		}
	}
}

func mockFailsTimes(n int) service.Function[int, int] {
	return func() service.Function[int, int] {
		var count int
		var mutex sync.Mutex
		return func(ctx context.Context, in int) (out int, err error) {
			mutex.Lock()
			defer mutex.Unlock()
			if count >= n {
				return 42, nil
			}
			count++
			return out, errors.New("error")
		}
	}()
}

func mockSucceedsTimes(n int) service.Function[int, int] {
	return func() service.Function[int, int] {
		var count int
		var mutex sync.Mutex
		return func(ctx context.Context, in int) (out int, err error) {
			mutex.Lock()
			defer mutex.Unlock()
			if count >= n {
				return out, errors.New("error")
			}
			count++
			return 42, nil
		}
	}()
}
