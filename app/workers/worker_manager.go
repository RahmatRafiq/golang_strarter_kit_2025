package workers

import (
	"context"
	"fmt"
	"sync"

	"golang_starter_kit_2025/config"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

// WorkerManager manages the lifecycle of async workers
type WorkerManager struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	config *config.QueueConfig
	mu     sync.Mutex
}

// NewWorkerManager creates a new worker manager instance
func NewWorkerManager(cfg *config.QueueConfig) *WorkerManager {
	server := asynq.NewServer(
		cfg.GetRedisClientOpt(),
		asynq.Config{
			Concurrency: cfg.Concurrency,
			Queues:      cfg.Queues,
			// Retry configuration
			RetryDelayFunc: asynq.DefaultRetryDelayFunc,
			// Error handler
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("task_type", task.Type()).
					Str("task_payload", string(task.Payload())).
					Msg("Task processing failed")
			}),
			// Logger
			Logger: &AsynqLogger{},
		},
	)

	mux := asynq.NewServeMux()

	return &WorkerManager{
		server: server,
		mux:    mux,
		config: cfg,
	}
}

// RegisterHandler registers a task handler for a specific task type
func (wm *WorkerManager) RegisterHandler(taskType string, handler asynq.Handler) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.mux.Handle(taskType, handler)
	log.Info().
		Str("task_type", taskType).
		Msg("Task handler registered")
}

// Start starts the worker server
func (wm *WorkerManager) Start() error {
	log.Info().
		Int("concurrency", wm.config.Concurrency).
		Interface("queues", wm.config.Queues).
		Msg("Starting worker manager")

	if err := wm.server.Start(wm.mux); err != nil {
		return fmt.Errorf("failed to start worker server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the worker server
func (wm *WorkerManager) Shutdown() {
	log.Info().Msg("Shutting down worker manager...")
	wm.server.Shutdown()
	log.Info().Msg("Worker manager shutdown complete")
}

// AsynqLogger adapts zerolog to asynq.Logger interface
type AsynqLogger struct{}

func (l *AsynqLogger) Debug(args ...any) {
	log.Debug().Msg(fmt.Sprint(args...))
}

func (l *AsynqLogger) Info(args ...any) {
	log.Info().Msg(fmt.Sprint(args...))
}

func (l *AsynqLogger) Warn(args ...any) {
	log.Warn().Msg(fmt.Sprint(args...))
}

func (l *AsynqLogger) Error(args ...any) {
	log.Error().Msg(fmt.Sprint(args...))
}

func (l *AsynqLogger) Fatal(args ...any) {
	log.Fatal().Msg(fmt.Sprint(args...))
}
