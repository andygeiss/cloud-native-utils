package resource_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
)

func Test_JsonFileAccess_With_CreateAlreadyExists_Should_ReturnError(t *testing.T) {
	// Arrange
	path := "./json_file_access_create_already_exists.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)
	err2 := a.Create(ctx, "key", 21)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceAlreadyExists)
}

func Test_JsonFileAccess_With_CreateDirectoryNotFound_Should_ReturnError(t *testing.T) {
	// Arrange
	path := "./directory_not_found/json_file.json"
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)

	// Assert
	assert.That(t, "err must be correct", errors.Is(err, os.ErrNotExist), true)
}

func Test_JsonFileAccess_With_CreateValidKey_Should_Succeed(t *testing.T) {
	// Arrange
	path := "./json_file_access_create.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)
	err2 := a.Create(ctx, "key", 42)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceAlreadyExists)
}

func Test_JsonFileAccess_With_DeleteMissingKey_Should_ReturnError(t *testing.T) {
	// Arrange
	path := "./json_file_access_delete_not_exists.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)
	err2 := a.Delete(ctx, "key2")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceNotFound)
}

func Test_JsonFileAccess_With_DeleteValidKey_Should_RemoveResource(t *testing.T) {
	// Arrange
	path := "./json_file_access_delete.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)
	err2 := a.Delete(ctx, "key")
	v, err3 := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be nil", err2, nil)
	assert.That(t, "err3 must be correct", err3.Error(), resource.ErrorResourceNotFound)
	assert.That(t, "v must be nil", v == nil, true)
}

func Test_JsonFileAccess_With_ReadAllMultipleKeys_Should_ReturnAllValues(t *testing.T) {
	// Arrange
	path := "./json_file_access_read_all.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key1", 42)
	err2 := a.Create(ctx, "key2", 21)
	m, err3 := a.ReadAll(ctx)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be nil", err2, nil)
	assert.That(t, "err3 must be nil", err3, nil)
	assert.That(t, "m must be 2", len(m), 2)
}

func Test_JsonFileAccess_With_ReadDirectoryNotFound_Should_ReturnError(t *testing.T) {
	// Arrange
	path := "./directory_not_found/json_file.json"
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)

	// Assert
	assert.That(t, "err must be correct", errors.Is(err, os.ErrNotExist), true)
}

func Test_JsonFileAccess_With_ReadMissingKey_Should_ReturnError(t *testing.T) {
	// Arrange
	path := "./json_file_access_read_not_exists.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)
	v, err2 := a.Read(ctx, "key2")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceNotFound)
	assert.That(t, "v must be nil", v == nil, true)
}

func Test_JsonFileAccess_With_ReadValidKey_Should_ReturnValue(t *testing.T) {
	// Arrange
	path := "./json_file_access_read.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)
	v, err2 := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be nil", err2, nil)
	assert.That(t, "v must be 42", *v, 42)
}

func Test_JsonFileAccess_With_UpdateMissingKey_Should_ReturnError(t *testing.T) {
	// Arrange
	path := "./json_file_access_update_not_exists.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)
	err2 := a.Update(ctx, "key2", 21)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceNotFound)
}

func Test_JsonFileAccess_With_UpdateValidKey_Should_UpdateValue(t *testing.T) {
	// Arrange
	path := "./json_file_access_update.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()

	// Act
	err := a.Create(ctx, "key", 42)
	err2 := a.Update(ctx, "key", 21)
	v, err3 := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be nil", err2, nil)
	assert.That(t, "err3 must be nil", err3, nil)
	assert.That(t, "v must be 21", *v, 21)
}
