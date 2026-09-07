package security_test

import (
	"fmt"

	"github.com/andygeiss/cloud-native-utils/security"
)

func ExampleEncrypt() {
	key := security.GenerateKey()

	// Each Encrypt uses a fresh nonce, so the ciphertext differs every time.
	// Only the round trip is worth asserting.
	ciphertext := security.Encrypt([]byte("secret"), key)
	plaintext, err := security.Decrypt(ciphertext, key)

	fmt.Println(string(plaintext), err)
	// Output: secret <nil>
}

func ExamplePassword() {
	hashed, err := security.Password([]byte("correct horse battery staple"))
	fmt.Println(err)
	fmt.Println(security.IsPasswordValid(hashed, []byte("correct horse battery staple")))
	fmt.Println(security.IsPasswordValid(hashed, []byte("wrong")))
	// Output:
	// <nil>
	// true
	// false
}
