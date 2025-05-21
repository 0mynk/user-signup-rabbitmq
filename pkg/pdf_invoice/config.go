package pdf_invoice

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	OutputDir = "data/invoices"
)

type Config struct {
	UniDocAPIKey string
}

func LoadConfig() (*Config, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found ", err)
	}

	apiKey := os.Getenv("UNIDOC_LICENSE_API_KEY")
	if apiKey == "" {
		return nil, errors.New("UNIDOC_LICENSE_API_KEY is not set")
	}

	// Create output directory
	if err := os.MkdirAll(OutputDir, 0755); err != nil {
		return nil, err
	}

	return &Config{
		UniDocAPIKey: apiKey,
	}, nil
}
