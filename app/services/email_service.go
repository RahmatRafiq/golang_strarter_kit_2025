package services

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/config"

	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

// EmailService handles email sending operations
type EmailService struct {
	config *config.EmailConfig
}

// NewEmailService creates a new email service instance
func NewEmailService() *EmailService {
	return &EmailService{
		config: config.GetEmailConfig(),
	}
}

// EmailData holds data for email templates
type EmailData struct {
	AppName     string
	Subject     string
	Body        template.HTML
	Username    string
	VerifyURL   string
	ResetURL    string
	Year        int
	ExpiryHours int
}

// SendEmail sends an email with plain text body
func (s *EmailService) SendEmail(to, subject, body string) error {
	if !s.config.IsConfigured() {
		log.Warn().Msg("Email not configured, skipping email send")
		return fmt.Errorf("email service not configured")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromAddress))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	return s.send(m)
}

// SendHTMLEmail sends an email with HTML body
func (s *EmailService) SendHTMLEmail(to, subject, htmlBody string) error {
	if !s.config.IsConfigured() {
		log.Warn().Msg("Email not configured, skipping email send")
		return fmt.Errorf("email service not configured")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromAddress))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	return s.send(m)
}

// SendTemplatedEmail sends an email using a template
func (s *EmailService) SendTemplatedEmail(to, subject, templateName string, data map[string]interface{}) error {
	if !s.config.IsConfigured() {
		log.Warn().Msg("Email not configured, skipping email send")
		return fmt.Errorf("email service not configured")
	}

	// Render template
	htmlBody, err := s.renderTemplate(templateName, subject, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return s.SendHTMLEmail(to, subject, htmlBody)
}

// SendWelcomeEmail sends a welcome email to new users
func (s *EmailService) SendWelcomeEmail(to, username string) error {
	data := map[string]interface{}{
		"Username": username,
	}
	return s.SendTemplatedEmail(to, "Welcome to "+helpers.GetEnv("APP_NAME", "Our Platform"), "welcome.html", data)
}

// SendVerificationEmail sends an email verification link
func (s *EmailService) SendVerificationEmail(to, username, verifyURL string, expiryHours int) error {
	data := map[string]interface{}{
		"Username":    username,
		"VerifyURL":   verifyURL,
		"ExpiryHours": expiryHours,
	}
	return s.SendTemplatedEmail(to, "Verify Your Email Address", "verify_email.html", data)
}

// SendPasswordResetEmail sends a password reset link
func (s *EmailService) SendPasswordResetEmail(to, username, resetURL string, expiryHours int) error {
	data := map[string]interface{}{
		"Username":    username,
		"ResetURL":    resetURL,
		"ExpiryHours": expiryHours,
	}
	return s.SendTemplatedEmail(to, "Reset Your Password", "reset_password.html", data)
}

// renderTemplate renders an email template with base layout
func (s *EmailService) renderTemplate(templateName, subject string, data map[string]interface{}) (string, error) {
	// Get template paths
	templatesDir := filepath.Join("app", "templates", "emails")
	basePath := filepath.Join(templatesDir, "base.html")
	templatePath := filepath.Join(templatesDir, templateName)

	// Check if files exist
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return "", fmt.Errorf("base template not found: %s", basePath)
	}
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template not found: %s", templatePath)
	}

	// Parse templates
	tmpl, err := template.ParseFiles(basePath, templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse templates: %w", err)
	}

	// Render content template first
	var contentBuf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&contentBuf, filepath.Base(templateName), data); err != nil {
		return "", fmt.Errorf("failed to execute content template: %w", err)
	}

	// Prepare data for base template
	baseData := EmailData{
		AppName: helpers.GetEnv("APP_NAME", "Application"),
		Subject: subject,
		Body:    template.HTML(contentBuf.String()), // #nosec G203 - Content is from our own templates, not user input
		Year:    time.Now().Year(),
	}

	// Render base template with content
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base.html", baseData); err != nil {
		return "", fmt.Errorf("failed to execute base template: %w", err)
	}

	return buf.String(), nil
}

// send sends the email message
func (s *EmailService) send(m *gomail.Message) error {
	d := gomail.NewDialer(
		s.config.SMTPHost,
		s.config.SMTPPort,
		s.config.SMTPUsername,
		s.config.SMTPPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		log.Error().
			Err(err).
			Str("smtp_host", s.config.SMTPHost).
			Int("smtp_port", s.config.SMTPPort).
			Msg("Failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().
		Str("to", m.GetHeader("To")[0]).
		Str("subject", m.GetHeader("Subject")[0]).
		Msg("Email sent successfully")

	return nil
}
