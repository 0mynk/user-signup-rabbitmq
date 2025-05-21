
### Project Structure
```
user-signup-rabbitmq/
├── cmd/
│   ├── email_service/
│   │   └── main.go
│   └── publisher/
│       └── main.go
├── pkg/
│   └── common/
│       └── rabbitmq.go
├── .env
├── docker-compose.yml
├── go.mod
└── go.sum
```

### 1. docker-compose.yml
```yaml
version: '3'

services:
  mailhog:
    image: mailhog/mailhog
    ports:
      - "1025:1025"  # SMTP
      - "8025:8025"  # Web UI
  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"  # AMQP
      - "15672:15672" # Management UI
```


only for RabbitMQ on Docker
```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

### 2. .env
```env
# MailHog Configuration
SMTP_HOST=mailhog
SMTP_PORT=1025

# RabbitMQ Configuration
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
```

### 3. Run the System
```bash
# Start dependencies
docker-compose up -d

# Build and run services
go run cmd/email_service/main.go
go run cmd/publisher/main.go

# Access interfaces:
# - RabbitMQ Management: http://localhost:15672
# - MailHog UI: http://localhost:8025
```

### 6. Verify in MailHog
1. Open http://localhost:8025
2. You'll see emails like this:
```
From: noreply@example.com
To: user@example.com
Subject: Welcome to Our Service!

Welcome, John Doe!
Thank you for signing up!
```

### Key Features
1. **Local SMTP Testing**: Uses MailHog to capture emails
2. **Dockerized Environment**: RabbitMQ + MailHog in one command
3. **HTML Email Support**: Rich formatted emails
4. **Environment Variables**: Safe configuration management
5. **Error Handling**: Skip email if SMTP not configured

This setup allows you to:
- Test email content without real SMTP server
- Develop offline
- Inspect full email headers and content
- Safely experiment with email templates

**Troubleshooting Tips**:
- Check `docker-compose ps` to verify services are running
- Use `docker-compose logs` to view service logs
- Verify environment variables are loaded properly
- Check firewall settings if services can't connect