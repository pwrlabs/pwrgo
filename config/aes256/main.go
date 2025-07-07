package aes256

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iterationCount = 65536
	keyLength      = 32 // 256 bits
	salt           = "your-salt-value"
)

func generateKey(password string) []byte {
	return pbkdf2.Key([]byte(password), []byte(salt), iterationCount, keyLength, sha256.New)
}

func generateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return iv, nil
}

func Encrypt(data []byte, password string) ([]byte, error) {
	key := generateKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv, err := generateIV()
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	padding := blockSize - len(data)%blockSize
	paddedData := make([]byte, len(data)+padding)
	copy(paddedData, data)
	for i := len(data); i < len(paddedData); i++ {
		paddedData[i] = byte(padding)
	}

	ciphertext := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	result := make([]byte, len(iv)+len(ciphertext))
	copy(result, iv)
	copy(result[len(iv):], ciphertext)

	return result, nil
}

func Decrypt(encryptedData []byte, password string) ([]byte, error) {
	key := generateKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := encryptedData[:aes.BlockSize]
	ciphertext := encryptedData[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	padding := int(plaintext[len(plaintext)-1])
	return plaintext[:len(plaintext)-padding], nil
}
