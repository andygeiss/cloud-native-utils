package resource

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"
)

// SqliteAccess provides a simple key-value store using SQLite.
type SqliteAccess[K comparable, V any] struct {
	db    *sql.DB
	mutex sync.RWMutex
}

// NewSqliteAccess creates a new instance of SqliteAccess.
func NewSqliteAccess[K comparable, V any](db *sql.DB) *SqliteAccess[K, V] {
	return &SqliteAccess[K, V]{
		db: db,
	}
}

// Create inserts a new key-value pair into the table.
func (a *SqliteAccess[K, V]) Create(ctx context.Context, key K, value V) (err error) {
	// Skip if context is canceled or timed out.
	if err := ctx.Err(); err != nil {
		return err
	}

	// Ensure that the table is not modified concurrently.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Encode the value and insert it into the table.
	encoded, _ := json.Marshal(value)
	valueAsString := string(encoded)

	// Ensure that the value is inserted atomically by using a transaction.
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO kv_store (key, value) VALUES (?, ?)", key, valueAsString)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Delete removes the key-value pair associated with the given key.
func (a *SqliteAccess[K, V]) Delete(ctx context.Context, key K) (err error) {
	// Skip if context is canceled or timed out.
	if err := ctx.Err(); err != nil {
		return err
	}

	// Ensure that the table is not modified concurrently.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Ensure that the value is deleted atomically by using a transaction.
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM kv_store WHERE key = ?", key)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Init initializes the table and index.
func (a *SqliteAccess[K, V]) Init(ctx context.Context) (err error) {
	// Skip if context is canceled or timed out.
	if err := ctx.Err(); err != nil {
		return err
	}

	// Ensure that the table is not modified concurrently.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Drop the table if it exists.
	_, err = a.db.ExecContext(ctx, "DROP TABLE IF EXISTS kv_store;")
	if err != nil {
		return err
	}

	// Create the table.
	_, err = a.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS kv_store (key TEXT PRIMARY KEY, value TEXT);")
	if err != nil {
		return err
	}

	// Create the index.
	_, err = a.db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_kv_store_key ON kv_store (key);")
	if err != nil {
		return err
	}

	return nil
}

// Read returns the value associated with the given key.
func (a *SqliteAccess[K, V]) Read(ctx context.Context, key K) (ptr *V, err error) {
	// Skip if context is canceled or timed out.
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Ensure that read operations can be performed concurrently.
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// Query the value from the table.
	var valueAsString string
	err = a.db.QueryRowContext(ctx, "SELECT value FROM kv_store WHERE key = ?", key).Scan(&valueAsString)
	if err != nil {
		return nil, err
	}

	// Unmarshal the value from JSON.
	var value V
	err = json.Unmarshal([]byte(valueAsString), &value)
	return &value, err
}

// ReadAll returns all values from the table.
func (a *SqliteAccess[K, V]) ReadAll(ctx context.Context) (values []V, err error) {
	// Skip if context is canceled or timed out.
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Ensure that read operations can be performed concurrently.
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// Query all values from the table.
	rows, err := a.db.QueryContext(ctx, "SELECT value FROM kv_store")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Store all values in a slice.
	for rows.Next() {
		var valueAsString string
		if err := rows.Scan(&valueAsString); err != nil {
			return nil, err
		}
		var value V
		err = json.Unmarshal([]byte(valueAsString), &value)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, rows.Err()
}

// Update updates the value associated with the given key.
func (a *SqliteAccess[K, V]) Update(ctx context.Context, key K, value V) (err error) {
	// Skip if context is canceled or timed out.
	if err := ctx.Err(); err != nil {
		return err
	}

	// Ensure that the table is not modified concurrently.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Encode the value as JSON.
	valueAsString, _ := json.Marshal(value)

	// Ensure that the value is updated atomically by using a transaction.
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE kv_store SET value = ? WHERE key = ?", valueAsString, key)
	if err != nil {
		return err
	}

	return tx.Commit()
}
