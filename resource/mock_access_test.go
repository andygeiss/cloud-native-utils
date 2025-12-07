package resource_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
)

func TestMockAccessCreate_OK(t *testing.T) {
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithCreateFn(func(ctx context.Context, key int, value int) error {
			isCalled = true
			return nil
		})
	err := a.Create(context.Background(), 0, 42)
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "create must be called", isCalled, true)
}

func TestMockAccessCreate_Error(t *testing.T) {
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithCreateFn(func(ctx context.Context, key int, value int) error {
			isCalled = true
			return errors.New("error")
		})
	err := a.Create(context.Background(), 0, 42)
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "create must be called", isCalled, true)
}

func TestMockAccessRead_OK(t *testing.T) {
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithReadFn(func(ctx context.Context, key int) (*int, error) {
			isCalled = true
			value := 42
			return &value, nil
		})
	value, err := a.Read(context.Background(), 0)
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "read must be called", isCalled, true)
	assert.That(t, "value must be correct", *value, 42)
}

func TestMockAccessRead_Error(t *testing.T) {
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithReadFn(func(ctx context.Context, key int) (*int, error) {
			isCalled = true
			return nil, errors.New("error")
		})
	value, err := a.Read(context.Background(), 0)
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "read must be called", isCalled, true)
	assert.That(t, "value must be nil", value == nil, true)
}

func TestMockAccessReadAll_OK(t *testing.T) {
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithReadAllFn(func(ctx context.Context) ([]int, error) {
			isCalled = true
			value := 42
			return []int{value}, nil
		})
	value, err := a.ReadAll(context.Background())
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "read must be called", isCalled, true)
	assert.That(t, "value must be correct", value, []int{42})
}

func TestMockAccessReadAll_Error(t *testing.T) {
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithReadAllFn(func(ctx context.Context) ([]int, error) {
			isCalled = true
			return nil, errors.New("error")
		})
	value, err := a.ReadAll(context.Background())
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "read must be called", isCalled, true)
	assert.That(t, "value must be nil", value == nil, true)
}

func TestMockAccessUpdate_OK(t *testing.T) {
	isCalled := false
	val := 0
	a := resource.NewMockAccess[int, int]().
		WithUpdateFn(func(ctx context.Context, key int, value int) error {
			isCalled = true
			val = value
			return nil
		})
	err := a.Update(context.Background(), 0, 1)
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "update must be called", isCalled, true)
	assert.That(t, "value must be correct", val, 1)
}

func TestMockAccessUpdate_Error(t *testing.T) {
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithUpdateFn(func(ctx context.Context, key int, value int) error {
			isCalled = true
			return errors.New("error")
		})
	err := a.Update(context.Background(), 0, 42)
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "update must be called", isCalled, true)
}

func TestMockAccessDelete_OK(t *testing.T) {
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithDeleteFn(func(ctx context.Context, key int) error {
			isCalled = true
			return nil
		})
	err := a.Delete(context.Background(), 0)
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "delete must be called", isCalled, true)
}

func TestMockAccessDelete_Error(t *testing.T) {
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithDeleteFn(func(ctx context.Context, key int) error {
			isCalled = true
			return errors.New("error")
		})
	err := a.Delete(context.Background(), 0)
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "delete must be called", isCalled, true)
}
