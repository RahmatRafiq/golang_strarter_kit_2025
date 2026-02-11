package config

import (
	"os"
	"strconv"
)

// EmailConfig holds the SMTP configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPUsername string
	SMTPPassword string
	FromAddress  string
	FromName     string
	SMTPPort     int
}

// GetEmailConfig returns email configuration from environment variables
func GetEmailConfig() *EmailConfig {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil || port == 0 {
		port = 587 // Default SMTP port for TLS
	}

	fromName := os.Getenv("SMTP_FROM_NAME")
	if fromName == "" {
		fromName = os.Getenv("APP_NAME")
	}

	return &EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     port,
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromAddress:  os.Getenv("SMTP_FROM_ADDRESS"),
		FromName:     fromName,
	}
}

// IsConfigured checks if email is properly configured
func (c *EmailConfig) IsConfigured() bool {
	return c.SMTPHost != "" &&
		c.SMTPPort > 0 &&
		c.SMTPUsername != "" &&
		c.SMTPPassword != "" &&
		c.FromAddress != ""
}
