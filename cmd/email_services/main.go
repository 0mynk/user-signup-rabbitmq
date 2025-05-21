package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"user-signup-rabbitmq/pkg/common"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found ", err)
	}
}

func processEmail(event common.UserEvent) {
	// Get SMTP configuration
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	if smtpHost == "" || smtpPort == 0 {
		log.Println("SMTP configuration missing, skipping email send")
		return
	}

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", "ancilar@gmail.com")
	m.SetHeader("To", event.Email)
	m.SetHeader("Subject", "Welcome to Our Service!")

	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<h1>Welcome, %s %s!</h1>
			<p>Thank you for signing up!</p>
			<p>We're excited to have you on board.</p>
		</body>
		</html>
	`, event.FirstName, event.LastName)

	m.SetBody("text/html", htmlBody)

	// Send email using MailHog
	d := gomail.NewDialer(smtpHost, smtpPort, "", "")
	// d.Timeout = 10 * time.Second

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", event.Email, err)
		return
	}

	log.Printf("Email sent to %s", event.Email)
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
