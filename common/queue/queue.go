package queue

import (
	"github.com/spyrosmoux/cicd/common/helpers"
	"log"
	"log/slog"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn             *amqp.Connection
	channel          *amqp.Channel
	rabbitmqHost     = helpers.LoadEnvVariable("RABBITMQ_HOST")
	rabbitmqUser     = helpers.LoadEnvVariable("RABBITMQ_USER")
	rabbitmqPassword = helpers.LoadEnvVariable("RABBITMQ_PASSWORD")
	rabbitmqPort     = helpers.LoadEnvVariable("RABBITMQ_PORT")
)

// InitRabbitMQPublisher initializes the connection to RabbitMQ for the Api
func InitRabbitMQPublisher() {
	channel = establishConnection()

	_, err := channel.QueueDeclare(
		"jobs",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
}

func PublishJob(pipelineRunId string, body []byte) {
	err := channel.Publish(
		"",
		"jobs",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			DeliveryMode:  amqp.Persistent,
			CorrelationId: pipelineRunId,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
}

// InitRabbitMQConsumer initializes the connection to RabbitMQ for the Runner
func InitRabbitMQConsumer() <-chan amqp.Delivery {
	ch := establishConnection()

	q, err := ch.QueueDeclare(
		"jobs",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Error("Failed to declare a queue: " + err.Error())
	}

	// Set QoS (Quality of Service) to prefetch 1 message at a time
	err = ch.Qos(1, 0, false)
	if err != nil {
		slog.Error("Failed to set QoS: " + err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Error("Failed to register a consumer: " + err.Error())
	}

	return msgs
}

func establishConnection() *amqp.Channel {
	var err error
	conn, err = amqp.Dial("amqp://" + rabbitmqUser + ":" + rabbitmqPassword + "@" + rabbitmqHost + ":" + rabbitmqPort + "/")
	if err != nil {
		slog.Error("failed to connect to RabbitMQ " + err.Error())
		os.Exit(1)
	}

	channel, err = conn.Channel()
	if err != nil {
		slog.Error("failed to open a channel: " + err.Error())
		os.Exit(1)
	}

	return channel
}
