package main

import (
	"encoding/json"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"github.com/spyrosmoux/cicd/api/sdk"
	"github.com/spyrosmoux/cicd/common/dto"
	"github.com/spyrosmoux/cicd/common/helpers"
	"github.com/spyrosmoux/cicd/common/queue"
	"github.com/spyrosmoux/cicd/runner/dirmanagement"
	"github.com/spyrosmoux/cicd/runner/pipelines"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"time"
)

var (
	apiBaseUrl string
)

func init() {
	apiBaseUrl = helpers.LoadEnvVariable("API_BASE_URL")
	err := dirmanagement.InitGlobalDM()
	if err != nil {
		slog.Error("Failed to initialize file system: " + err.Error())
		os.Exit(1)
	}
}

func main() {
	msgs := queue.InitRabbitMQConsumer()

	var forever chan struct{}

	client := sdk.NewClient(apiBaseUrl)
	svc := pipelines.NewService()

	go func() {
		for d := range msgs {
			slog.Info("Received a message with correlation id: " + d.CorrelationId)

			var publishRunDto dto.PublishRunDto
			err := json.Unmarshal(d.Body, &publishRunDto)
			if err != nil {
				slog.Error(err.Error())
			}

			_, err = client.UpdatePipelineRun(d.CorrelationId, dto.UpdatePipelineRunDto{
				Status:      pipelineruns.RUNNING.String(),
				TimeStarted: time.Now().Unix(),
			})
			if err != nil {
				slog.Error("Failed to update pipeline with error: " + err.Error())
			}

			var pipeline pipelines.Pipeline
			err = yaml.Unmarshal(publishRunDto.PipelineAsBytes, &pipeline)
			if err != nil {
				slog.Error(err.Error())
			}

			runResult := true
			err = svc.RunPipeline(pipeline, publishRunDto.Metadata)
			if err != nil {
				runResult = false
				slog.Error("Failed to run pipeline with error: " + err.Error())
			}

			// Acknowledge the message after successful processing
			err = d.Ack(false)
			if err != nil {
				slog.Error("Failed to acknowledge message: " + err.Error())
			}

			if runResult {
				_, err = client.UpdatePipelineRun(d.CorrelationId, dto.UpdatePipelineRunDto{
					Status:    pipelineruns.COMPLETED.String(),
					TimeEnded: time.Now().Unix(),
				})
				if err != nil {
					slog.Error("Failed to update pipeline with error: " + err.Error())
				}
			} else {
				_, err = client.UpdatePipelineRun(d.CorrelationId, dto.UpdatePipelineRunDto{
					Status:    pipelineruns.FAILED.String(),
					TimeEnded: time.Now().Unix(),
				})
				if err != nil {
					slog.Error("Failed to update pipeline with error: " + err.Error())
				}
			}
		}
	}()

	slog.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
