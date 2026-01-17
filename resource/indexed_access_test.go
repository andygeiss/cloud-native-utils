package resource //nolint:testpackage // internal package tests for unexported types

import (
	"context"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
)

type testEntity struct {
	ID    string
	Email string
	Role  string
}

func Test_IndexedAccess_With_ValidEntity_Should_Create(t *testing.T) {
	// Arrange
	base := NewInMemoryAccess[string, *testEntity]()
	store := NewIndexedAccess(base)
	store.AddIndex("email", func(e *testEntity) string { return e.Email })
	entity := &testEntity{ID: "1", Email: "test@example.com", Role: "admin"}
	// Act
	err := store.Create(context.Background(), "1", entity)
	// Assert
	assert.That(t, "error", err, nil)
}

func Test_IndexedAccess_With_MultipleEntities_Should_FindByIndex(t *testing.T) {
	// Arrange
	base := NewInMemoryAccess[string, *testEntity]()
	store := NewIndexedAccess(base)
	store.AddIndex("email", func(e *testEntity) string { return e.Email })
	store.AddIndex("role", func(e *testEntity) string { return e.Role })
	entity1 := &testEntity{ID: "1", Email: "alice@example.com", Role: "admin"}
	entity2 := &testEntity{ID: "2", Email: "bob@example.com", Role: "admin"}
	_ = store.Create(context.Background(), "1", entity1)
	_ = store.Create(context.Background(), "2", entity2)
	// Act
	byEmail, _ := store.FindByIndex(context.Background(), "email", "alice@example.com")
	byRole, _ := store.FindByIndex(context.Background(), "role", "admin")
	notFound, _ := store.FindByIndex(context.Background(), "email", "unknown@example.com")
	// Assert
	assert.That(t, "by email count", len(byEmail), 1)
	assert.That(t, "by email ID", byEmail[0].ID, "1")
	assert.That(t, "by role count", len(byRole), 2)
	assert.That(t, "not found", len(notFound), 0)
}

func Test_IndexedAccess_With_ExistingEntity_Should_FindOneByIndex(t *testing.T) {
	// Arrange
	base := NewInMemoryAccess[string, *testEntity]()
	store := NewIndexedAccess(base)
	store.AddIndex("email", func(e *testEntity) string { return e.Email })
	entity := &testEntity{ID: "1", Email: "alice@example.com", Role: "admin"}
	_ = store.Create(context.Background(), "1", entity)
	// Act
	found, ok := store.FindOneByIndex(context.Background(), "email", "alice@example.com")
	_, notOk := store.FindOneByIndex(context.Background(), "email", "unknown@example.com")
	// Assert
	assert.That(t, "found ok", ok, true)
	assert.That(t, "found ID", (*found).ID, "1")
	assert.That(t, "not found ok", notOk, false)
}

func Test_IndexedAccess_With_UpdatedEmail_Should_UpdateIndex(t *testing.T) {
	// Arrange
	base := NewInMemoryAccess[string, *testEntity]()
	store := NewIndexedAccess(base)
	store.AddIndex("email", func(e *testEntity) string { return e.Email })
	entity := &testEntity{ID: "1", Email: "old@example.com", Role: "admin"}
	_ = store.Create(context.Background(), "1", entity)
	// Act
	updated := &testEntity{ID: "1", Email: "new@example.com", Role: "admin"}
	err := store.Update(context.Background(), "1", updated)
	// Assert
	assert.That(t, "error", err, nil)
	oldResult, _ := store.FindByIndex(context.Background(), "email", "old@example.com")
	newResult, _ := store.FindByIndex(context.Background(), "email", "new@example.com")
	assert.That(t, "old email removed", len(oldResult), 0)
	assert.That(t, "new email indexed", len(newResult), 1)
}

func Test_IndexedAccess_With_DeletedEntity_Should_RemoveFromIndex(t *testing.T) {
	// Arrange
	base := NewInMemoryAccess[string, *testEntity]()
	store := NewIndexedAccess(base)
	store.AddIndex("email", func(e *testEntity) string { return e.Email })
	entity := &testEntity{ID: "1", Email: "test@example.com", Role: "admin"}
	_ = store.Create(context.Background(), "1", entity)
	// Act
	err := store.Delete(context.Background(), "1")
	// Assert
	assert.That(t, "error", err, nil)
	result, _ := store.FindByIndex(context.Background(), "email", "test@example.com")
	assert.That(t, "index removed", len(result), 0)
}

func Test_IndexedAccess_With_ExistingKey_Should_Read(t *testing.T) {
	// Arrange
	base := NewInMemoryAccess[string, *testEntity]()
	store := NewIndexedAccess(base)
	entity := &testEntity{ID: "1", Email: "test@example.com", Role: "admin"}
	_ = store.Create(context.Background(), "1", entity)
	// Act
	result, err := store.Read(context.Background(), "1")
	// Assert
	assert.That(t, "error", err, nil)
	assert.That(t, "result ID", (*result).ID, "1")
}

func Test_IndexedAccess_With_MultipleEntities_Should_ReadAll(t *testing.T) {
	// Arrange
	base := NewInMemoryAccess[string, *testEntity]()
	store := NewIndexedAccess(base)
	_ = store.Create(context.Background(), "1", &testEntity{ID: "1"})
	_ = store.Create(context.Background(), "2", &testEntity{ID: "2"})
	// Act
	result, err := store.ReadAll(context.Background())
	// Assert
	assert.That(t, "error", err, nil)
	assert.That(t, "count", len(result), 2)
}

func Test_IndexedAccess_With_UnknownIndex_Should_ReturnEmpty(t *testing.T) {
	// Arrange
	base := NewInMemoryAccess[string, *testEntity]()
	store := NewIndexedAccess(base)
	// Act
	result, err := store.FindByIndex(context.Background(), "unknown", "value")
	// Assert
	assert.That(t, "error", err, nil)
	assert.That(t, "result length", len(result), 0)
}

func Test_IndexedAccess_With_EmptyIndexKey_Should_NotIndex(t *testing.T) {
	// Arrange
	base := NewInMemoryAccess[string, *testEntity]()
	store := NewIndexedAccess(base)
	store.AddIndex("email", func(e *testEntity) string { return e.Email })
	entity := &testEntity{ID: "1", Email: "", Role: "admin"} // Empty email
	// Act
	err := store.Create(context.Background(), "1", entity)
	// Assert
	assert.That(t, "error", err, nil)
	result, _ := store.FindByIndex(context.Background(), "email", "")
	assert.That(t, "empty key not indexed", len(result), 0)
}
