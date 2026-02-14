package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"golang_starter_kit_2025/app/workers/tasks"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewSendEmailTask tests task creation
func TestNewSendEmailTask(t *testing.T) {
	t.Run("create task with valid data", func(t *testing.T) {
		to := "test@example.com"
		subject := "Test Subject"
		body := "Test Body"

		task, err := tasks.NewSendEmailTask(to, subject, body)

		require.NoError(t, err, "should create task without error")
		assert.Equal(t, tasks.TypeSendEmail, task.Type(), "task type should match")

		// Verify payload
		var payload tasks.SendEmailPayload
		err = json.Unmarshal(task.Payload(), &payload)
		require.NoError(t, err)
		assert.Equal(t, to, payload.To)
		assert.Equal(t, subject, payload.Subject)
		assert.Equal(t, body, payload.Body)
		assert.Empty(t, payload.Template)
	})

	t.Run("create task with template", func(t *testing.T) {
		to := "test@example.com"
		subject := "Welcome"
		template := "welcome_email"

		task, err := tasks.NewSendEmailTaskWithTemplate(to, subject, template)

		require.NoError(t, err)
		assert.Equal(t, tasks.TypeSendEmail, task.Type())

		var payload tasks.SendEmailPayload
		err = json.Unmarshal(task.Payload(), &payload)
		require.NoError(t, err)
		assert.Equal(t, template, payload.Template)
	})
}

// mockEmailService implements the email service interface for testing
type mockEmailService struct {
	sendEmailCalled     bool
	sendHTMLEmailCalled bool
	shouldFail          bool
	lastTo              string
	lastSubject         string
	lastBody            string
}

func (m *mockEmailService) SendEmail(to, subject, body string) error {
	m.sendEmailCalled = true
	m.lastTo = to
	m.lastSubject = subject
	m.lastBody = body

	if m.shouldFail {
		return errors.New("mock email send failed")
	}
	return nil
}

func (m *mockEmailService) SendHTMLEmail(to, subject, htmlBody string) error {
	m.sendHTMLEmailCalled = true
	m.lastTo = to
	m.lastSubject = subject
	m.lastBody = htmlBody

	if m.shouldFail {
		return errors.New("mock HTML email send failed")
	}
	return nil
}

// TestHandleSendEmailTask_ProcessTask tests email task processing
func TestHandleSendEmailTask_ProcessTask(t *testing.T) {
	t.Run("process plain text email successfully", func(t *testing.T) {
		mockService := &mockEmailService{}
		handler := tasks.NewHandleSendEmailTask(mockService)

		task, err := tasks.NewSendEmailTask("test@example.com", "Subject", "Body")
		require.NoError(t, err)

		ctx := context.Background()
		err = handler.ProcessTask(ctx, task)

		assert.NoError(t, err, "should process task without error")
		assert.True(t, mockService.sendEmailCalled, "should call SendEmail")
		assert.False(t, mockService.sendHTMLEmailCalled, "should not call SendHTMLEmail")
		assert.Equal(t, "test@example.com", mockService.lastTo)
		assert.Equal(t, "Subject", mockService.lastSubject)
		assert.Equal(t, "Body", mockService.lastBody)
	})

	t.Run("process HTML email with template", func(t *testing.T) {
		mockService := &mockEmailService{}
		handler := tasks.NewHandleSendEmailTask(mockService)

		task, err := tasks.NewSendEmailTaskWithTemplate("test@example.com", "Subject", "template")
		require.NoError(t, err)

		ctx := context.Background()
		err = handler.ProcessTask(ctx, task)

		assert.NoError(t, err)
		assert.False(t, mockService.sendEmailCalled, "should not call SendEmail")
		assert.True(t, mockService.sendHTMLEmailCalled, "should call SendHTMLEmail")
	})

	t.Run("handle email service failure", func(t *testing.T) {
		mockService := &mockEmailService{shouldFail: true}
		handler := tasks.NewHandleSendEmailTask(mockService)

		task, err := tasks.NewSendEmailTask("test@example.com", "Subject", "Body")
		require.NoError(t, err)

		ctx := context.Background()
		err = handler.ProcessTask(ctx, task)

		assert.Error(t, err, "should return error on email send failure")
		assert.Contains(t, err.Error(), "failed to send email")
	})

	t.Run("handle nil email service", func(t *testing.T) {
		handler := tasks.NewHandleSendEmailTask(nil)

		task, err := tasks.NewSendEmailTask("test@example.com", "Subject", "Body")
		require.NoError(t, err)

		ctx := context.Background()
		err = handler.ProcessTask(ctx, task)

		// Should not error but log warning
		assert.NoError(t, err, "should not error with nil service")
	})

	t.Run("handle invalid payload", func(t *testing.T) {
		mockService := &mockEmailService{}
		handler := tasks.NewHandleSendEmailTask(mockService)

		// Create task with invalid JSON payload
		task := asynq.NewTask(tasks.TypeSendEmail, []byte("invalid json"))

		ctx := context.Background()
		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err, "should error on invalid payload")
		assert.Contains(t, err.Error(), "failed to unmarshal payload")
	})
}

// TestSendEmailPayload tests payload structure
func TestSendEmailPayload(t *testing.T) {
	t.Run("marshal and unmarshal payload", func(t *testing.T) {
		payload := tasks.SendEmailPayload{
			To:       "test@example.com",
			Subject:  "Test",
			Body:     "Body",
			Template: "template",
		}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		var unmarshaled tasks.SendEmailPayload
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, payload.To, unmarshaled.To)
		assert.Equal(t, payload.Subject, unmarshaled.Subject)
		assert.Equal(t, payload.Body, unmarshaled.Body)
		assert.Equal(t, payload.Template, unmarshaled.Template)
	})
}

// BenchmarkNewSendEmailTask benchmarks task creation
func BenchmarkNewSendEmailTask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = tasks.NewSendEmailTask("test@example.com", "Subject", "Body")
	}
}

// BenchmarkHandleSendEmailTask_ProcessTask benchmarks task processing
func BenchmarkHandleSendEmailTask_ProcessTask(b *testing.B) {
	mockService := &mockEmailService{}
	handler := tasks.NewHandleSendEmailTask(mockService)

	task, _ := tasks.NewSendEmailTask("test@example.com", "Subject", "Body")
	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = handler.ProcessTask(ctx, task)
	}
}
