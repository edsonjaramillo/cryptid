// Package main is the entry point for the CLI application
package main

import (
	"log"
	"os"

	"github.com/edsonjaramillo/hyde/backend/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "Hyde",
		Version: "0.1.0",
		Usage:   `Hyde is a versatile CLI tool designed encrypting and decrypting files.`,
		Commands: []*cli.Command{
			commands.DecryptCommand,
			commands.EncryptCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
