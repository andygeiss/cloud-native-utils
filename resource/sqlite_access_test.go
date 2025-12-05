package resource_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
	_ "modernc.org/sqlite"
)

func init() {
	os.Mkdir("testdata", 0755)
}

func Test_SqliteAccess_Create_Given_Key_Should_Return_Without_Error(t *testing.T) {
	path := "testdata/test.sqlite"
	db, _ := sql.Open("sqlite", path)
	defer db.Close()
	a := resource.NewSqliteAccess[string, string](db)
	err := a.Init()
	err2 := a.Create("key", "value")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be nil", err2, nil)
}

func Test_SqliteAccess_Create_Given_Key_Should_Return_Error_If_Key_Exists(t *testing.T) {
	path := "testdata/test.sqlite"
	db, _ := sql.Open("sqlite", path)
	defer db.Close()
	a := resource.NewSqliteAccess[string, string](db)
	_ = a.Init()
	_ = a.Create("key", "value")
	err := a.Create("key", "value")
	assert.That(t, "err must be correct", err.Error(), "constraint failed: UNIQUE constraint failed: kv_store.key (1555)")
}

func Test_SqliteAccess_Read_Given_Key_Should_Return_Value_Without_Error(t *testing.T) {
	path := "testdata/test.sqlite"
	db, _ := sql.Open("sqlite", path)
	defer db.Close()
	a := resource.NewSqliteAccess[string, string](db)
	_ = a.Init()
	_ = a.Create("key", "value")
	value, err := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be 'value'", value, "value")
}

func Test_SqliteAccess_Read_Given_Key_Should_Return_Error_If_Key_Not_Exists(t *testing.T) {
	path := "testdata/test.sqlite"
	db, _ := sql.Open("sqlite", path)
	defer db.Close()
	a := resource.NewSqliteAccess[string, string](db)
	_ = a.Init()
	_ = a.Create("key", "value")
	value, err := a.Read("key2")
	assert.That(t, "err must not be nil", err != nil, true)
	assert.That(t, "value must be empty", value, "")
}

func Test_SqliteAccess_ReadAll_Given_Keys_Should_Return_List_Without_Error(t *testing.T) {
	path := "testdata/test.sqlite"
	db, _ := sql.Open("sqlite", path)
	defer db.Close()
	a := resource.NewSqliteAccess[string, string](db)
	_ = a.Init()
	_ = a.Create("key1", "value1")
	_ = a.Create("key2", "value2")
	values, err := a.ReadAll()
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "values must be ['value1', 'value2']", values, []string{"value1", "value2"})
}

func Test_SqliteAccess_Update_Given_Key_Should_Update_Value_Without_Error(t *testing.T) {
	path := "testdata/test.sqlite"
	db, _ := sql.Open("sqlite", path)
	defer db.Close()
	a := resource.NewSqliteAccess[string, string](db)
	_ = a.Init()
	_ = a.Create("key", "value")
	err := a.Update("key", "value2")
	value, _ := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be 'value2'", value, "value2")
}

func Test_SqliteAccess_Delete_Given_Key_Should_Delete_Value_Without_Error(t *testing.T) {
	path := "testdata/test.sqlite"
	db, _ := sql.Open("sqlite", path)
	defer db.Close()
	a := resource.NewSqliteAccess[string, string](db)
	_ = a.Init()
	_ = a.Create("key", "value")
	err := a.Delete("key")
	value, err2 := a.Read("key")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must not be nil", err2.Error(), "sql: no rows in result set")
	assert.That(t, "value must be empty", value, "")
}
