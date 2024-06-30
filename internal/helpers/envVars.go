package helpers

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnvVariable(variable string) string {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, attempting to read from host environment variables")
	}

	variable = getEnvOrExit("API_PORT")

	return variable
}

func getEnvOrExit(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}
