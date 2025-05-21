package emails

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"user-signup-rabbitmq/pkg/common"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func ProcessEmail(event common.UserEvent) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found ", err)
	}
	// Email configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@yourdomain.com")
	m.SetHeader("To", event.Email)
	m.SetHeader("Subject", "Welcome to Our Service!")

	// HTML email content
	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<h1>Welcome, %s %s!</h1>
			<p>Thank you for signing up with us!</p>
			<p>We're excited to have you on board.</p>
			<p>Your account email: %s</p>
		</body>
		</html>
	`, event.FirstName, event.LastName, event.Email)

	// Plain text alternative
	textBody := fmt.Sprintf(
		"Welcome %s %s!\n\nThank you for signing up with us!\nYour account email: %s",
		event.FirstName, event.LastName, event.Email,
	)

	m.SetBody("text/plain", textBody)
	m.AddAlternative("text/html", htmlBody)

	// Send email
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", event.Email, err)
		return
	}

	log.Printf("Successfully sent welcome email to %s", event.Email)
}
