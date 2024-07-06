package main

import (
	"log"
	"spyrosmoux/api/internal/helpers"
	"spyrosmoux/api/internal/queue"
	"spyrosmoux/api/internal/routers"
	"strconv"
)

var (
	apiPort        string
	pemFilePath    string
	clientId       string
	installationId int64
)

func init() {
	apiPort = helpers.LoadEnvVariable("API_PORT")
	pemFilePath = helpers.LoadEnvVariable("PEM_FILE_PATH")
	clientId = helpers.LoadEnvVariable("CLIENT_ID")
	installationId, _ = strconv.ParseInt(helpers.LoadEnvVariable("INSTALLATION_ID"), 10, 64)
}

func main() {
	// Initialize RabbitMQ
	queue.InitRabbitMQ()

	// Setup routes
	router := routers.SetupRouter()

	log.Printf("Starting server on port %s", apiPort)
	log.Fatal(router.Run(":" + apiPort))
}
