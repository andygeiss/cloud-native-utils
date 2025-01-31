package resource_test

import (
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
)

func TestJsonFileAccess_Create(t *testing.T) {
	path := "./json_file_access_create.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	err := a.Create("key", 42)
	err2 := a.Create("key", 42)
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceAlreadyExists)
}

func TestJsonFileAccess_Read(t *testing.T) {
	path := "./json_file_access_read.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	_ = a.Create("key", 42)
	v, err := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "v must be 42", *v, 42)
}

func TestJsonFileAccess_Update(t *testing.T) {
	path := "./json_file_access_update.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	_ = a.Create("key", 42)
	err := a.Update("key", 21)
	v, _ := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "v must be 21", *v, 21)
}

func TestJsonFileAccess_Delete(t *testing.T) {
	path := "./json_file_access_delete.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	_ = a.Create("key", 42)
	err := a.Delete("key")
	v, err2 := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceNotFound)
	assert.That(t, "v must be nil", v == nil, true)
}
