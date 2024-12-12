package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/api/config"
	"github.com/spyrosmoux/cicd/api/middlewares"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"github.com/spyrosmoux/cicd/api/routes"
	"github.com/spyrosmoux/cicd/common/helpers"
	"github.com/spyrosmoux/cicd/common/queue"
	"log"
)

var (
	apiPort string
	dbHost  string
	dbPort  string
	dbUser  string
	dbPass  string
	dbName  string
	router  *gin.Engine
)

func init() {
	apiPort = helpers.LoadEnvVariable("API_PORT")
	dbHost = helpers.LoadEnvVariable("DB_HOST")
	dbPort = helpers.LoadEnvVariable("DB_PORT")
	dbUser = helpers.LoadEnvVariable("DB_USER")
	dbPass = helpers.LoadEnvVariable("DB_PASS")
	dbName = helpers.LoadEnvVariable("DB_NAME")

	// Initialize Db Connection
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"
	config.Init(dsn, &pipelineruns.PipelineRun{})

	// Initialize RabbitMQ
	queue.InitRabbitMQ()

	// Initialize SuperTokens
	middlewares.InitSuperTokens()

	// Setup routes
	router = routes.SetupRouter()
}

func main() {
	log.Printf("Starting server on port %s", apiPort)
	log.Fatal(router.Run(":" + apiPort))
}
