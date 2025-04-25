package encryption

import (
	"bytes"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize = 16
	keySize  = 32
)

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 4096, keySize, sha256.New)
}

func EncryptFile(inputFileStream []byte, password string) ([]byte, error) {
	// Generate a random salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// Derive a key from the password
	key := deriveKey(password, salt)

	// Create a new AES cipher using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Use GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt the data
	ciphertext := aesGCM.Seal(nil, nonce, inputFileStream, nil)

	// Prepend salt and nonce to the ciphertext
	finalData := bytes.Join([][]byte{salt, nonce, ciphertext}, nil)

	// return the encrypted data
	return finalData, nil
}

func DecryptFile(inputFileStream []byte, password string) ([]byte, error) {
	// Read the encrypted file
	data := inputFileStream

	// Extract salt, nonce, and ciphertext
	salt := data[:saltSize]
	nonceSize := 12 // GCM standard nonce size
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	// Derive the key using the same method
	key := deriveKey(password, salt)

	// Create the cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		println("Authentication failed")
		return nil, err
	}

	// Write the plaintext to the output file
	return plaintext, nil
}
