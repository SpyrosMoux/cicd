package main

import (
	"log"
	"os"
	"spyrosmoux/api/internal/queue"
	"spyrosmoux/api/internal/routers"

	"github.com/joho/godotenv"
)

var (
	apiPort string
)

func init() {
	loadEnvVariables()
}

func loadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, attempting to read from host environment variables")
	}

	apiPort = getEnvOrExit("API_PORT")
}

func getEnvOrExit(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}

func main() {
	// Initialize RabbitMQ
	queue.InitRabbitMQ()

	// Setup routes
	router := routers.SetupRouter()

	log.Printf("Starting server on port %s", apiPort)
	log.Fatal(router.Run(":" + apiPort))
}
