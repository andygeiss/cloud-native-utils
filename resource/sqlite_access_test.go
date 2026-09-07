package resource_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
	_ "modernc.org/sqlite"
)

// newTestDB opens an empty SQLite database for one test. Every test gets its
// own file, so the tests hold no order between them and leave nothing behind.
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatalf("opening the database failed: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func Test_SqliteAccess_With_CreateDuplicateKey_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewSqliteAccess[string, string](newTestDB(t))
	ctx := context.Background()
	_ = a.Init(ctx)
	_ = a.Create(ctx, "key", "value")

	// Act
	err := a.Create(ctx, "key", "value")

	// Assert
	assert.That(t, "err must be correct", err.Error(), "constraint failed: UNIQUE constraint failed: kv_store.key (1555)")
}

func Test_SqliteAccess_With_CreateValidKey_Should_Succeed(t *testing.T) {
	// Arrange
	a := resource.NewSqliteAccess[string, string](newTestDB(t))
	ctx := context.Background()

	// Act
	err := a.Init(ctx)
	err2 := a.Create(ctx, "key", "value")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be nil", err2, nil)
}

func Test_SqliteAccess_With_DeleteValidKey_Should_RemoveValue(t *testing.T) {
	// Arrange
	a := resource.NewSqliteAccess[string, string](newTestDB(t))
	ctx := context.Background()
	_ = a.Init(ctx)
	_ = a.Create(ctx, "key", "value")

	// Act
	err := a.Delete(ctx, "key")
	_, err2 := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must not be nil", err2 != nil, true)
}

func Test_SqliteAccess_With_ReadAllMultipleKeys_Should_ReturnAllValues(t *testing.T) {
	// Arrange
	a := resource.NewSqliteAccess[string, string](newTestDB(t))
	ctx := context.Background()
	_ = a.Init(ctx)
	_ = a.Create(ctx, "key1", "value1")
	_ = a.Create(ctx, "key2", "value2")

	// Act
	values, err := a.ReadAll(ctx)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "values must be ['value1', 'value2']", values, []string{"value1", "value2"})
}

func Test_SqliteAccess_With_ReadMissingKey_Should_ReturnError(t *testing.T) {
	// Arrange
	a := resource.NewSqliteAccess[string, string](newTestDB(t))
	ctx := context.Background()
	_ = a.Init(ctx)
	_ = a.Create(ctx, "key", "value")

	// Act
	_, err := a.Read(ctx, "key2")

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
}

func Test_SqliteAccess_With_ReadValidKey_Should_ReturnValue(t *testing.T) {
	// Arrange
	a := resource.NewSqliteAccess[string, string](newTestDB(t))
	ctx := context.Background()
	_ = a.Init(ctx)
	_ = a.Create(ctx, "key", "value")

	// Act
	value, err := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be 'value'", *value, "value")
}

func Test_SqliteAccess_With_UpdateValidKey_Should_UpdateValue(t *testing.T) {
	// Arrange
	a := resource.NewSqliteAccess[string, string](newTestDB(t))
	ctx := context.Background()
	_ = a.Init(ctx)
	_ = a.Create(ctx, "key", "value")

	// Act
	err := a.Update(ctx, "key", "value2")
	value, _ := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be 'value2'", *value, "value2")
}
