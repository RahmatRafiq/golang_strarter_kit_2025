package tasks

import (
	"context"
	"fmt"

	"golang_starter_kit_2025/app/workers"

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
	return workers.NewTaskFromPayload(TypeSendEmail, payload)
}

// NewSendEmailTaskWithTemplate creates a new send email task with template
func NewSendEmailTaskWithTemplate(to, subject, template string) (*asynq.Task, error) {
	payload := SendEmailPayload{
		To:       to,
		Subject:  subject,
		Template: template,
	}
	return workers.NewTaskFromPayload(TypeSendEmail, payload)
}

// HandleSendEmailTask processes send email tasks
type HandleSendEmailTask struct {
	emailService EmailService
}

// NewHandleSendEmailTask creates a new send email task handler with dependencies
func NewHandleSendEmailTask(emailService EmailService) *HandleSendEmailTask {
	return &HandleSendEmailTask{
		emailService: emailService,
	}
}

// ProcessTask implements asynq.Handler interface
func (h *HandleSendEmailTask) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload SendEmailPayload
	if err := workers.UnmarshalTaskPayload(task, &payload); err != nil {
		return err
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
