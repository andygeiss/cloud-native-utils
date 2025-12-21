package i18n

import "sync"

// SessionStore stores values of any type per session ID.
// This is a generic version of SessionLanguageStore that can store any type.
type SessionStore[V any] struct {
	data map[string]V
	mu   sync.RWMutex
}

// NewSessionStore creates a new generic session store.
func NewSessionStore[V any]() *SessionStore[V] {
	return &SessionStore[V]{
		data: make(map[string]V),
	}
}

// Get returns the value for the given session ID.
// Returns the zero value and false if not found.
func (a *SessionStore[V]) Get(sessionID string) (V, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	val, ok := a.data[sessionID]
	return val, ok
}

// Set sets the value for the given session ID.
func (a *SessionStore[V]) Set(sessionID string, value V) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.data[sessionID] = value
}

// Clear removes the value for the given session ID.
func (a *SessionStore[V]) Clear(sessionID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.data, sessionID)
}
