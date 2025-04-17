package queue

import (
	"encoding/json"

	"github.com/spyrosmoux/cicd/common/dto"
	"github.com/spyrosmoux/cicd/common/helpers"
	"github.com/spyrosmoux/cicd/common/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn             *amqp.Connection
	channel          *amqp.Channel
	rabbitmqHost     = helpers.LoadEnvVariable("RABBITMQ_HOST")
	rabbitmqUser     = helpers.LoadEnvVariable("RABBITMQ_USER")
	rabbitmqPassword = helpers.LoadEnvVariable("RABBITMQ_PASSWORD")
	rabbitmqPort     = helpers.LoadEnvVariable("RABBITMQ_PORT")
	logs             = logger.NewLogger()
)

// InitRabbitMQPublisher initializes the connection to RabbitMQ for the Api
func InitRabbitMQPublisher(channel string) {
	ch := establishConnection()

	_, err := ch.QueueDeclare(
		channel,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logs.WithError(err).Fatal("unable to initialize publisher")
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
		logs.WithError(err).Error("failed to publish message")
	}
}

func PublishLog(pipelineRunId string, logEntryDto dto.LogEntryDto) error {
	body, err := json.Marshal(logEntryDto)
	if err != nil {
		logs.WithError(err).Error("failed to marshal LogEntryDto into []byte")
		return err
	}

	err = channel.Publish(
		"",
		"logs",
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
		logs.WithError(err).Error("failed to publish message")
		return err
	}

	return nil
}

// InitRabbitMQConsumer initializes the connection to RabbitMQ for the Runner
func InitRabbitMQConsumer(channel string, prefetchSize int) <-chan amqp.Delivery {
	ch := establishConnection()

	q, err := ch.QueueDeclare(
		channel,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logs.WithError(err).Error("failed to declare a queue")
	}

	err = ch.Qos(prefetchSize, 0, false)
	if err != nil {
		logs.WithError(err).Error("failed to set QoS")
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
		logs.WithError(err).Error("failed to register a consumer")
	}

	return msgs
}

func establishConnection() *amqp.Channel {
	var err error
	conn, err = amqp.Dial("amqp://" + rabbitmqUser + ":" + rabbitmqPassword + "@" + rabbitmqHost + ":" + rabbitmqPort + "/")
	if err != nil {
		logs.WithError(err).Fatal("failed to connect to RabbitMQ")
	}

	channel, err = conn.Channel()
	if err != nil {
		logs.WithError(err).Fatal("failed to open a channel")
	}

	return channel
}
