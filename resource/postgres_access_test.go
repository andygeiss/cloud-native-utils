package resource_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/resource"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func getPostgresDSN() string {
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		return ""
	}
	return dsn
}

func Test_PostgresAccess_With_CreateDuplicateKey_Should_ReturnError(t *testing.T) {
	// Arrange
	dsn := getPostgresDSN()
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN not set, skipping PostgreSQL tests")
	}
	db, _ := sql.Open("pgx", dsn)
	defer func() { _ = db.Close() }()
	a := resource.NewPostgresAccess[string, string](db)
	ctx := context.Background()
	_ = a.Init(ctx)
	_ = a.Create(ctx, "key", "value")

	// Act
	err := a.Create(ctx, "key", "value")

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
}

func Test_PostgresAccess_With_CreateValidKey_Should_Succeed(t *testing.T) {
	// Arrange
	dsn := getPostgresDSN()
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN not set, skipping PostgreSQL tests")
	}
	db, _ := sql.Open("pgx", dsn)
	defer func() { _ = db.Close() }()
	a := resource.NewPostgresAccess[string, string](db)
	ctx := context.Background()

	// Act
	err := a.Init(ctx)
	err2 := a.Create(ctx, "key", "value")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "err2 must be nil", err2, nil)
}

func Test_PostgresAccess_With_DeleteValidKey_Should_RemoveValue(t *testing.T) {
	// Arrange
	dsn := getPostgresDSN()
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN not set, skipping PostgreSQL tests")
	}
	db, _ := sql.Open("pgx", dsn)
	defer func() { _ = db.Close() }()
	a := resource.NewPostgresAccess[string, string](db)
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

func Test_PostgresAccess_With_ReadAllMultipleKeys_Should_ReturnAllValues(t *testing.T) {
	// Arrange
	dsn := getPostgresDSN()
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN not set, skipping PostgreSQL tests")
	}
	db, _ := sql.Open("pgx", dsn)
	defer func() { _ = db.Close() }()
	a := resource.NewPostgresAccess[string, string](db)
	ctx := context.Background()
	_ = a.Init(ctx)
	_ = a.Create(ctx, "key1", "value1")
	_ = a.Create(ctx, "key2", "value2")

	// Act
	values, err := a.ReadAll(ctx)

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "values must have 2 elements", len(values), 2)
}

func Test_PostgresAccess_With_ReadMissingKey_Should_ReturnError(t *testing.T) {
	// Arrange
	dsn := getPostgresDSN()
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN not set, skipping PostgreSQL tests")
	}
	db, _ := sql.Open("pgx", dsn)
	defer func() { _ = db.Close() }()
	a := resource.NewPostgresAccess[string, string](db)
	ctx := context.Background()
	_ = a.Init(ctx)
	_ = a.Create(ctx, "key", "value")

	// Act
	_, err := a.Read(ctx, "key2")

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
}

func Test_PostgresAccess_With_ReadValidKey_Should_ReturnValue(t *testing.T) {
	// Arrange
	dsn := getPostgresDSN()
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN not set, skipping PostgreSQL tests")
	}
	db, _ := sql.Open("pgx", dsn)
	defer func() { _ = db.Close() }()
	a := resource.NewPostgresAccess[string, string](db)
	ctx := context.Background()
	_ = a.Init(ctx)
	_ = a.Create(ctx, "key", "value")

	// Act
	value, err := a.Read(ctx, "key")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be 'value'", *value, "value")
}

func Test_PostgresAccess_With_UpdateValidKey_Should_UpdateValue(t *testing.T) {
	// Arrange
	dsn := getPostgresDSN()
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN not set, skipping PostgreSQL tests")
	}
	db, _ := sql.Open("pgx", dsn)
	defer func() { _ = db.Close() }()
	a := resource.NewPostgresAccess[string, string](db)
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
