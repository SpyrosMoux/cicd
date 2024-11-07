package main

import (
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"github.com/spyrosmoux/cicd/common/helpers"
	"github.com/spyrosmoux/cicd/common/queue"
	"github.com/spyrosmoux/cicd/runner/pipelines"
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

	client := pipelineruns.NewClient(apiBaseUrl)

	go func() {
		for d := range msgs {
			slog.Info("Received a message with correlation id: " + d.CorrelationId)

			_, err := client.UpdatePipelineRunStatus(d.CorrelationId, pipelineruns.RUNNING.String())
			if err != nil {
				slog.Error("Failed to update pipeline with error: " + err.Error())
			}

			runResult := true
			err = pipelines.RunPipeline(string(d.Body))
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
				_, err = client.UpdatePipelineRunStatus(d.CorrelationId, pipelineruns.COMPLETED.String())
				if err != nil {
					slog.Error("Failed to update pipeline with error: " + err.Error())
				}
			}

			_, err = client.UpdatePipelineRunStatus(d.CorrelationId, pipelineruns.FAILED.String())
			if err != nil {
				slog.Error("Failed to update pipeline with error: " + err.Error())
			}
		}
	}()

	slog.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
