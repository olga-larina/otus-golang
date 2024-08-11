package rabbit

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	logger       Logger
	uri          string
	exchangeName string
	exchangeType string
	queueName    string
	routingKey   string

	connection *amqp.Connection
	channel    *amqp.Channel
}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, msg string, args ...any)
}

func NewQueue(logger Logger, uri string, exchangeName string, exchangeType string, queueName string, routingKey string) *Queue {
	return &Queue{
		logger:       logger,
		uri:          uri,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queueName:    queueName,
		routingKey:   routingKey,
	}
}

func (q *Queue) Start(ctx context.Context) error {
	q.logger.Info(ctx, "starting rabbit queue")

	var err error

	q.connection, err = amqp.Dial(q.uri)
	if err != nil {
		return err
	}
	q.logger.Info(ctx, "got rabbit queue connection, getting channel")

	q.channel, err = q.connection.Channel()
	if err != nil {
		return err
	}
	q.logger.Info(ctx, "got rabbit queue channel, declaring exchange", "exchangeName", q.exchangeName, "exchangeType", q.exchangeType)

	err = q.channel.ExchangeDeclare(
		q.exchangeName, // name
		q.exchangeType, // type
		true,           // durable
		false,          // autoDelete
		false,          // internal
		false,          // noWait
		nil,            // arguments
	)
	if err != nil {
		return err
	}
	q.logger.Info(ctx, "rabbit exchange declared, declaring queue", "queueName", q.queueName)

	queue, err := q.channel.QueueDeclare(
		q.queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	q.logger.Info(ctx, "declared new rabbit queue, declaring binding", "routingKey", q.routingKey)

	err = q.channel.QueueBind(
		queue.Name,
		q.routingKey,
		q.exchangeName,
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	q.logger.Info(ctx, "queue bound to exchange, rabbit queue started")
	return nil
}

func (q *Queue) Stop(ctx context.Context) error {
	q.logger.Info(ctx, "stopping rabbit queue")

	err := q.connection.Close()
	if err != nil {
		return err
	}

	q.logger.Info(ctx, "stopped rabbit queue")
	return nil
}
