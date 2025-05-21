package main

import (
	"encoding/json"
	"log"
	"time"

	"user-signup-rabbitmq/pkg/common"
	"user-signup-rabbitmq/pkg/pdf_invoice"
)

func generatePDF(event common.UserEvent) {
	log.Printf("Generating PDF for %s %s", event.FirstName, event.LastName)

	// Initialize PDF generator
	cfg, err := pdf_invoice.LoadConfig()
	if err != nil {
		log.Printf("PDF config error: %v", err)
		return
	}

	generator := pdf_invoice.NewGenerator(cfg)

	// Convert to PDF data structure
	pdfData := &common.UserEvent{
		Email:     event.Email,
		FirstName: event.FirstName,
		LastName:  event.LastName,
	}

	// Generate and save PDF
	filename, err := generator.GenerateInvoice(pdfData)
	if err != nil {
		log.Printf("PDF generation failed: %v", err)
		return
	}

	log.Printf("PDF generated successfully: %s", filename)
	time.Sleep(2 * time.Second)
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
		"pdf_service",
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

	log.Println("PDF service started. Waiting for events...")

	for d := range msgs {
		var event common.UserEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Printf("Error decoding message: %s", err)
			d.Nack(false, true) // Requeue message
			continue
		}

		go func() {
			generatePDF(event)
			d.Ack(false)
		}()
	}
}
