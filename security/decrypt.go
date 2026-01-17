package security

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// Decrypt takes an encrypted byte slice (ciphertext) and a 256-bit AES key,
// and decrypts the ciphertext using AES-GCM.
func Decrypt(ciphertext []byte, key [32]byte) ([]byte, error) {
	// Create a new AES cipher block with the provided key.
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	// Create a new GCM (Galois/Counter Mode) cipher based on the AES block.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Ensure the ciphertext is long enough to include the nonce.
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	// Decrypt the ciphertext using the nonce and the encrypted message.
	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()], // Extract the nonce from the ciphertext.
		ciphertext[gcm.NonceSize():], // The encrypted message follows the nonce.
		nil,
	)
}
