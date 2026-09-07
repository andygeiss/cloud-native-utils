package security_test

import (
	"bytes"
	"testing"

	"github.com/andygeiss/cloud-native-utils/security"
)

// fuzzKey is fixed so a crashing input stays reproducible across runs.
var fuzzKey = func() [32]byte {
	var key [32]byte
	for i := range key {
		key[i] = byte(i)
	}
	return key
}()

func FuzzDecrypt(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte("short"))
	f.Add(make([]byte, 12)) // exactly the nonce size, with no payload
	f.Add(security.Encrypt([]byte(""), fuzzKey))
	f.Add(security.Encrypt([]byte("hello"), fuzzKey))

	f.Fuzz(func(t *testing.T, ciphertext []byte) {
		plaintext, err := security.Decrypt(ciphertext, fuzzKey)
		if err != nil {
			return // Rejecting a forged or truncated ciphertext is success.
		}

		// Whatever authenticated must survive another round trip unchanged.
		again, err := security.Decrypt(security.Encrypt(plaintext, fuzzKey), fuzzKey)
		if err != nil {
			t.Fatalf("re-encrypted plaintext failed to decrypt: %v", err)
		}
		if !bytes.Equal(again, plaintext) {
			t.Errorf("round trip changed the plaintext: %q became %q", plaintext, again)
		}
	})
}
