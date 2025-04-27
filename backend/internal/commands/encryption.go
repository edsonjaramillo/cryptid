// Package commands contains the commands for the CLI
package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/edsonjaramillo/hyde/backend/internal/encryption"
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
	encryptedData, err := encryption.EncryptData(filedata, password)
	if err != nil {
		return err
	}

	outputFilename := addEncExtension(file)
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
	plaintext, err := encryption.DecryptData(filedata, password)
	if err != nil {
		return err
	}

	outputFilename := removeEncExtension(encryptedFile)
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

// if .enc is not in the filename, add it, if it already has .enc then remove it
func addEncExtension(file string) string {
	if len(file) < 4 || file[len(file)-4:] != ".enc" {
		return file + ".enc"
	}
	return file
}

func removeEncExtension(file string) string {
	if len(file) < 4 || file[len(file)-4:] != ".enc" {
		return file
	}
	return file[:len(file)-4]
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
