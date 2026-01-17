package web_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/web"
)

func Test_ServerSessions_With_Create_Should_ReturnSessionWithCorrectID(t *testing.T) {
	// Arrange
	sessions := web.NewServerSessions()

	// Act
	session := sessions.Create("the unique id", nil)

	// Assert
	assert.That(t, "id is correct", session.ID, "the unique id")
}

func Test_ServerSessions_With_Delete_Should_RemoveSession(t *testing.T) {
	// Arrange
	sessions := web.NewServerSessions()
	session := sessions.Create("the unique id", nil)

	// Act
	sessions.Delete(session.ID)
	_, found := sessions.Read(session.ID)

	// Assert
	assert.That(t, "session must not be found", found, false)
}

func Test_ServerSessions_With_Read_Should_ReturnExistingSession(t *testing.T) {
	// Arrange
	sessions := web.NewServerSessions()
	session := sessions.Create("the unique id", nil)

	// Act
	current, found := sessions.Read(session.ID)

	// Assert
	assert.That(t, "session must be found", found, true)
	assert.That(t, "session is correct", *current, session)
}

func Test_ServerSessions_With_Update_Should_ModifySessionData(t *testing.T) {
	// Arrange
	sessions := web.NewServerSessions()
	session := sessions.Create("the unique id", nil)
	session.Data = "value"

	// Act
	sessions.Update(session)
	current, found := sessions.Read(session.ID)

	// Assert
	assert.That(t, "session must be found", found, true)
	assert.That(t, "session is correct", *current, session)
}
