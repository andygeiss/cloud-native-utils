package resource_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
)

func Test_MockAccess_With_CreateError_Should_ReturnError(t *testing.T) {
	// Arrange
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithCreateFn(func(ctx context.Context, key int, value int) error {
			isCalled = true
			return errors.New("error")
		})

	// Act
	err := a.Create(context.Background(), 0, 42)

	// Assert
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "create must be called", isCalled, true)
}

func Test_MockAccess_With_CreateSuccess_Should_Succeed(t *testing.T) {
	// Arrange
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithCreateFn(func(ctx context.Context, key int, value int) error {
			isCalled = true
			return nil
		})

	// Act
	err := a.Create(context.Background(), 0, 42)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "create must be called", isCalled, true)
}

func Test_MockAccess_With_DeleteError_Should_ReturnError(t *testing.T) {
	// Arrange
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithDeleteFn(func(ctx context.Context, key int) error {
			isCalled = true
			return errors.New("error")
		})

	// Act
	err := a.Delete(context.Background(), 0)

	// Assert
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "delete must be called", isCalled, true)
}

func Test_MockAccess_With_DeleteSuccess_Should_Succeed(t *testing.T) {
	// Arrange
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithDeleteFn(func(ctx context.Context, key int) error {
			isCalled = true
			return nil
		})

	// Act
	err := a.Delete(context.Background(), 0)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "delete must be called", isCalled, true)
}

func Test_MockAccess_With_ReadAllError_Should_ReturnError(t *testing.T) {
	// Arrange
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithReadAllFn(func(ctx context.Context) ([]int, error) {
			isCalled = true
			return nil, errors.New("error")
		})

	// Act
	value, err := a.ReadAll(context.Background())

	// Assert
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "read must be called", isCalled, true)
	assert.That(t, "value must be nil", value == nil, true)
}

func Test_MockAccess_With_ReadAllSuccess_Should_ReturnValues(t *testing.T) {
	// Arrange
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithReadAllFn(func(ctx context.Context) ([]int, error) {
			isCalled = true
			value := 42
			return []int{value}, nil
		})

	// Act
	value, err := a.ReadAll(context.Background())

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "read must be called", isCalled, true)
	assert.That(t, "value must be correct", value, []int{42})
}

func Test_MockAccess_With_ReadError_Should_ReturnError(t *testing.T) {
	// Arrange
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithReadFn(func(ctx context.Context, key int) (*int, error) {
			isCalled = true
			return nil, errors.New("error")
		})

	// Act
	value, err := a.Read(context.Background(), 0)

	// Assert
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "read must be called", isCalled, true)
	assert.That(t, "value must be nil", value == nil, true)
}

func Test_MockAccess_With_ReadSuccess_Should_ReturnValue(t *testing.T) {
	// Arrange
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithReadFn(func(ctx context.Context, key int) (*int, error) {
			isCalled = true
			value := 42
			return &value, nil
		})

	// Act
	value, err := a.Read(context.Background(), 0)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "read must be called", isCalled, true)
	assert.That(t, "value must be correct", *value, 42)
}

func Test_MockAccess_With_UpdateError_Should_ReturnError(t *testing.T) {
	// Arrange
	isCalled := false
	a := resource.NewMockAccess[int, int]().
		WithUpdateFn(func(ctx context.Context, key int, value int) error {
			isCalled = true
			return errors.New("error")
		})

	// Act
	err := a.Update(context.Background(), 0, 1)

	// Assert
	assert.That(t, "err must be correct", err.Error(), "error")
	assert.That(t, "update must be called", isCalled, true)
}

func Test_MockAccess_With_UpdateSuccess_Should_UpdateValue(t *testing.T) {
	// Arrange
	isCalled := false
	val := 0
	a := resource.NewMockAccess[int, int]().
		WithUpdateFn(func(ctx context.Context, key int, value int) error {
			isCalled = true
			val = value
			return nil
		})

	// Act
	err := a.Update(context.Background(), 0, 1)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "update must be called", isCalled, true)
	assert.That(t, "val must be correct", val, 1)
}
