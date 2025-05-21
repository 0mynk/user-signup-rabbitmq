package pdf_invoice

import (
	"fmt"
	"path/filepath"
	"time"
	"user-signup-rabbitmq/pkg/common"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
)

type PDFGenerator struct {
	config *Config
}

var counter = 1

func NewGenerator(cfg *Config) *PDFGenerator {
	return &PDFGenerator{config: cfg}
}

// Add this license initialization function
func initUniDoc(apiKey string) error {
	if err := license.SetMeteredKey(apiKey); err != nil {
		return fmt.Errorf("failed to set UniDoc license: %v", err)
	}
	return nil
}

func (g *PDFGenerator) GenerateInvoice(data *common.UserEvent) (string, error) {
	// Initialize UniDoc license
	initUniDoc(g.config.UniDocAPIKey)

	c := creator.New()
	c.SetPageMargins(50, 50, 50, 50)

	// Create invoice content
	if err := g.createHeader(c); err != nil {
		return "", err
	}

	if err := g.createCustomerInfo(c, data); err != nil {
		return "", err
	}

	// Generate filename
	filename := filepath.Join(OutputDir,
		fmt.Sprintf("%d-invoice-%s-%d.pdf",
			counter,
			data.LastName,
			time.Now().Unix(),
		))

	counter++

	// Save PDF
	if err := c.WriteToFile(filename); err != nil {
		return "", fmt.Errorf("failed to save PDF: %w", err)
	}

	return filename, nil
}

func (g *PDFGenerator) createHeader(c *creator.Creator) error {
	// Simple header with a title
	para := c.NewParagraph("Invoice")
	para.SetFontSize(24)
	para.SetMargins(0, 0, 20, 10)
	para.SetTextAlignment(creator.TextAlignmentCenter)

	return c.Draw(para)
}

func (g *PDFGenerator) createCustomerInfo(c *creator.Creator, data *common.InvoiceData) error {
	// Basic customer info: name and email
	fullName := fmt.Sprintf("%s %s", data.FirstName, data.LastName)
	para := c.NewParagraph(fmt.Sprintf("Customer Name: %s\nEmail: %s", fullName, data.Email))
	para.SetFontSize(12)
	para.SetMargins(0, 0, 10, 10)

	return c.Draw(para)
}
