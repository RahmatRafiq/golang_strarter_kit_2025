package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TypeSendEmail = "email:send"
)

// SendEmailPayload represents the payload for sending an email
type SendEmailPayload struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Template string `json:"template,omitempty"` // Optional template name
}

// NewSendEmailTask creates a new send email task
func NewSendEmailTask(to, subject, body string) (*asynq.Task, error) {
	payload := SendEmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	return newSendEmailTaskFromPayload(payload)
}

// NewSendEmailTaskWithTemplate creates a new send email task with template
func NewSendEmailTaskWithTemplate(to, subject, template string) (*asynq.Task, error) {
	payload := SendEmailPayload{
		To:       to,
		Subject:  subject,
		Template: template,
	}
	return newSendEmailTaskFromPayload(payload)
}

func newSendEmailTaskFromPayload(payload SendEmailPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeSendEmail, data), nil
}

// HandleSendEmailTask processes send email tasks
type HandleSendEmailTask struct {
	emailService interface {
		SendEmail(to, subject, body string) error
		SendHTMLEmail(to, subject, htmlBody string) error
	}
}

// NewHandleSendEmailTask creates a new send email task handler with dependencies
func NewHandleSendEmailTask(emailService interface {
	SendEmail(to, subject, body string) error
	SendHTMLEmail(to, subject, htmlBody string) error
}) *HandleSendEmailTask {
	return &HandleSendEmailTask{
		emailService: emailService,
	}
}

// ProcessTask implements asynq.Handler interface
func (h *HandleSendEmailTask) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload SendEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Info().
		Str("to", payload.To).
		Str("subject", payload.Subject).
		Str("template", payload.Template).
		Msg("Processing send email task")

	// Send email using email service
	if h.emailService != nil {
		var err error
		if payload.Template != "" {
			// If template is specified, body is HTML
			err = h.emailService.SendHTMLEmail(payload.To, payload.Subject, payload.Body)
		} else {
			// Plain text email
			err = h.emailService.SendEmail(payload.To, payload.Subject, payload.Body)
		}

		if err != nil {
			log.Error().
				Err(err).
				Str("to", payload.To).
				Str("subject", payload.Subject).
				Msg("Failed to send email")
			return fmt.Errorf("failed to send email: %w", err)
		}

		log.Info().
			Str("to", payload.To).
			Str("subject", payload.Subject).
			Msg("Email sent successfully via worker")
	} else {
		log.Warn().Msg("Email service not configured, email not sent")
	}

	return nil
}
