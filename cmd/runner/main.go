package main

import (
	"github.com/spyrosmoux/cicd/api/entities"
	"github.com/spyrosmoux/cicd/api/sdk"
	"github.com/spyrosmoux/cicd/common/helpers"
	"github.com/spyrosmoux/cicd/common/queue"
	"github.com/spyrosmoux/cicd/runner/pipelines"
	"gopkg.in/yaml.v3"
	"log/slog"
)

var (
	apiBaseUrl string
)

func init() {
	apiBaseUrl = helpers.LoadEnvVariable("API_BASE_URL")
}

func main() {
	msgs := queue.InitRabbitMQRunner()

	var forever chan struct{}

	client := sdk.NewClient(apiBaseUrl)

	go func() {
		for d := range msgs {
			slog.Info("Received a message with correlation id: " + d.CorrelationId)

			_, err := client.UpdatePipelineRunStatus(d.CorrelationId, entities.RUNNING.String())
			if err != nil {
				slog.Error("Failed to update pipeline with error: " + err.Error())
			}

			var pipeline pipelines.Pipeline
			err = yaml.Unmarshal(d.Body, &pipeline)
			if err != nil {
				slog.Error(err.Error())
			}

			runResult := true
			err = pipelines.RunPipeline(pipeline)
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
				_, err = client.UpdatePipelineRunStatus(d.CorrelationId, entities.COMPLETED.String())
				if err != nil {
					slog.Error("Failed to update pipeline with error: " + err.Error())
				}
			} else {
				_, err = client.UpdatePipelineRunStatus(d.CorrelationId, entities.FAILED.String())
				if err != nil {
					slog.Error("Failed to update pipeline with error: " + err.Error())
				}
			}
		}
	}()

	slog.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
