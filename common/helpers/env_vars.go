package helpers

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spyrosmoux/cicd/common/logger"
)

var logs = logger.NewLogger()

func LoadEnvVariable(variable string) string {
	if err := godotenv.Load(); err != nil {
		logs.Warn("no .env file found, attempting to read variables from host environment")
	}

	variableValue := getEnvOrExit(variable)

	return variableValue
}

func getEnvOrExit(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logs.WithFields(logrus.Fields{
			"variable": key,
		}).Fatal("environment variable not set")
	}
	return value
}
