package jobs

import (
	"fmt"
	"golang_starter_kit_2025/app/queue"
	"log"
)

// EmailJob represents an email sending job
type EmailJob struct {
	queue.BaseJob
}

// NewEmailJob creates a new email job
func NewEmailJob() queue.JobInterface {
	return &EmailJob{
		BaseJob: queue.BaseJob{
			Name: "email_job",
		},
	}
}

// Handle processes the email job
func (e *EmailJob) Handle() error {
	data := e.GetData()

	to, ok := data["to"].(string)
	if !ok {
		return fmt.Errorf("missing 'to' field in email job data")
	}

	subject, ok := data["subject"].(string)
	if !ok {
		return fmt.Errorf("missing 'subject' field in email job data")
	}

	body, ok := data["body"].(string)
	if !ok {
		return fmt.Errorf("missing 'body' field in email job data")
	}

	// Simulate email sending
	log.Printf("📧 Sending email to: %s", to)
	log.Printf("📧 Subject: %s", subject)
	log.Printf("📧 Body: %s", body)

	// Here you would integrate with your actual email service
	// For example: SendGrid, AWS SES, SMTP, etc.

	return nil
}

// SendEmailJob is a helper function to queue an email job
func SendEmailJob(to, subject, body string) error {
	job := NewEmailJob()
	job.SetData(map[string]interface{}{
		"to":      to,
		"subject": subject,
		"body":    body,
	})

	// You would get the queue instance from your app context
	// For now, this is just a placeholder
	return nil
}
