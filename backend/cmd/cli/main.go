// Package main is the entry point for the CLI application
package main

import (
	"log"
	"os"

	"github.com/edsonjaramillo/crpytid/backend/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "cryptid",
		Version: "0.1.0",
		Usage:   `Cryptid is a versatile CLI tool designed encrypting and decrypting files.`,
		Commands: []*cli.Command{
			commands.EncryptCommand,
			commands.DecryptCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
