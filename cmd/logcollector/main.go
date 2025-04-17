package main

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spyrosmoux/cicd/common/db"
	"github.com/spyrosmoux/cicd/common/dto"
	"github.com/spyrosmoux/cicd/common/helpers"
	"github.com/spyrosmoux/cicd/common/logger"
	"github.com/spyrosmoux/cicd/common/queue"
	"github.com/spyrosmoux/cicd/logcollector"
)

var (
	apiPort       string
	dbHost        string
	dbPort        string
	dbUser        string
	dbPass        string
	dbName        string
	router        *gin.Engine
	logs          = logger.NewLogger()
	logRepository logcollector.LogRepository
	logService    logcollector.LogService
)

func init() {
	apiPort = helpers.LoadEnvVariable("LOGCOLLECTOR_PORT")
	dbHost = helpers.LoadEnvVariable("LOGCOLLECTOR_DB_HOST")
	dbPort = helpers.LoadEnvVariable("LOGCOLLECTOR_DB_PORT")
	dbUser = helpers.LoadEnvVariable("LOGCOLLECTOR_DB_USER")
	dbPass = helpers.LoadEnvVariable("LOGCOLLECTOR_DB_PASS")
	dbName = helpers.LoadEnvVariable("LOGCOLLECTOR_DB_NAME")

	// Initialize Db Connection
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"
	err := db.Init(dsn, "logcollector", &logcollector.LogEntry{})
	if err != nil {
		logs.Fatal(err)
	}

	router = logcollector.SetupRouter()

	logRepository = logcollector.NewLogRepository(db.DB)
	logService = logcollector.NewLogService(logRepository)
}

func main() {
	msgs := queue.InitRabbitMQConsumer("logs", 50)

	go func() {
		for d := range msgs {
			consumeLog(d)
			err := d.Ack(false)
			if err != nil {
				logs.WithError(err).Error("failed to acknowledge message")
			}
		}
	}()

	logs.WithFields(logrus.Fields{
		"port": apiPort,
	}).Info("starting server")
	logs.Fatal(router.Run(":" + apiPort))
}

func consumeLog(d amqp091.Delivery) {
	var logEntryDto dto.LogEntryDto
	err := json.Unmarshal(d.Body, &logEntryDto)
	if err != nil {
		logs.WithError(err).Error("failed to unmarshal logEntryDto")
		return
	}

	logEntry, err := logService.AddLog(logEntryDto)
	if err != nil {
		logs.Error(err)
		return
	}

	logs.WithFields(logrus.Fields{
		"logId": logEntry.Id,
		"runId": logEntry.RunId,
	}).Info("saved")
}
