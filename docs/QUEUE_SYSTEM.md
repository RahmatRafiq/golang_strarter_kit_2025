# Queue System Documentation

## 📋 Overview

This starter kit includes a production-ready queue system built on [Asynq](https://github.com/hibiken/asynq) with Redis as the message broker. The system is designed to be **modular**, **scalable**, **testable**, and **easy to use**.

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Application Layer                        │
│  (Services: PasswordResetService, EmailVerificationService)  │
└────────────────────┬────────────────────────────────────────┘
                     │ uses
                     ↓
┌─────────────────────────────────────────────────────────────┐
│                   Queue Client (Producer)                    │
│  - Enqueue()                                                 │
│  - EnqueueWithPriority()                                     │
│  - EnqueueIn() / EnqueueAt()                                 │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ↓
┌─────────────────────────────────────────────────────────────┐
│                      Redis (Broker)                          │
│  - critical queue (60% workers)                              │
│  - default queue  (30% workers)                              │
│  - low queue      (10% workers)                              │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ↓
┌─────────────────────────────────────────────────────────────┐
│                 Worker Manager (Consumer)                    │
│  - Task Registry                                             │
│  - Metrics Tracking                                          │
│  - Error Handling                                            │
└────────────────────┬────────────────────────────────────────┘
                     │ executes
                     ↓
┌─────────────────────────────────────────────────────────────┐
│                      Task Handlers                           │
│  - HandleSendEmailTask                                       │
│  - ... (add more tasks here)                                 │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### 1. Configuration

Add to your `.env`:

```env
# Queue Configuration
QUEUE_ENABLED=true
QUEUE_REDIS_ADDR=localhost:6379
QUEUE_REDIS_PASSWORD=
QUEUE_REDIS_DB=1
QUEUE_CONCURRENCY=10
```

### 2. Enqueue a Task

```go
import (
    "context"
    "golang_starter_kit_2025/app/workers/tasks"
    "golang_starter_kit_2025/facades"
)

// Get queue client
queueClient := facades.GetQueueClient()

// Create task
task, err := tasks.NewSendEmailTask(
    "user@example.com",
    "Welcome!",
    "Welcome to our platform",
)
if err != nil {
    return err
}

// Enqueue with priority
ctx := context.Background()
info, err := queueClient.EnqueueWithPriority(ctx, task, "critical")
```

### 3. Create Custom Task

#### Step 1: Define Task Type and Payload

```go
// app/workers/tasks/process_file_task.go
package tasks

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/hibiken/asynq"
)

const TypeProcessFile = "file:process"

type ProcessFilePayload struct {
    FilePath string `json:"file_path"`
    UserID   uint   `json:"user_id"`
}

func NewProcessFileTask(filePath string, userID uint) (*asynq.Task, error) {
    payload := ProcessFilePayload{
        FilePath: filePath,
        UserID:   userID,
    }

    data, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal payload: %w", err)
    }

    return asynq.NewTask(TypeProcessFile, data), nil
}
```

#### Step 2: Create Task Handler

```go
// app/workers/tasks/process_file_task.go (continued)

type HandleProcessFileTask struct {
    fileService interface {
        ProcessFile(path string, userID uint) error
    }
}

func NewHandleProcessFileTask(fs interface {
    ProcessFile(path string, userID uint) error
}) *HandleProcessFileTask {
    return &HandleProcessFileTask{
        fileService: fs,
    }
}

func (h *HandleProcessFileTask) ProcessTask(ctx context.Context, task *asynq.Task) error {
    var payload ProcessFilePayload
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return fmt.Errorf("failed to unmarshal payload: %w", err)
    }

    // Process file
    if err := h.fileService.ProcessFile(payload.FilePath, payload.UserID); err != nil {
        return fmt.Errorf("failed to process file: %w", err)
    }

    return nil
}
```

#### Step 3: Register Handler

```go
// bootstrap/main.go (in Init function)

if helpers.GetEnv("QUEUE_ENABLED", "false") == "true" {
    queueCfg := config.GetQueueConfig()
    workerManager = workers.NewWorkerManager(queueCfg)

    // Register email handler
    emailService := services.NewEmailService()
    workerManager.RegisterHandler(tasks.TypeSendEmail, tasks.NewHandleSendEmailTask(emailService))

    // Register file processing handler
    fileService := services.NewFileService()
    workerManager.RegisterHandler(tasks.TypeProcessFile, tasks.NewHandleProcessFileTask(fileService))

    // Start workers
    go func() {
        if err := workerManager.Start(); err != nil {
            log.Fatal().Err(err).Msg("Failed to start worker manager")
        }
    }()
}
```

## 📊 Priority Queues

Tasks can be enqueued to different priority queues:

| Queue    | Weight | Use Case                                    |
|----------|--------|---------------------------------------------|
| critical | 60%    | Password resets, OTP, critical emails       |
| default  | 30%    | Regular emails, notifications               |
| low      | 10%    | Analytics, cleanup, non-urgent tasks        |

```go
// Critical priority
queueClient.EnqueueWithPriority(ctx, task, "critical")

// Default priority
queueClient.EnqueueWithPriority(ctx, task, "default")

// Low priority
queueClient.EnqueueWithPriority(ctx, task, "low")
```

## ⏰ Scheduling Tasks

### Delay Execution

```go
// Process after 5 minutes
delay := 5 * time.Minute
info, err := queueClient.EnqueueIn(ctx, task, delay)
```

### Schedule at Specific Time

```go
// Process tomorrow at 9 AM
processAt := time.Now().Add(24 * time.Hour).
    Truncate(24 * time.Hour).
    Add(9 * time.Hour)

info, err := queueClient.EnqueueAt(ctx, task, processAt)
```

## 🔄 Retry Configuration

```go
import "github.com/hibiken/asynq"

// Custom retry count
info, err := queueClient.Enqueue(ctx, task,
    asynq.MaxRetry(5),
)

// Disable retry
info, err := queueClient.Enqueue(ctx, task,
    asynq.MaxRetry(0),
)

// With timeout
info, err := queueClient.Enqueue(ctx, task,
    asynq.Timeout(30 * time.Second),
)
```

Default retry behavior:
- Max retries: 25
- Backoff: Exponential (seconds → minutes → hours)
- Retry delay function: `asynq.DefaultRetryDelayFunc`

## 📈 Monitoring & Metrics

### Get Worker Metrics

```go
metrics := workerManager.GetMetrics()
if metrics != nil {
    fmt.Printf("Tasks Processed: %d\n", metrics.TasksProcessed)
    fmt.Printf("Tasks Failed: %d\n", metrics.TasksFailed)
    fmt.Printf("Tasks Retried: %d\n", metrics.TasksRetried)
    fmt.Printf("Average Latency: %.2fms\n", metrics.AverageLatencyMs)
    fmt.Printf("Active Workers: %d\n", metrics.ActiveWorkers)
}
```

### Metrics Available

- `TasksProcessed`: Total successfully processed tasks
- `TasksFailed`: Total failed tasks
- `TasksRetried`: Total retry attempts
- `AverageLatencyMs`: Average task processing time
- `ActiveWorkers`: Number of concurrent workers
- `LastUpdated`: Last metrics update timestamp

## 🧪 Testing

### Unit Tests

```go
// tests/unit/my_task_test.go
package services_test

import (
    "context"
    "testing"
    "golang_starter_kit_2025/app/workers/tasks"
    "github.com/stretchr/testify/assert"
)

func TestMyTask_Creation(t *testing.T) {
    task, err := tasks.NewMyTask("param1", "param2")

    assert.NoError(t, err)
    assert.Equal(t, tasks.TypeMyTask, task.Type())
}

func TestMyTaskHandler_Process(t *testing.T) {
    // Create mock service
    mockService := &mockMyService{}
    handler := tasks.NewHandleMyTask(mockService)

    // Create task
    task, _ := tasks.NewMyTask("test", "data")

    // Process
    ctx := context.Background()
    err := handler.ProcessTask(ctx, task)

    assert.NoError(t, err)
    assert.True(t, mockService.called)
}
```

### Integration Tests

```go
// tests/integration/queue_flow_test.go
package integration

import (
    "testing"
    "time"
    testhelpers "golang_starter_kit_2025/tests/helpers"
)

func TestQueueFlow_MyTask(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Setup test queue
    tq := testhelpers.SetupTestQueue(t)
    defer testhelpers.CleanupTestQueue(t, tq)

    // Register handler
    mockHandler := &testhelpers.MockTaskHandler{}
    testhelpers.RegisterTestHandler(t, tq, "my:task", mockHandler)

    // Start workers
    testhelpers.StartTestWorkers(t, tq)

    // Enqueue task
    task := testhelpers.CreateTestTask("my:task", nil)
    testhelpers.EnqueueTestTask(t, tq, task)

    // Wait for completion
    completed := testhelpers.WaitForTaskCompletion(t,
        mockHandler.WasProcessed, 5*time.Second)

    assert.True(t, completed)
}
```

### Run Tests

```bash
# All queue tests
go test ./tests/unit/... -v -run "Queue|Task|Registry"

# Integration tests
go test ./tests/integration/... -v -run "QueueFlow"

# With Redis requirement
go test ./tests/integration/... -v

# Benchmarks
go test ./tests/unit/... -bench=. -benchmem
```

## 🔧 Advanced Features

### Task Deduplication

```go
// Prevent duplicate tasks within 24 hours
info, err := queueClient.Enqueue(ctx, task,
    asynq.TaskID("send-email-" + userEmail),
    asynq.Unique(24 * time.Hour),
)
```

### Custom Queue Config

```go
// config/queue_config.go
Queues: map[string]int{
    "critical":    6,  // 60% workers
    "high":        4,  // 40% workers
    "default":     3,  // 30% workers
    "low":         2,  // 20% workers
    "background":  1,  // 10% workers
}
```

### Task Context

```go
func (h *Handler) ProcessTask(ctx context.Context, task *asynq.Task) error {
    // Get retry count
    retried, _ := asynq.GetRetryCount(ctx)

    // Get task ID
    taskID, _ := asynq.GetTaskID(ctx)

    // Get queue name
    queueName, _ := asynq.GetQueueName(ctx)

    // Use context for cancellation
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Process task
    }

    return nil
}
```

## 🎯 Best Practices

### 1. Task Design

✅ **DO:**
- Keep tasks small and focused
- Make tasks idempotent (safe to retry)
- Use clear, descriptive task names
- Include all necessary data in payload

❌ **DON'T:**
- Pass large objects in payload (use IDs instead)
- Make tasks dependent on global state
- Use tasks for synchronous operations
- Include sensitive data in task payload

### 2. Error Handling

```go
func (h *Handler) ProcessTask(ctx context.Context, task *asynq.Task) error {
    // Permanent errors (don't retry)
    if validation.Failed(payload) {
        return fmt.Errorf("validation failed: %w", asynq.SkipRetry)
    }

    // Transient errors (retry)
    if err := externalAPI.Call(); err != nil {
        return fmt.Errorf("API call failed (will retry): %w", err)
    }

    return nil
}
```

### 3. Monitoring

```go
// Log important events
log.Info().
    Str("task_id", taskID).
    Str("task_type", task.Type()).
    Int("retry_count", retried).
    Msg("Processing task")

// Track metrics
metrics := workerManager.GetMetrics()
// Export to monitoring system (Prometheus, Datadog, etc)
```

### 4. Testing

- Write unit tests for task creation and handler logic
- Write integration tests for end-to-end flow
- Use `testhelpers.SetupTestQueue` for consistent test setup
- Mock external dependencies (email service, APIs, etc)
- Test retry behavior and error scenarios

## 📚 References

- [Asynq Documentation](https://github.com/hibiken/asynq)
- [Asynq Best Practices](https://github.com/hibiken/asynq/wiki/Best-Practices)
- [Testing Guide](https://github.com/hibiken/asynq/wiki/Testing)
- [Task Lifecycle](https://github.com/hibiken/asynq/wiki/Task-Lifecycle)

## 🆘 Troubleshooting

### Redis Connection Issues

```bash
# Check Redis is running
redis-cli ping
# Should return: PONG

# Check connection
redis-cli -h localhost -p 6379 -n 1 info
```

### Tasks Not Processing

1. Check `QUEUE_ENABLED=true` in `.env`
2. Verify workers are started in `bootstrap/main.go`
3. Check task handler is registered
4. Check Redis connection
5. View logs for errors

### High Memory Usage

- Reduce `QUEUE_CONCURRENCY`
- Check for memory leaks in task handlers
- Monitor queue depth (too many pending tasks)

### Slow Processing

- Increase `QUEUE_CONCURRENCY`
- Optimize task handler code
- Check external service response times
- Consider task batching

## 🔄 Migration from Old System

If you're migrating from a synchronous system:

1. Identify long-running operations
2. Extract to task handlers
3. Replace direct calls with enqueue
4. Add proper error handling
5. Test thoroughly
6. Monitor metrics
7. Adjust concurrency as needed

---

**Need help?** Open an issue on GitHub or check the [main README](../README.md).
