package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
)

func InitRabbitMQ() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
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
