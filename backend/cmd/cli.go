package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var Usage = `Cryptid is a versatile CLI tool designed encrypting and decrypting files.`

func main() {
	app := &cli.App{
		Name:    "cryptid",
		Version: "0.0.1",
		Usage:   Usage,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
