// Package commands contains the commands for the CLI
package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/edsonjaramillo/hyde/backend/internal/aes"
	"github.com/urfave/cli/v3"
)

// EncryptCommand adds the encryption command to the CLI
var EncryptCommand = &cli.Command{
	Name:      "encrypt",
	Usage:     "AES file encryption",
	UsageText: "hyde encrypt <file> --password <password>\n",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "input",
			UsageText: "Enter input file.",
		},
	},
	Action: encryptAction,
	Flags:  []cli.Flag{passwordFlag, deleteFlag},
}

// DecryptCommand adds the decryption command to the CLI
var DecryptCommand = &cli.Command{
	Name:      "decrypt",
	Usage:     "AES file decryption",
	UsageText: "hyde decrypt <file> --password <password>\n",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:      "input",
			UsageText: "Enter input file.",
		},
	},
	Action: decryptAction,
	Flags:  []cli.Flag{passwordFlag, deleteFlag},
}

// Actions

func encryptAction(_ context.Context, cmd *cli.Command) error {
	file := cmd.StringArg("input")
	filedata, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	password := cmd.String("password")
	encryptedData, err := aes.EncryptData(filedata, password)
	if err != nil {
		return err
	}

	outputFilename := AddEncExtension(file)
	if err := os.WriteFile(outputFilename, encryptedData, 0644); err != nil {
		return err
	}

	deleteEnabled := cmd.Bool("delete")
	if deleteEnabled {
		if err := os.Remove(file); err != nil {
			fmt.Printf("Error removing file: %v\n", err)
		}
	}

	return nil
}

func decryptAction(_ context.Context, cmd *cli.Command) error {
	encryptedFile := cmd.StringArg("input")
	filedata, err := os.ReadFile(encryptedFile)
	if err != nil {
		return err
	}

	password := cmd.String("password")
	plaintext, err := aes.DecryptData(filedata, password)
	if err != nil {
		return err
	}

	outputFilename := RemoveEncExtension(encryptedFile)
	if err := os.WriteFile(outputFilename, plaintext, 0644); err != nil {
		return err
	}

	deleteEnabled := cmd.Bool("delete")
	if deleteEnabled {
		if err := os.Remove(encryptedFile); err != nil {
			fmt.Printf("Error removing file: %v\n", err)
		}
	}

	return nil
}

// AddEncExtension adds the .enc extension to the file name if it doesn't already have it
func AddEncExtension(file string) string {
	base := filepath.Base(file)

	if len(base) < 4 || file[len(base)-4:] != ".enc" {
		return base + ".enc"
	}
	return base
}

// RemoveEncExtension removes the .enc extension from the file name if it has it
func RemoveEncExtension(file string) string {
	base := filepath.Base(file)

	if len(base) < 4 || base[len(base)-4:] != ".enc" {
		return base
	}
	return base[:len(base)-4]
}

// Flags

var passwordFlag = &cli.StringFlag{
	Name:     "password",
	Aliases:  []string{"p"},
	Usage:    "Enter password.",
	OnlyOnce: true,
	Required: true,
}

var deleteFlag = &cli.BoolFlag{
	Name:     "delete",
	Usage:    "Remove the original file after encryption.",
	OnlyOnce: true,
	Required: false,
}
