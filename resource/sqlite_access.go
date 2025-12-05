package resource

import (
	"database/sql"
	"encoding/json"
)

// sqliteAccess provides a simple key-value store using SQLite.
type sqliteAccess[K comparable, V any] struct {
	db *sql.DB
}

// NewSqliteAccess creates a new instance of sqliteAccess.
func NewSqliteAccess[K comparable, V any](db *sql.DB) *sqliteAccess[K, V] {
	return &sqliteAccess[K, V]{
		db: db,
	}
}

// Create inserts a new key-value pair into the table.
func (a *sqliteAccess[K, V]) Create(key K, value V) error {
	encoded, _ := json.Marshal(value)
	valueAsString := string(encoded)
	_, err := a.db.Exec("INSERT INTO kv_store (key, value) VALUES (?, ?)", key, valueAsString)
	return err
}

// Init initializes the table and index.
func (a *sqliteAccess[K, V]) Init() error {
	_, _ = a.db.Exec("DROP TABLE IF EXISTS kv_store;")
	_, _ = a.db.Exec("CREATE TABLE IF NOT EXISTS kv_store (key TEXT PRIMARY KEY, value TEXT);")
	_, _ = a.db.Exec("CREATE INDEX IF NOT EXISTS idx_kv_store_key ON kv_store (key);")
	return nil
}

// Read returns the value associated with the given key.
func (a *sqliteAccess[K, V]) Read(key K) (V, error) {
	var value V
	var valueAsString string
	err := a.db.QueryRow("SELECT value FROM kv_store WHERE key = ?", key).Scan(&valueAsString)
	if err != nil {
		return value, err
	}
	err = json.Unmarshal([]byte(valueAsString), &value)
	return value, err
}

// ReadAll returns all values from the table.
func (a *sqliteAccess[K, V]) ReadAll() ([]V, error) {
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
	valueAsString, _ := json.Marshal(value)
	_, err := a.db.Exec("UPDATE kv_store SET value = ? WHERE key = ?", valueAsString, key)
	return err
}

// Delete removes the key-value pair associated with the given key.
func (a *sqliteAccess[K, V]) Delete(key K) error {
	_, err := a.db.Exec("DELETE FROM kv_store WHERE key = ?", key)
	return err
}
