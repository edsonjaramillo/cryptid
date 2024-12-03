package main

import (
	"log"
	"os"

	"github.com/edsonjaramillo/crpytid/internal/commands/aes"
	"github.com/edsonjaramillo/crpytid/internal/commands/jwt"
	"github.com/edsonjaramillo/crpytid/internal/commands/password"
	"github.com/urfave/cli/v2"
)

var Usage = `Cryptid is a versatile CLI tool designed for secure & efficient cryptographic operations. 
It supports, generating complex passwords, creating JWT secrets, AES file encryption & decryption. 
This tool is ideal for developers & security professionals who need reliable cryptographic utilities in their workflow.`

func main() {
	app := &cli.App{
		Name:     "cryptid",
		Version:  "v1.1.0",
		Usage:    Usage,
		Commands: []*cli.Command{password.PasswordCommand, jwt.JWTSecretsCommand, aes.AESCommand},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
