package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GoDotEnvVariable(key string) string {

	var env = os.Getenv("GO_ENV")
	var err error
	if env == "dev" {
		err = godotenv.Load(".env.dev")
	} else if env == "prod" {
		err = godotenv.Load(".env.prod")
	} else {
		err = godotenv.Load(".env")
	}
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
