package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("Starting Peril server...")
	connUrl := "amqp://guest:guest@0.0.0.0:5672/"
	conn, err := amqp091.Dial(connUrl)
	if err != nil {
		slog.Error("failed to connect to rabbitmq server", slog.String("error_msg", err.Error()))
		os.Exit(-1)
	}
	defer conn.Close()
	fmt.Println("connection successful")

	rmqpChan, err := conn.Channel()
	if err != nil {
		slog.Error("failed to create channel", slog.String("error_msg", err.Error()))
	}
	defer rmqpChan.Close()

	pubsub.PublishJSON(rmqpChan, routing.ExchangePerilDirect, routing.PauseKey,
		routing.PlayingState{IsPaused: true})

	<-signalChan
}
