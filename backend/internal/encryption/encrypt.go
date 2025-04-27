// Package encryption provides functions to encrypt and decrypt data using AES-GCM.
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize         = 16
	keySize          = 32
	pbkdf2Iterations = 600_000
	gcmNonceSize     = 12
)

// ErrInvalidCiphertext is returned when the ciphertext is invalid or too short.
var ErrInvalidCiphertext = errors.New("encryption: invalid ciphertext")

// ErrDecryptionFailed is returned when decryption fails due to an authentication error.
var ErrDecryptionFailed = errors.New("encryption: decryption failed (authentication error)")

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, keySize, sha256.New)
}

// EncryptData encrypts the input data using AES-GCM with a key derived from the password.
func EncryptData(plainData []byte, password string) ([]byte, error) {
	// 1. Generate a random salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// 2. Derive a key from the password and salt
	key := deriveKey(password, salt)

	// 3. Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		// Avoid leaking key material in error messages
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// 4. Create GCM AEAD
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	// 5. Generate a random nonce
	nonceSize := aesGCM.NonceSize()
	if nonceSize != gcmNonceSize {
		return nil, fmt.Errorf("unexpected GCM nonce size: %d", nonceSize)
	}
	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 6. Encrypt the data using GCM Seal
	ciphertext := aesGCM.Seal(nil, nonce, plainData, nil)

	// 7. Assemble the final result
	encryptedData := make([]byte, 0, saltSize+nonceSize+len(ciphertext))
	encryptedData = append(encryptedData, salt...)
	encryptedData = append(encryptedData, nonce...)
	encryptedData = append(encryptedData, ciphertext...)

	return encryptedData, nil
}

// DecryptData decrypts data previously encrypted with EncryptData.
func DecryptData(encryptedData []byte, password string) ([]byte, error) {
	minLen := saltSize + gcmNonceSize + 1
	if len(encryptedData) < minLen {
		return nil, fmt.Errorf("%w: input data too short (%d bytes)", ErrInvalidCiphertext, len(encryptedData))
	}

	// 2. Extract salt, nonce, and ciphertext
	salt := encryptedData[:saltSize]
	nonce := encryptedData[saltSize : saltSize+gcmNonceSize]
	ciphertext := encryptedData[saltSize+gcmNonceSize:]

	// 3. Derive the key using the same parameters
	key := deriveKey(password, salt)

	// 4. Create the AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// 5. Create GCM AEAD
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	// 6. Check if the extracted nonce size matches the expected GCM nonce size.
	if len(nonce) != aesGCM.NonceSize() {
		return nil, fmt.Errorf("%w: incorrect nonce size (%d bytes)", ErrInvalidCiphertext, len(nonce))
	}

	// 7. Decrypt the data using GCM Open
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecryptionFailed, err)
	}

	// 8. Return the plaintext
	return plaintext, nil
}
