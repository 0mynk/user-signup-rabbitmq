package main

import (
	"encoding/json"
	"log"
	"time"

	"user-signup-rabbitmq/pkg/common"
)

func processEmail(event common.UserEvent) {
	// Get SMTP configuration
	// Stimulate email sending
	log.Printf("Sending email to %s", event.Email)
	log.Printf("From the Event %s", event)
	time.Sleep(1 * time.Second)
}

func main() {
	conn := common.ConnectRabbitMQ()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Channel creation failed: %s", err)
	}
	defer ch.Close()

	common.SetupExchange(ch)

	q, err := ch.QueueDeclare(
		"email_service",
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Queue declaration failed: %s", err)
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"user_signups",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Queue bind failed: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false, // Auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Consume failed: %s", err)
	}

	log.Println("Email service started. Waiting for events...")

	for d := range msgs {
		var event common.UserEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Printf("Error decoding message: %s", err)
			d.Nack(false, true) // Requeue message
			continue
		}

		go func() {
			processEmail(event)
			d.Ack(false)
		}()
	}
}
