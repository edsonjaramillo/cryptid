// Package commands contains the commands for the CLI
package commands

import (
	"context"
	"os"

	"github.com/edsonjaramillo/hyde/backend/internal/encryption"
	"github.com/urfave/cli/v3"
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

	inputFileStream, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	encryptedData, err := encryption.EncryptFile(inputFileStream, password)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputFile, encryptedData, 0644); err != nil {
		return err
	}

	return nil
}

func decryptAction(_ context.Context, cmd *cli.Command) error {
	password := cmd.String("password")
	inputFile := cmd.String("input")
	outputFile := cmd.String("output")

	inputFileStream, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	plaintext, err := encryption.DecryptFile(inputFileStream, password)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputFile, plaintext, 0644); err != nil {
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
