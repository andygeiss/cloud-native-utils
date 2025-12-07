package resource

import (
	"database/sql"
	"encoding/json"
	"sync"
)

// sqliteAccess provides a simple key-value store using SQLite.
type sqliteAccess[K comparable, V any] struct {
	db    *sql.DB
	mutex sync.RWMutex
}

// NewSqliteAccess creates a new instance of sqliteAccess.
func NewSqliteAccess[K comparable, V any](db *sql.DB) *sqliteAccess[K, V] {
	return &sqliteAccess[K, V]{
		db: db,
	}
}

// Create inserts a new key-value pair into the table.
func (a *sqliteAccess[K, V]) Create(key K, value V) error {
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

// Init initializes the table and index.
func (a *sqliteAccess[K, V]) Init() error {
	// Ensure that the table is not modified concurrently.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Drop the table if it exists, create a new one and an index.
	_, _ = a.db.Exec("DROP TABLE IF EXISTS kv_store;")
	_, _ = a.db.Exec("CREATE TABLE IF NOT EXISTS kv_store (key TEXT PRIMARY KEY, value TEXT);")
	_, _ = a.db.Exec("CREATE INDEX IF NOT EXISTS idx_kv_store_key ON kv_store (key);")
	return nil
}

// Read returns the value associated with the given key.
func (a *sqliteAccess[K, V]) Read(key K) (*V, error) {
	// Ensure that read operations can be performed concurrently.
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// Query the value from the table.
	var value V
	var valueAsString string
	err := a.db.QueryRow("SELECT value FROM kv_store WHERE key = ?", key).Scan(&valueAsString)
	if err != nil {
		return &value, err
	}

	// Unmarshal the value from JSON.
	err = json.Unmarshal([]byte(valueAsString), &value)
	return &value, err
}

// ReadAll returns all values from the table.
func (a *sqliteAccess[K, V]) ReadAll() ([]V, error) {
	// Ensure that read operations can be performed concurrently.
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// Query all values from the table.
	rows, err := a.db.Query("SELECT value FROM kv_store")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Store all values in a slice.
	var values []V
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
func (a *sqliteAccess[K, V]) Update(key K, value V) error {
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

	_, err = tx.Exec("UPDATE kv_store SET value = ? WHERE key = ?", valueAsString, key)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Delete removes the key-value pair associated with the given key.
func (a *sqliteAccess[K, V]) Delete(key K) error {
	// Ensure that the table is not modified concurrently.
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Ensure that the value is deleted atomically by using a transaction.
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM kv_store WHERE key = ?", key)
	if err != nil {
		return err
	}

	return tx.Commit()
}
