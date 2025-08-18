package security_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestServerSessions_Create(t *testing.T) {
	sessions := security.NewServerSessions()
	session := sessions.Create("the unique id", nil)
	id := session.ID
	assert.That(t, "id is correct", id, "the unique id")
}

func TestServerSessions_Read(t *testing.T) {
	sessions := security.NewServerSessions()
	session := sessions.Create("the unique id", nil)
	id := session.ID
	current, found := sessions.Read(id)
	assert.That(t, "session must be found", found, true)
	assert.That(t, "session is correct", *current, session)
}

func TestServerSessions_Update(t *testing.T) {
	sessions := security.NewServerSessions()
	session := sessions.Create("the unique id", nil)
	id := session.ID
	session.Data = "value"
	sessions.Update(session)
	current, found := sessions.Read(id)
	assert.That(t, "session must be found", found, true)
	assert.That(t, "session is correct", *current, session)
}

func TestServerSessions_Delete(t *testing.T) {
	sessions := security.NewServerSessions()
	session := sessions.Create("the unique id", nil)
	id := session.ID
	sessions.Delete(id)
	_, found := sessions.Read(id)
	assert.That(t, "session must not be found", found, false)
}
