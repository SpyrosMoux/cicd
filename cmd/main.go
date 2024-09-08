package main

import (
	"log"
	"spyrosmoux/api/internal/helpers"
	"spyrosmoux/api/internal/queue"
	"spyrosmoux/api/internal/routers"
	"spyrosmoux/api/internal/supertokens"
)

var (
	apiPort string
)

func init() {
	apiPort = helpers.LoadEnvVariable("API_PORT")
}

func main() {
	// Initialize RabbitMQ
	queue.InitRabbitMQ()

	// Initialize SuperTokens
	supertokens.InitSuperTokens()

	// Setup routes
	router := routers.SetupRouter()

	log.Printf("Starting server on port %s", apiPort)
	log.Fatal(router.Run(":" + apiPort))
}
