// Package env provides a typesafe way to parse environment variables using zog
package env

import (
	"os"
)

// Schema is a struct that holds the environment variables for the application.
type Schema struct {
	AllowedOrgins string
	APIPort       string
}

var initEnv = Schema{
	AllowedOrgins: os.Getenv("ALLOWED_ORGINS"),
	APIPort:       os.Getenv("API_PORT"),
}

// Values is a global variable that holds the parsed environment variables.
var Values Schema

func init() {
	// Validate ALLOWED_ORGINS
	if initEnv.AllowedOrgins == "" {
		Values.AllowedOrgins = "http://localhost:3000"
	}

	// Validate API_PORT
	if initEnv.APIPort == "" {
		Values.APIPort = "8080"
	}

}
