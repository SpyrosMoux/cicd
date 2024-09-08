package main

import (
	"log"

	"github.com/spyrosmoux/api/internal/helpers"
	"github.com/spyrosmoux/api/internal/queue"
	"github.com/spyrosmoux/api/internal/routers"
	"github.com/spyrosmoux/api/internal/supertokens"
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
