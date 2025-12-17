package security_test

import (
	"reflect"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/security"
)

func Test_Hash_With_DifferentData_Should_ReturnDifferentHashes(t *testing.T) {
	// Arrange
	tag := "secure-tag"
	data1 := []byte("Data 1")
	data2 := []byte("Data 2")

	// Act
	hash1 := security.Hash(tag, data1)
	hash2 := security.Hash(tag, data2)

	// Assert
	assert.That(t, "hashes must be different", !reflect.DeepEqual(hash1, hash2), true)
}

func Test_Hash_With_DifferentTags_Should_ReturnDifferentHashes(t *testing.T) {
	// Arrange
	data := []byte("Test data")
	tag1 := "tag1"
	tag2 := "tag2"

	// Act
	hash1 := security.Hash(tag1, data)
	hash2 := security.Hash(tag2, data)

	// Assert
	assert.That(t, "hashes must be different", !reflect.DeepEqual(hash1, hash2), true)
}

func Test_Hash_With_EmptyData_Should_ReturnNonEmptyHash(t *testing.T) {
	// Arrange
	tag := "secure-tag"
	data := []byte("")

	// Act
	hash := security.Hash(tag, data)

	// Assert
	assert.That(t, "hash length must be not equal 0", len(hash) > 0, true)
}

func Test_Hash_With_SameInput_Should_ReturnSameHash(t *testing.T) {
	// Arrange
	tag := "secure-tag"
	data := []byte("Test data")

	// Act
	hash1 := security.Hash(tag, data)
	hash2 := security.Hash(tag, data)

	// Assert
	assert.That(t, "hashes must be equal", hash1, hash2)
}
