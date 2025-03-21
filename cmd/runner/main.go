package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"github.com/spyrosmoux/cicd/api/sdk"
	"github.com/spyrosmoux/cicd/common/dto"
	"github.com/spyrosmoux/cicd/common/helpers"
	"github.com/spyrosmoux/cicd/common/logger"
	"github.com/spyrosmoux/cicd/common/queue"
	"github.com/spyrosmoux/cicd/runner/dirmanagement"
	"github.com/spyrosmoux/cicd/runner/pipelines"
	"gopkg.in/yaml.v3"
)

var (
	apiBaseUrl string
	logs       *logrus.Logger
	runOnce    string
)

func init() {
	logs = logger.NewLogger()
	apiBaseUrl = helpers.LoadEnvVariable("API_BASE_URL")
	err := dirmanagement.InitGlobalDM()
	if err != nil {
		logs.WithError(err).Error("failed to initialize file system")
		os.Exit(1)
	}
	pipelines.SetPredefinedVars()
	runOnce = helpers.LoadEnvVariable("RUN_ONCE")
}

func main() {
	msgs := queue.InitRabbitMQConsumer()

	var forever chan struct{}

	client := sdk.NewClient(apiBaseUrl)
	svc := pipelines.NewService(logs)

	go func() {
		for d := range msgs {
			runner(client, svc, d)
			if runOnce == "true" {
				os.Exit(0)
			}
		}
	}()

	logs.Info("[*] Waiting for messages")
	<-forever
}

func runner(client *sdk.Client, svc pipelines.Service, d amqp091.Delivery) {
	logs.WithFields(logrus.Fields{
		"id": d.CorrelationId,
	}).Info("received a message with correlation id")

	var publishRunDto dto.PublishRunDto
	err := json.Unmarshal(d.Body, &publishRunDto)
	if err != nil {
		logs.WithError(err).Error("failed to unmarshal publishRunDto")
	}

	_, err = client.UpdatePipelineRun(d.CorrelationId, dto.UpdatePipelineRunDto{
		Status:      pipelineruns.RUNNING.String(),
		TimeStarted: time.Now().Unix(),
	})
	if err != nil {
		logs.WithError(err).Error("failed to update pipelineRun")
	}

	var pipeline pipelines.Pipeline
	err = yaml.Unmarshal(publishRunDto.PipelineAsBytes, &pipeline)
	if err != nil {
		logs.WithError(err).Error("filaed to unmarshal pipeline")
	}

	runResult := true
	runError := svc.RunPipeline(pipeline, publishRunDto.Metadata)
	if runError != nil {
		runResult = false
	}

	// Acknowledge the message after successful processing
	err = d.Ack(false)
	if err != nil {
		logs.WithError(err).Error("failed to acknowledge message")
	}

	if runResult {
		_, err = client.UpdatePipelineRun(d.CorrelationId, dto.UpdatePipelineRunDto{
			Status:    pipelineruns.COMPLETED.String(),
			TimeEnded: time.Now().Unix(),
		})
		if err != nil {
			logs.WithError(err).Error("failed to update pipeline")
		}
	} else {
		_, err = client.UpdatePipelineRun(d.CorrelationId, dto.UpdatePipelineRunDto{
			Status:    pipelineruns.FAILED.String(),
			Error:     runError.Error(),
			TimeEnded: time.Now().Unix(),
		})
		if err != nil {
			logs.WithError(err).Error("failed to update pipeline")
		}
	}
}
