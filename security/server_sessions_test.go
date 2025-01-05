package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestServerSessions_Update(t *testing.T) {
	sessions := security.NewServerSessions()
	token := sessions.Update(security.ServerSession{AvatarURL: "avatar_url", EMail: "email", Login: "login", Name: "name"})
	assert.That(t, "token is correct", len(token), 64)
}

func TestServerSessions_Get(t *testing.T) {
	sessions := security.NewServerSessions()
	session := security.ServerSession{AvatarURL: "avatar_url", EMail: "email", Login: "login", Name: "name"}
	token := sessions.Update(session)
	current, found := sessions.Get(token)
	assert.That(t, "session must be found", found, true)
	assert.That(t, "session is correct", current, session)
}
