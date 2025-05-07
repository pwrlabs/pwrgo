package encode

import (
	"crypto/sha512"

	"golang.org/x/crypto/pbkdf2"
)

// GenerateSeed generates a seed from a mnemonic phrase using PBKDF2
func GenerateSeed(mnemonic []byte, passphrase string) []byte {
	// Use PBKDF2 with SHA512 to generate the seed
	// Parameters:
	// - mnemonic as the password
	// - "mnemonic" + passphrase as the salt
	// - 2048 iterations
	// - 64 bytes output (512 bits)
	salt := append([]byte("mnemonic"), []byte(passphrase)...)
	return pbkdf2.Key(mnemonic, salt, 2048, 64, sha512.New)
}
