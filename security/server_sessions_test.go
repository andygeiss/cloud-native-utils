package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestServerSessions_Update(t *testing.T) {
	sessions := security.NewServerSessions()
	id := sessions.Update(security.ServerSession{AvatarURL: "avatar_url", Name: "name"})
	assert.That(t, "id is correct", len(id), 64)
}

func TestServerSessions_Get(t *testing.T) {
	sessions := security.NewServerSessions()
	session := security.ServerSession{AvatarURL: "avatar_url", Name: "name"}
	id := sessions.Update(session)
	session.ID = id
	current, found := sessions.Get(id)
	assert.That(t, "session must be found", found, true)
	assert.That(t, "session is correct", *current, session)
}

func TestServerSessions_Remove(t *testing.T) {
	sessions := security.NewServerSessions()
	id := sessions.Update(security.ServerSession{AvatarURL: "avatar_url", Name: "name"})
	sessions.Remove(id)
	_, found := sessions.Get(id)
	assert.That(t, "session must not be found", found, false)
}
