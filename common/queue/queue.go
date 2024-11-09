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

// TODO(spyrosmoux) rename to InitRabbitMQPublisher and make any changes if necessary to generify this
func InitRabbitMQ() {
	var err error
	conn, err = amqp.Dial("amqp://" + rabbitmqUser + ":" + rabbitmqPassword + "@" + rabbitmqHost + ":" + rabbitmqPort + "/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	_, err = channel.QueueDeclare(
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
			ContentType:   "text/plain",
			Body:          body,
			DeliveryMode:  amqp.Persistent,
			CorrelationId: pipelineRunId,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
}

// InitRabbitMQRunner initializes the connection to RabbitMQ for the Runner
// TODO(spyrosmoux) rename to InitRabbitMQConsumer and make any changes if necessary to generify this
func InitRabbitMQRunner() <-chan amqp.Delivery {
	conn, err := amqp.Dial("amqp://" + rabbitmqUser + ":" + rabbitmqPassword + "@" + rabbitmqHost + ":" + rabbitmqPort + "/")
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ " + err.Error())
		os.Exit(1)
	}

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("Failed to open a channel: " + err.Error())
	}

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
