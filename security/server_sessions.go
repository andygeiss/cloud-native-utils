package security

import (
	"encoding/hex"
	"sync"
)

// ServerSession is a session for a user.
type ServerSession struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// ServerSessions is a thread-safe map of email addresses to tokens.
type ServerSessions struct {
	sessions map[string]ServerSession
	mutex    sync.RWMutex
}

// NewServerSessions creates a new serverSessions.
func NewServerSessions() *ServerSessions {
	return &ServerSessions{
		sessions: make(map[string]ServerSession),
	}
}

// Create adds a new session to the serverSessions.
func (a *ServerSessions) Create() (s ServerSession) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	bytes := GenerateKey()
	sessionID := hex.EncodeToString(bytes[:])
	s.ID = sessionID
	a.sessions[sessionID] = s
	return s
}

// Get returns the session for the given sessionID.
func (a *ServerSessions) Read(id string) (*ServerSession, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	info, ok := a.sessions[id]
	return &info, ok
}

// Update adds a new session to the serverSessions.
func (a *ServerSessions) Update(s ServerSession) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	session := a.sessions[s.ID]
	session.Value = s.Value
	a.sessions[s.ID] = session
	return
}

// Delete removes the session with the given sessionID.
func (a *ServerSessions) Delete(id string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	delete(a.sessions, id)
}
