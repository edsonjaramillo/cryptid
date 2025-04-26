package env

import (
	"log"
	"os"
	"regexp"

	z "github.com/Oudwins/zog"
)

type EnvSchema struct {
	ALLOWED_ORIGINS string
	API_PORT        string
}

var initEnv = EnvSchema{
	ALLOWED_ORIGINS: os.Getenv("ALLOWED_ORIGINS"),
	API_PORT:        os.Getenv("API_PORT"),
}

// allowed_origins default is http://localhost:3000
var envSchema = z.Struct(z.Schema{
	"ALLOWED_ORIGINS": z.String().Optional().Default("http://localhost:3000"),
	"API_PORT":        z.String().Optional().Match(regexp.MustCompile(`^\d{1,5}$`)).Default("8080"),
})

func ConfigEnv() EnvSchema {
	blank := EnvSchema{}
	var issues = envSchema.Parse(initEnv, &blank)
	if issues != nil {
		log.Fatalf("Error parsing env variables: %v", issues)
		os.Exit(1)
	}

	return blank
}
