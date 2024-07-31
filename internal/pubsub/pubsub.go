package pubsub

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp091.Channel, exchange, key string, val T) error {
	valJSON, err := json.Marshal(val)
	if err != nil {
		slog.Error("failed to parse json", slog.String("error_msg", err.Error()))
		return err
	}

	ch.PublishWithContext(context.Background(), exchange, key, false, false,
		amqp091.Publishing{ContentType: "application/json", Body: valJSON})
	return nil
}

func DeclareAndBind(conn *amqp091.Connection, exchange, queueName, key string, simpleQueueType int) (*amqp091.Channel, amqp091.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		slog.Error("failed to create channel", slog.String("error_msg", err.Error()))
		return nil, amqp091.Queue{}, err
	}
	defer ch.Close()

	var durable, autoDelete, exclusive bool
	if simpleQueueType == 0 {
		durable = false
		autoDelete = true
		exclusive = true
	} else {
		durable = true
		autoDelete = false
		exclusive = false
	}

	queue, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		slog.Error("failed to declare a queue", slog.String("error_msg", err.Error()))
		return nil, amqp091.Queue{}, err
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		slog.Error("failed to bind queue to exchange", slog.String("error_msg", err.Error()))
		return nil, amqp091.Queue{}, err
	}

	return ch, queue, nil
}
