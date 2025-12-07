package resource_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
)

func Test_JsonFileAccess_Create_OK(t *testing.T) {
	path := "./json_file_access_create.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	err2 := a.Create(ctx, "key", 42)
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceAlreadyExists)
}

func Test_JsonFileAccess_Create_Already_Exists(t *testing.T) {
	path := "./json_file_access_create_already_exists.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	err2 := a.Create(ctx, "key", 21)
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceAlreadyExists)
}

func Test_JsonFileAccess_Create_Directory_Not_Found(t *testing.T) {
	path := "./directory_not_found/json_file.json"
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	assert.That(t, "err must be correct", errors.Is(err, os.ErrNotExist), true)
}

func Test_JsonFileAccess_Read_OK(t *testing.T) {
	path := "./json_file_access_read.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	v, err2 := a.Read(ctx, "key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be nil", err2, nil)
	assert.That(t, "v must be 42", *v, 42)
}

func Test_JsonFileAccess_ReadAll_OK(t *testing.T) {
	path := "./json_file_access_read_all.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key1", 42)
	err2 := a.Create(ctx, "key2", 21)
	m, err3 := a.ReadAll(ctx)
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be nil", err2, nil)
	assert.That(t, "err3 must be nil", err3, nil)
	assert.That(t, "m must be 2", len(m), 2)
}

func Test_JsonFileAccess_Read_Directory_Not_Found(t *testing.T) {
	path := "./directory_not_found/json_file.json"
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	assert.That(t, "err must be correct", errors.Is(err, os.ErrNotExist), true)
}

func Test_JsonFileAccess_Read_ResourceNotFound(t *testing.T) {
	path := "./json_file_access_read_not_exists.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	v, err2 := a.Read(ctx, "key2")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceNotFound)
	assert.That(t, "v must be nil", v == nil, true)
}

func Test_JsonFileAccess_Update_OK(t *testing.T) {
	path := "./json_file_access_update.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	err2 := a.Update(ctx, "key", 21)
	v, err3 := a.Read(ctx, "key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be nil", err2, nil)
	assert.That(t, "err3 must be nil", err3, nil)
	assert.That(t, "v must be 21", *v, 21)
}

func Test_JsonFileAccess_Update_DirectoryNotFound(t *testing.T) {
	path := "./directory_not_found/json_file.json"
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	err2 := a.Update(ctx, "key", 21)
	assert.That(t, "err must be correct", errors.Is(err, os.ErrNotExist), true)
	assert.That(t, "err2 must be correct", errors.Is(err2, os.ErrNotExist), true)
}

func Test_JsonFileAccess_Update_ResourceNotFound(t *testing.T) {
	path := "./json_file_access_update_not_exists.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	err2 := a.Update(ctx, "key2", 21)
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceNotFound)
}

func Test_JsonFileAccess_Delete_OK(t *testing.T) {
	path := "./json_file_access_delete.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	err2 := a.Delete(ctx, "key")
	v, err3 := a.Read(ctx, "key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be nil", err2, nil)
	assert.That(t, "err3 must be correct", err3.Error(), resource.ErrorResourceNotFound)
	assert.That(t, "v must be nil", v == nil, true)
}

func Test_JsonFileAccess_Delete_DirectoryNotFound(t *testing.T) {
	path := "./directory_not_found/json_file.json"
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	err2 := a.Delete(ctx, "key")
	assert.That(t, "err must be correct", errors.Is(err, os.ErrNotExist), true)
	assert.That(t, "err2 must be correct", errors.Is(err2, os.ErrNotExist), true)
}

func Test_JsonFileAccess_Delete_ResourceNotFound(t *testing.T) {
	path := "./json_file_access_delete_not_exists.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	err2 := a.Delete(ctx, "key2")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceNotFound)
}

func Test_JsonFileAccess_Delete_ResourceNotFoundWithContext(t *testing.T) {
	path := "./json_file_access_delete_not_exists.json"
	defer os.Remove(path)
	a := resource.NewJsonFileAccess[string, int](path)
	ctx := context.Background()
	err := a.Create(ctx, "key", 42)
	err2 := a.Delete(ctx, "key2")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be correct", err2.Error(), resource.ErrorResourceNotFound)
}
