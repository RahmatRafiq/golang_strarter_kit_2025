package interfaces

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
)

// QueueClient defines the interface for job queue operations
// This abstraction allows for easier testing and potential swapping of queue backends
type QueueClient interface {
	// Enqueue adds a task to the queue with default options
	Enqueue(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)

	// EnqueueWithPriority enqueues a task to a specific priority queue
	EnqueueWithPriority(ctx context.Context, task *asynq.Task, priority string, opts ...asynq.Option) (*asynq.TaskInfo, error)

	// EnqueueIn schedules a task to be processed after a delay
	EnqueueIn(ctx context.Context, task *asynq.Task, delay time.Duration, opts ...asynq.Option) (*asynq.TaskInfo, error)

	// EnqueueAt schedules a task to be processed at a specific time
	EnqueueAt(ctx context.Context, task *asynq.Task, processAt time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error)

	// Close closes the client connection
	Close() error
}

// TaskHandler defines the interface for task processing
// This allows for dependency injection and easier testing
type TaskHandler interface {
	ProcessTask(ctx context.Context, task *asynq.Task) error
}

// WorkerManager defines the interface for managing workers
// This abstraction allows for different worker implementations
type WorkerManager interface {
	// RegisterHandler registers a task handler for a specific task type
	RegisterHandler(taskType string, handler asynq.Handler)

	// Start starts the worker server
	Start() error

	// Shutdown gracefully shuts down the worker server
	Shutdown()

	// GetMetrics returns current worker metrics (if monitoring enabled)
	GetMetrics() *WorkerMetrics
}

// WorkerMetrics holds worker performance metrics
type WorkerMetrics struct {
	TasksProcessed   int64
	TasksFailed      int64
	TasksRetried     int64
	AverageLatencyMs float64
	QueueDepth       map[string]int64
	ActiveWorkers    int
	LastUpdated      time.Time
}
