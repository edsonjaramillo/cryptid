// Package commands contains the commands for the CLI
package commands

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"os"

	"github.com/urfave/cli/v3"
	"golang.org/x/crypto/pbkdf2"
)

// EncryptCommand adds the encryption command to the CLI
var EncryptCommand = &cli.Command{
	Name:   "encrypt",
	Usage:  "AES file encryption",
	Action: encryptAction,
	Flags:  []cli.Flag{passwordFlag, inputFileFlag, outputFileFlag},
}

// DecryptCommand adds the decryption command to the CLI
var DecryptCommand = &cli.Command{
	Name:   "decrypt",
	Usage:  "AES file decryption",
	Action: decryptAction,
	Flags:  []cli.Flag{passwordFlag, inputFileFlag, outputFileFlag},
}

// Actions

func encryptAction(_ context.Context, cmd *cli.Command) error {
	password := cmd.String("password")
	inputFile := cmd.String("input")
	outputFile := cmd.String("output")

	if err := encryptFile(inputFile, outputFile, password); err != nil {
		return err
	}

	return nil
}

func decryptAction(_ context.Context, cmd *cli.Command) error {
	password := cmd.String("password")
	inputFile := cmd.String("input")
	outputFile := cmd.String("output")

	if err := decryptFile(inputFile, outputFile, password); err != nil {
		return err
	}

	return nil
}

// Flags
var passwordFlag = &cli.StringFlag{
	Name:     "password",
	Aliases:  []string{"p"},
	Usage:    "Enter password.",
	OnlyOnce: true,
	Required: true,
}

var inputFileFlag = &cli.StringFlag{
	Name:      "input",
	Aliases:   []string{"i"},
	Usage:     "Enter input file.",
	OnlyOnce:  true,
	TakesFile: true,
	Required:  true,
}

var outputFileFlag = &cli.StringFlag{
	Name:      "output",
	Aliases:   []string{"o"},
	Usage:     "Enter output file.",
	OnlyOnce:  true,
	TakesFile: true,
	Required:  true,
}

const (
	saltSize = 16
	keySize  = 32
)

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 4096, keySize, sha256.New)
}

func encryptFile(inputFile, outputFile, password string) error {
	// Read the plaintext file
	plaintext, err := os.ReadFile(inputFile)
	if err != nil {
		println("Could not read file")
		return err
	}

	// Generate a random salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	// Derive a key from the password
	key := deriveKey(password, salt)

	// Create a new AES cipher using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Use GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Generate a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	// Encrypt the data
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	// Prepend salt and nonce to the ciphertext
	finalData := bytes.Join([][]byte{salt, nonce, ciphertext}, nil)

	// Write the encrypted data to the output file
	return os.WriteFile(outputFile, finalData, 0644)
}

func decryptFile(inputFile, outputFile, password string) error {
	// Read the encrypted file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		println("Could not read file")
		return err
	}

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
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		println("Authentication failed")
		return err
	}

	// Write the plaintext to the output file
	return os.WriteFile(outputFile, plaintext, 0644)
}
