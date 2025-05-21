# User Signup Event Processing System

This project implements a microservices-based system for processing user signup events using RabbitMQ as the message broker. The system consists of multiple services that handle different aspects of user signup processing, including email notifications and PDF invoice generation.

## System Architecture

The system is composed of the following components:

1. **Publisher Service**: Publishes user signup events to RabbitMQ
2. **Email Service**: Consumes events and sends welcome emails to new users
3. **PDF Service**: Consumes events and generates PDF invoices for new users

## Prerequisites

- Go 1.16 or higher
- RabbitMQ server
- SMTP server (or MailHog for local development)

## Project Structure

```
.
├── cmd/
│   ├── publisher/         # Event publisher service
│   ├── email_services/    # Email processing service
│   └── pdf_services/      # PDF generation service
├── pkg/
│   ├── common/           # Shared utilities and types
│   ├── emails/           # Email templates and utilities
│   └── pdf_invoice/      # PDF generation utilities
└── .env                  # Environment configuration
```

## Setup Guide

### 1. Install Dependencies

```bash
# Initialize Go modules (if not already done)
go mod init user-signup-rabbitmq

# Install and tidy up dependencies
go mod tidy
```

### 2. Configure Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# RabbitMQ Configuration
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# SMTP Configuration
SMTP_HOST=localhost
SMTP_PORT=1025  # Default MailHog port
```

### 3. Start RabbitMQ

```bash
# Using Docker
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

### 4. Start MailHog (for local email testing)

```bash
# Using Docker
docker run -d --name mailhog -p 1025:1025 -p 8025:8025 mailhog/mailhog
```

### 5. Run the Services

Open separate terminal windows for each service:

```bash
# Terminal 1 - Publisher Service
cd cmd/publisher
go run main.go

# Terminal 2 - Email Service
cd cmd/email_services
go run main.go

# Terminal 3 - PDF Service
cd cmd/pdf_services
go run main.go
```

## Testing the System

1. The publisher service will automatically send a test user signup event
2. The email service will process the event and send a welcome email
3. The PDF service will generate an invoice PDF for the new user

## Monitoring

- RabbitMQ Management Interface: http://localhost:15672 (guest/guest)
- MailHog Web Interface: http://localhost:8025

## Development

### Adding New Services

1. Create a new directory under `cmd/`
2. Implement the service using the common package utilities
3. Set up RabbitMQ consumer/producer as needed

### Common Package Usage

The `pkg/common` package provides shared functionality:

- `ConnectRabbitMQ()`: Establishes connection to RabbitMQ
- `SetupExchange()`: Configures the message exchange
- `UserEvent`: Common event structure for user signups

## Error Handling

- Services implement retry mechanisms for failed operations
- Messages are acknowledged only after successful processing
- Failed messages are requeued for retry

