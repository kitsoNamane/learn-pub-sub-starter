package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	connUrl := "amqp://guest:guest@0.0.0.0:5672/"
	conn, err := amqp091.Dial(connUrl)
	if err != nil {
		slog.Error("failed to connect to rabbitmq server", slog.String("error_msg", err.Error()))
		os.Exit(-1)
	}
	defer conn.Close()
	fmt.Println("connection successful")

	name, err := gamelogic.ClientWelcome()
	if err != nil {
		slog.Error("failed to get customer name", slog.String("error_msg", err.Error()))
		os.Exit(-1)
	}

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, fmt.Sprintf("pause.%s", name), routing.PauseKey, 0)
	if err != nil {
		slog.Error("failed to create and bind a queue", slog.String("error_msg", err.Error()))
		os.Exit(-1)
	}

	<-signalChan
}
