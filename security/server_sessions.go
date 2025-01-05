package security

import (
	"encoding/hex"
	"sync"
)

// ServerSession is a session for a user.
type ServerSession struct {
	AvatarURL string `json:"avatar_url"`
	EMail     string `json:"email"`
	Login     string `json:"login"`
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

// Update adds a new session to the serverSessions.
func (a *ServerSessions) Update(info ServerSession) (sessionID string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	bytes := GenerateKey()
	sessionID = hex.EncodeToString(bytes[:])
	a.sessions[sessionID] = info
	return sessionID
}

// Get returns the session for the given sessionID.
func (a *ServerSessions) Get(sessionID string) (*ServerSession, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	info, ok := a.sessions[sessionID]
	return &info, ok
}
