package security

import (
	"encoding/hex"
	"sync"
)

// ServerSession is a session for a user.
type ServerSession struct {
	ID        string `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
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
func (a *ServerSessions) Create() (session ServerSession) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	bytes := GenerateKey()
	sessionID := hex.EncodeToString(bytes[:])
	session.ID = sessionID
	a.sessions[sessionID] = session
	return session
}

// Get returns the session for the given sessionID.
func (a *ServerSessions) Read(id string) (*ServerSession, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	info, ok := a.sessions[id]
	return &info, ok
}

// Update adds a new session to the serverSessions.
func (a *ServerSessions) Update(info ServerSession) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	session := a.sessions[info.ID]
	session.AvatarURL = info.AvatarURL
	session.Name = info.Name
	a.sessions[info.ID] = session
	return
}

// Delete removes the session with the given sessionID.
func (a *ServerSessions) Delete(id string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	delete(a.sessions, id)
}
