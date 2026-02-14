package workers

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang_starter_kit_2025/app/workers/interfaces"
	"golang_starter_kit_2025/config"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

// WorkerManager manages the lifecycle of async workers
type WorkerManager struct {
	metricsStarted time.Time
	server         *asynq.Server
	mux            *asynq.ServeMux
	config         *config.QueueConfig
	tasksProcessed atomic.Int64
	tasksFailed    atomic.Int64
	tasksRetried   atomic.Int64
	totalLatencyMs atomic.Int64
	latencyCount   atomic.Int64
	mu             sync.Mutex
	metricsEnabled bool
}

// NewWorkerManager creates a new worker manager instance
func NewWorkerManager(cfg *config.QueueConfig) *WorkerManager {
	wm := &WorkerManager{
		mux:            asynq.NewServeMux(),
		config:         cfg,
		metricsEnabled: true,
		metricsStarted: time.Now(),
	}

	server := asynq.NewServer(
		cfg.GetRedisClientOpt(),
		asynq.Config{
			Concurrency: cfg.Concurrency,
			Queues:      cfg.Queues,
			// Retry configuration
			RetryDelayFunc: asynq.DefaultRetryDelayFunc,
			// Error handler with metrics tracking
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				wm.tasksFailed.Add(1)

				retried, _ := asynq.GetRetryCount(ctx)
				if retried > 0 {
					wm.tasksRetried.Add(1)
				}

				log.Error().
					Err(err).
					Str("task_type", task.Type()).
					Str("task_payload", string(task.Payload())).
					Int("retry_count", retried).
					Msg("Task processing failed")
			}),
			// Logger
			Logger: &AsynqLogger{},
		},
	)

	wm.server = server
	return wm
}

// NewWorkerManagerWithMetrics creates a worker manager with metrics enabled/disabled
func NewWorkerManagerWithMetrics(cfg *config.QueueConfig, enableMetrics bool) *WorkerManager {
	wm := NewWorkerManager(cfg)
	wm.metricsEnabled = enableMetrics
	return wm
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

// GetMetrics returns current worker metrics
func (wm *WorkerManager) GetMetrics() *interfaces.WorkerMetrics {
	if !wm.metricsEnabled {
		return nil
	}

	latencyCount := wm.latencyCount.Load()
	avgLatency := float64(0)
	if latencyCount > 0 {
		avgLatency = float64(wm.totalLatencyMs.Load()) / float64(latencyCount)
	}

	return &interfaces.WorkerMetrics{
		TasksProcessed:   wm.tasksProcessed.Load(),
		TasksFailed:      wm.tasksFailed.Load(),
		TasksRetried:     wm.tasksRetried.Load(),
		AverageLatencyMs: avgLatency,
		ActiveWorkers:    wm.config.Concurrency,
		LastUpdated:      time.Now(),
		QueueDepth:       make(map[string]int64), // Can be populated from Redis if needed
	}
}

// trackTaskCompletion increments the tasks processed counter and records latency
func (wm *WorkerManager) trackTaskCompletion(latencyMs int64) {
	if !wm.metricsEnabled {
		return
	}

	wm.tasksProcessed.Add(1)
	wm.totalLatencyMs.Add(latencyMs)
	wm.latencyCount.Add(1)
}

// WrapHandlerWithMetrics wraps a handler to track metrics
func (wm *WorkerManager) WrapHandlerWithMetrics(handler asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		start := time.Now()

		err := handler.ProcessTask(ctx, task)

		if wm.metricsEnabled {
			latencyMs := time.Since(start).Milliseconds()
			if err == nil {
				wm.trackTaskCompletion(latencyMs)
			}
		}

		return err
	})
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
