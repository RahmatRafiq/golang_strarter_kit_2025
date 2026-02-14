package providers

import (
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/app/workers"
	"golang_starter_kit_2025/app/workers/tasks"

	"github.com/rs/zerolog/log"
)

// RegisterQueueHandlers registers all task handlers to the global registry
// This provides auto-discovery and centralized task management
func RegisterQueueHandlers() {
	log.Info().Msg("Registering queue handlers...")

	// Initialize services
	emailService := services.NewEmailService()

	// Register task handlers to global registry
	workers.RegisterTask(tasks.TypeSendEmail, tasks.NewHandleSendEmailTask(emailService))

	// Add more task handlers here as your application grows
	// Example:
	// fileService := services.NewFileService()
	// workers.RegisterTask(tasks.TypeProcessFile, tasks.NewHandleProcessFileTask(fileService))

	log.Info().
		Int("handlers_count", workers.GetGlobalRegistry().Count()).
		Msg("Queue handlers registered successfully")
}
