package resource_test

import (
	"context"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
)

func Test_InMemoryAccess_With_CreateContextCancelled_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	err := a.Create(ctx, "key", 42)

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be context canceled", err.Error(), "context canceled")
}

func Test_InMemoryAccess_With_CreateDuplicateKey_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)

	// Act
	err := a.Create(ctx, "key", 21)

	// Assert
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceAlreadyExists)
}

func Test_InMemoryAccess_With_CreateValidKey_Should_Succeed(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)

	// Assert
	assert.That(t, "err must be nil", err, nil)
}

func Test_InMemoryAccess_With_DeleteContextCancelled_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	err := a.Delete(ctx, "key")

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be context canceled", err.Error(), "context canceled")
}

func Test_InMemoryAccess_With_DeleteMissingKey_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()

	// Act
	err := a.Delete(ctx, "key")

	// Assert
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}

func Test_InMemoryAccess_With_DeleteValidKey_Should_RemoveResource(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)

	// Act
	err := a.Delete(ctx, "key")
	value, _ := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be nil", value == nil, true)
}

func Test_InMemoryAccess_With_ReadAllContextCancelled_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	_, err := a.ReadAll(ctx)

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be context canceled", err.Error(), "context canceled")
}

func Test_InMemoryAccess_With_ReadAllEmptyStore_Should_ReturnEmptySlice(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()

	// Act
	values, err := a.ReadAll(ctx)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "values must be empty", len(values), 0)
}

func Test_InMemoryAccess_With_ReadAllMultipleKeys_Should_ReturnAllValues(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	_ = a.Create(ctx, "key1", 42)
	_ = a.Create(ctx, "key2", 21)

	// Act
	values, _ := a.ReadAll(ctx)

	// Assert
	assert.That(t, "values len must be 2", len(values), 2)
}

func Test_InMemoryAccess_With_ReadContextCancelled_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	_, err := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be context canceled", err.Error(), "context canceled")
}

func Test_InMemoryAccess_With_ReadMissingKey_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()

	// Act
	_, err := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}

func Test_InMemoryAccess_With_ReadValidKey_Should_ReturnValue(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)

	// Act
	v, err := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "v must be 42", *v, 42)
}

func Test_InMemoryAccess_With_UpdateContextCancelled_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	err := a.Update(ctx, "key", 42)

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "err must be context canceled", err.Error(), "context canceled")
}

func Test_InMemoryAccess_With_UpdateMissingKey_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()

	// Act
	err := a.Update(ctx, "key", 21)

	// Assert
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}

func Test_InMemoryAccess_With_UpdateValidKey_Should_UpdateValue(t *testing.T) {
	// Arrange
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)

	// Act
	err := a.Update(ctx, "key", 21)
	value, _ := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be 21", *value, 21)
}
