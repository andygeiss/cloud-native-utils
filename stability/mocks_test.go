package stability_test

import (
	"cloud-native/service"
	"context"
	"errors"
	"sync"
)

func mockAlwaysFails() service.Function[int] {
	return func() service.Function[int] {
		return func(ctx context.Context) (*int, error) {
			return nil, errors.New("error")
		}
	}()
}

func mockAlwaysSucceeds() service.Function[int] {
	return func() service.Function[int] {
		return func(ctx context.Context) (*int, error) {
			value := 42
			return &value, nil
		}
	}()
}

func mockFailsTimes(n int) service.Function[int] {
	return func() service.Function[int] {
		var count int
		var mutex sync.Mutex
		return func(ctx context.Context) (*int, error) {
			mutex.Lock()
			defer mutex.Unlock()
			if count >= n {
				value := 42
				return &value, nil
			}
			count++
			return nil, errors.New("error")
		}
	}()
}

func mockSucceedsTimes(n int) service.Function[int] {
	return func() service.Function[int] {
		var count int
		var mutex sync.Mutex
		return func(ctx context.Context) (*int, error) {
			mutex.Lock()
			defer mutex.Unlock()
			if count >= n {
				return nil, errors.New("error")
			}
			count++
			value := 42
			return &value, nil
		}
	}()
}
