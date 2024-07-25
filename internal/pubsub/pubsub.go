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
