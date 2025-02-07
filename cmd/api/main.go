package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spyrosmoux/cicd/api/config"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"github.com/spyrosmoux/cicd/api/routes"
	"github.com/spyrosmoux/cicd/common/helpers"
	"github.com/spyrosmoux/cicd/common/logger"
	"github.com/spyrosmoux/cicd/common/queue"
)

var (
	apiPort string
	dbHost  string
	dbPort  string
	dbUser  string
	dbPass  string
	dbName  string
	router  *gin.Engine
	logs    = logger.NewLogger()
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
	err := config.Init(dsn, &pipelineruns.PipelineRun{})
	if err != nil {
		logs.Fatal(err)
	}

	// Initialize RabbitMQ
	queue.InitRabbitMQPublisher()

	// Setup routes
	router = routes.SetupRouter()
}

func main() {
	logs.WithFields(logrus.Fields{
		"port": apiPort,
	}).Info("starting server")
	logs.Fatal(router.Run(":" + apiPort))
}
