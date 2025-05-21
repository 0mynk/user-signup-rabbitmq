package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"user-signup-rabbitmq/pkg/common"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn := common.ConnectRabbitMQ()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Channel creation failed: %s", err)
	}
	defer ch.Close()

	common.SetupExchange(ch)

	event := common.UserEvent{
		Email:     "johndoe@ancilar.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	body, _ := json.Marshal(event)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"user_signups",
		"",
		false,
		false,
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})

	if err != nil {
		log.Fatalf("Publish failed: %s", err)
	}

	log.Println("Published user signup event")
}
