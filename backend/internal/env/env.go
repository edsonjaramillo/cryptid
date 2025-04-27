// Package env provides a typesafe way to parse environment variables using zog
package env

import (
	"log"
	"os"
	"regexp"

	z "github.com/Oudwins/zog"
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

// allowed_origins default is http://localhost:3000
var envSchema = z.Struct(z.Schema{
	"AllowedOrgins": z.String().Optional().Default("http://localhost:3000"),
	"APIPort":       z.String().Optional().Match(regexp.MustCompile(`^\d{1,5}$`)).Default("8080"),
})

// ConfigEnv parses the environment variables and returns a Schema struct.
func ConfigEnv() Schema {
	blank := Schema{}
	var issues = envSchema.Parse(initEnv, &blank)
	if issues != nil {
		log.Fatalf("Error parsing env variables: %v", issues)
	}

	return blank
}
