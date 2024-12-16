package security_test

import (
	"bytes"
	"cloud-native/utils/security"
	"testing"
)

func TestHash(t *testing.T) {
	tag := "secure-tag"
	data := []byte("Test data")
	hash1 := security.Hash(tag, data)
	hash2 := security.Hash(tag, data)
	// Verify both hashes are equal.
	if !bytes.Equal(hash1, hash2) {
		t.Errorf("hashes must be equal, but got %x and %x", hash1, hash2)
	}
}

func TestHashDifferentTags(t *testing.T) {
	data := []byte("Test data")
	tag1 := "tag1"
	tag2 := "tag2"
	// Generate hashes with different tags.
	hash1 := security.Hash(tag1, data)
	hash2 := security.Hash(tag2, data)
	// Verify the hashes are different.
	if bytes.Equal(hash1, hash2) {
		t.Errorf("hashes must be different")
	}
}

func TestHash_Different_Data(t *testing.T) {
	tag := "secure-tag"
	data1 := []byte("Data 1")
	data2 := []byte("Data 2")
	// Generate hashes with different data.
	hash1 := security.Hash(tag, data1)
	hash2 := security.Hash(tag, data2)
	// Verify the hashes are different.
	if bytes.Equal(hash1, hash2) {
		t.Errorf("hashes must be different")
	}
}

func TestHash_Empty_Data(t *testing.T) {
	tag := "secure-tag"
	data := []byte("")
	// Generate hash for empty data.
	hash := security.Hash(tag, data)
	// Ensure hash is not nil.
	if len(hash) == 0 {
		t.Error("hash length must be not equal 0")
	}
}
