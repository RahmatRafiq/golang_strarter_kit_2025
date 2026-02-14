package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang_starter_kit_2025/app/workers"
	"golang_starter_kit_2025/config"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

// TestQueue wraps queue client and worker manager for testing
type TestQueue struct {
	Client        *workers.Client
	WorkerManager *workers.WorkerManager
	Config        *config.QueueConfig
	t             *testing.T
}

// SetupTestQueue initializes a test queue with Redis
// Similar to SetupTestDB for database testing
func SetupTestQueue(t *testing.T) *TestQueue {
	t.Helper()

	// Load .env.test for test environment
	if err := godotenv.Load("../../.env.test"); err != nil {
		t.Logf("Warning: .env.test not loaded: %v", err)
	}

	// Use centralized Redis config for testing
	redisConfig := config.GetTestRedisConfig()

	cfg := &config.QueueConfig{
		RedisAddr:     redisConfig.Addr,
		RedisPassword: redisConfig.Password,
		RedisDB:       redisConfig.DB,
		Concurrency:   2, // Low concurrency for tests
		Queues: map[string]int{
			"critical": 2,
			"default":  1,
		},
	}

	client := workers.NewClient(cfg)
	workerManager := workers.NewWorkerManager(cfg)

	testQueue := &TestQueue{
		Client:        client,
		WorkerManager: workerManager,
		Config:        cfg,
		t:             t,
	}

	t.Log("Test queue initialized with Redis DB 15")

	return testQueue
}

// CleanupTestQueue stops workers and closes client
func CleanupTestQueue(t *testing.T, tq *TestQueue) {
	t.Helper()

	if tq == nil {
		return
	}

	if tq.WorkerManager != nil {
		tq.WorkerManager.Shutdown()
	}

	if tq.Client != nil {
		if err := tq.Client.Close(); err != nil {
			t.Logf("Warning: Failed to close queue client: %v", err)
		}
	}

	t.Log("Test queue cleaned up")
}

// EnqueueTestTask enqueues a task and returns task info
func EnqueueTestTask(t *testing.T, tq *TestQueue, task *asynq.Task, opts ...asynq.Option) *asynq.TaskInfo {
	t.Helper()

	ctx := context.Background()
	info, err := tq.Client.Enqueue(ctx, task, opts...)
	if err != nil {
		t.Fatalf("Failed to enqueue test task: %v", err)
	}

	t.Logf("Enqueued task: ID=%s, Type=%s, Queue=%s", info.ID, info.Type, info.Queue)
	return info
}

// EnqueueTestTaskWithPriority enqueues a task to a specific priority queue
func EnqueueTestTaskWithPriority(t *testing.T, tq *TestQueue, task *asynq.Task, priority string, opts ...asynq.Option) *asynq.TaskInfo {
	t.Helper()

	ctx := context.Background()
	info, err := tq.Client.EnqueueWithPriority(ctx, task, priority, opts...)
	if err != nil {
		t.Fatalf("Failed to enqueue task with priority: %v", err)
	}

	t.Logf("Enqueued task with priority %s: ID=%s, Type=%s", priority, info.ID, info.Type)
	return info
}

// WaitForTaskCompletion waits for a task to complete with timeout
// Returns true if task completed, false if timeout reached
func WaitForTaskCompletion(t *testing.T, checkFunc func() bool, timeout time.Duration) bool {
	t.Helper()

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if checkFunc() {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}

	return false
}

// RegisterTestHandler registers a handler to the test worker manager
func RegisterTestHandler(t *testing.T, tq *TestQueue, taskType string, handler asynq.Handler) {
	t.Helper()

	tq.WorkerManager.RegisterHandler(taskType, handler)
	t.Logf("Registered test handler for task type: %s", taskType)
}

// StartTestWorkers starts the worker manager in background
func StartTestWorkers(t *testing.T, tq *TestQueue) {
	t.Helper()

	go func() {
		if err := tq.WorkerManager.Start(); err != nil {
			t.Logf("Worker manager error: %v", err)
		}
	}()

	// Give workers time to start
	time.Sleep(500 * time.Millisecond)

	t.Log("Test workers started")
}

// CreateTestTask creates a simple test task
func CreateTestTask(taskType string, payload []byte) *asynq.Task {
	if payload == nil {
		payload = []byte(fmt.Sprintf(`{"test": true, "created_at": "%s"}`, time.Now().Format(time.RFC3339)))
	}
	return asynq.NewTask(taskType, payload)
}

// MockTaskHandler is a simple handler for testing that tracks calls
type MockTaskHandler struct {
	ErrorMsg     string
	LastPayload  []byte
	ProcessCount int
	ProcessDelay time.Duration
	ShouldError  bool
}

// ProcessTask implements asynq.Handler
func (h *MockTaskHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	h.ProcessCount++
	h.LastPayload = task.Payload()

	if h.ProcessDelay > 0 {
		time.Sleep(h.ProcessDelay)
	}

	if h.ShouldError {
		return fmt.Errorf("%s", h.ErrorMsg)
	}

	return nil
}

// Reset resets the handler state
func (h *MockTaskHandler) Reset() {
	h.ProcessCount = 0
	h.LastPayload = nil
}

// WasProcessed checks if handler was called
func (h *MockTaskHandler) WasProcessed() bool {
	return h.ProcessCount > 0
}
