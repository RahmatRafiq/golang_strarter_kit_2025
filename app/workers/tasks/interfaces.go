package tasks

// EmailService defines the interface for email operations
// This allows for dependency injection and easier testing
type EmailService interface {
	SendEmail(to, subject, body string) error
	SendHTMLEmail(to, subject, htmlBody string) error
}
