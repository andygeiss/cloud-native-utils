package i18n

import "sync"

// SessionLanguageStore stores language preferences per session ID.
type SessionLanguageStore struct {
	languages map[string]string
	mu        sync.RWMutex
}

// NewSessionLanguageStore creates a new session language store.
func NewSessionLanguageStore() *SessionLanguageStore {
	return &SessionLanguageStore{
		languages: make(map[string]string),
	}
}

// Get returns the language for the given session ID.
// Returns empty string if not found.
func (a *SessionLanguageStore) Get(sessionID string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.languages[sessionID]
}

// Set sets the language for the given session ID.
func (a *SessionLanguageStore) Set(sessionID, lang string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.languages[sessionID] = lang
}

// Clear removes the language preference for the given session ID.
func (a *SessionLanguageStore) Clear(sessionID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.languages, sessionID)
}
