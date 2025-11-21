package resource_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
)

func TestInMemoryAccess_Create_OK(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	err := a.Create("key", 42)
	assert.That(t, "err must be nil", err, nil)
}

func TestInMemoryAccess_Create_AlreadyExists(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	_ = a.Create("key", 42)
	err := a.Create("key", 21)
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceAlreadyExists)
}

func TestInMemoryAccess_Read_OK(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	_ = a.Create("key", 42)
	v, err := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "v must be 42", *v, 42)
}

func TestInMemoryAccess_Read_ResourceNotFound(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	_, err := a.Read("key")
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}

func TestInMemoryAccess_ReadAll_OK(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	_ = a.Create("key1", 42)
	_ = a.Create("key2", 21)
	values, _ := a.ReadAll()
	assert.That(t, "values len must be 2", len(values), 2)
}

func TestInMemoryAccess_Update_OK(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	_ = a.Create("key", 42)
	err := a.Update("key", 21)
	value, _ := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be 21", *value, 21)
}

func TestInMemoryAccess_Update_ResourceNotFound(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	err := a.Update("key", 21)
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}

func TestInMemoryAccess_Delete_OK(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	_ = a.Create("key", 42)
	err := a.Delete("key")
	value, _ := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be nil", value == nil, true)
}

func TestInMemoryAccess_Delete_ResourceNotFound(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int]()
	err := a.Delete("key")
	assert.That(t, "err must be correct", err.Error(), resource.ErrorResourceNotFound)
}
