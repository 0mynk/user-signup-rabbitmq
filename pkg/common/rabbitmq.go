package common

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "user_signups"
)

func ConnectRabbitMQ() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect: %s", err)
	}
	return conn
}

func SetupExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		ExchangeName,
		"fanout",
		true,  // Durable
		false, // Auto-deleted
		false, // Internal
		false, // No-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Exchange declaration failed: %s", err)
	}
}
