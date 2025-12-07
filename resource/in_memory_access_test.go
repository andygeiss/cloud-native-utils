package resource_test

import (
	"context"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
)

func Test_InMemoryAccess_Create_OK(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	assert.That(t, "err must be nil", err, nil)
}

func Test_InMemoryAccess_Create_AlreadyExists(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)
	err := a.Create(ctx, "key", 21)
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceAlreadyExists)
}

func Test_InMemoryAccess_Read_OK(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)
	v, err := a.Read(ctx, "key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "v must be 42", *v, 42)
}

func Test_InMemoryAccess_Read_ResourceNotFound(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	_, err := a.Read(ctx, "key")
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}

func Test_InMemoryAccess_ReadAll_OK(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	_ = a.Create(ctx, "key1", 42)
	_ = a.Create(ctx, "key2", 21)
	values, _ := a.ReadAll(ctx)
	assert.That(t, "values len must be 2", len(values), 2)
}

func Test_InMemoryAccess_Update_OK(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)
	err := a.Update(ctx, "key", 21)
	value, _ := a.Read(ctx, "key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be 21", *value, 21)
}

func Test_InMemoryAccess_Update_ResourceNotFound(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	err := a.Update(ctx, "key", 21)
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}

func Test_InMemoryAccess_Delete_OK(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	_ = a.Create(ctx, "key", 42)
	err := a.Delete(ctx, "key")
	value, _ := a.Read(ctx, "key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be nil", value == nil, true)
}

func Test_InMemoryAccess_Delete_ResourceNotFound(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	ctx := context.Background()
	err := a.Delete(ctx, "key")
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}
