package security_test

import (
	"cloud-native/utils/assert"
	"cloud-native/utils/security"
	"reflect"
	"testing"
)

func TestHash(t *testing.T) {
	tag := "secure-tag"
	data := []byte("Test data")
	hash1 := security.Hash(tag, data)
	hash2 := security.Hash(tag, data)
	assert.That(t, "hashes must be equal", hash1, hash2)
}

func TestHash_Different_Tags(t *testing.T) {
	data := []byte("Test data")
	tag1 := "tag1"
	tag2 := "tag2"
	hash1 := security.Hash(tag1, data)
	hash2 := security.Hash(tag2, data)
	assert.That(t, "hashes must be different", !reflect.DeepEqual(hash1, hash2), true)
}

func TestHash_Different_Data(t *testing.T) {
	tag := "secure-tag"
	data1 := []byte("Data 1")
	data2 := []byte("Data 2")
	hash1 := security.Hash(tag, data1)
	hash2 := security.Hash(tag, data2)
	assert.That(t, "hashes must be different", !reflect.DeepEqual(hash1, hash2), true)
}

func TestHash_Empty_Data(t *testing.T) {
	tag := "secure-tag"
	data := []byte("")
	hash := security.Hash(tag, data)
	assert.That(t, "hash length must be not equal 0", len(hash) > 0, true)
}
