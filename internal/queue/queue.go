package queue

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spyrosmoux/api/internal/helpers"
)

var (
	conn             *amqp.Connection
	channel          *amqp.Channel
	rabbitmqHost     = helpers.LoadEnvVariable("RABBITMQ_HOST")
	rabbitmqUser     = helpers.LoadEnvVariable("RABBITMQ_USER")
	rabbitmqPassword = helpers.LoadEnvVariable("RABBITMQ_PASSWORD")
	rabbitmqPort     = helpers.LoadEnvVariable("RABBITMQ_PORT")
)

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
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
}

func PublishJob(body string) {
	err := channel.Publish(
		"",
		"jobs",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
}
