package resource_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
)

func TestInMemoryAccess_Create(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int](8)
	err := a.Create("key", 42)
	assert.That(t, "err must be nil", err, nil)
}

func TestInMemoryAccess_Read(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int](8)
	_ = a.Create("key", 42)
	v, err := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "v must be 42", *v, 42)
}

func TestInMemoryAccess_ReadAll(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int](8)
	_ = a.Create("key1", 42)
	_ = a.Create("key2", 21)
	values, _ := a.ReadAll()
	assert.That(t, "values must be 42 and 21", values, []int{42, 21})
}

func TestInMemoryAccess_Update(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int](8)
	_ = a.Create("key", 42)
	err := a.Update("key", 21)
	value, _ := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be 21", *value, 21)
}

func TestInMemoryAccess_Delete(t *testing.T) {
	a := resource.NewInMemoryAccess[string, int](8)
	_ = a.Create("key", 42)
	err := a.Delete("key")
	value, _ := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be nil", value == nil, true)
}
