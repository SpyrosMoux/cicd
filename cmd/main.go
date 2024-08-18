package main

import (
	"log"
	"spyrosmoux/api/internal"
	"spyrosmoux/api/internal/auth"
	"spyrosmoux/api/internal/db"
	"spyrosmoux/api/internal/helpers"
	"spyrosmoux/api/internal/queue"
)

var (
	apiPort string
)

func init() {
	apiPort = helpers.LoadEnvVariable("API_PORT")
}

func main() {
	// Initialize DB connection
	dbConnection, err := db.InitDB()
	if err != nil {
		log.Fatal("Error initializing database with error: " + err.Error())
	}

	// Initialize RabbitMQ
	queue.InitRabbitMQ()

	// Initialize SuperTokens
	auth.InitSuperTokens()

	// Setup routes
	// Pass dbConnection to initialize handlers/services/repositories
	router := internal.SetupRouter(dbConnection)

	log.Printf("Starting server on port %s", apiPort)
	log.Fatal(router.Run(":" + apiPort))
}
