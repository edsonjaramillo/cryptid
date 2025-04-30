// Package main is the entry point for the CLI application
package main

import (
	"context"
	"log"
	"os"

	"github.com/edsonjaramillo/hyde/backend/internal/commands"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:    "hyde",
		Version: "0.1.0",
		Usage:   `Hyde is a versatile CLI tool designed encrypting and decrypting files.`,
		Commands: []*cli.Command{
			commands.DecryptCommand,
			commands.EncryptCommand,
		},
		HideHelpCommand: true,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
