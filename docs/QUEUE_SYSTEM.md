# Queue System Documentation

## Overview

Production-ready queue system built on [Asynq](https://github.com/hibiken/asynq) with Redis. Modular, scalable, and easy to use.

## Quick Start

### 1. Configuration

```env
QUEUE_ENABLED=true
QUEUE_REDIS_ADDR=localhost:6379
QUEUE_REDIS_PASSWORD=
QUEUE_REDIS_DB=1
QUEUE_CONCURRENCY=10
```

### 2. Enqueue a Task

```go
import (
    "golang_starter_kit_2025/app/workers/tasks"
    "golang_starter_kit_2025/facades"
)

// Get queue client
queueClient := facades.GetQueueClient()

// Create and enqueue task
task, err := tasks.NewSendEmailTask(
    "user@example.com",
    "Welcome!",
    "Welcome to our platform",
)
if err != nil {
    return err
}

// Enqueue with priority
info, err := queueClient.EnqueueWithPriority(ctx, task, "critical")
```

### 3. Create Custom Task

#### Define Task & Payload

```go
// app/workers/tasks/my_task.go
package tasks

import (
    "golang_starter_kit_2025/app/workers"
    "github.com/hibiken/asynq"
)

const TypeMyTask = "my:task"

type MyTaskPayload struct {
    UserID uint   `json:"user_id"`
    Action string `json:"action"`
}

// Using base helper
func NewMyTask(userID uint, action string) (*asynq.Task, error) {
    payload := MyTaskPayload{UserID: userID, Action: action}
    return workers.NewTaskFromPayload(TypeMyTask, payload)
}
```

#### Create Handler

```go
// app/workers/tasks/my_task.go (continued)
package tasks

import (
    "context"
    "golang_starter_kit_2025/app/workers"
)

type HandleMyTask struct {
    myService MyService
}

func NewHandleMyTask(service MyService) *HandleMyTask {
    return &HandleMyTask{myService: service}
}

func (h *HandleMyTask) ProcessTask(ctx context.Context, task *asynq.Task) error {
    var payload MyTaskPayload
    if err := workers.UnmarshalTaskPayload(task, &payload); err != nil {
        return err
    }

    // Process task
    return h.myService.DoSomething(payload.UserID, payload.Action)
}
```

#### Register Handler

```go
// app/providers/queue_provider.go
package providers

import (
    "golang_starter_kit_2025/app/workers"
    "golang_starter_kit_2025/app/workers/tasks"
)

func RegisterQueueHandlers() {
    // Initialize services
    myService := services.NewMyService()

    // Register to global registry
    workers.RegisterTask(tasks.TypeMyTask, tasks.NewHandleMyTask(myService))
}
```

Handlers auto-register via bootstrap:

```go
// bootstrap/main.go
if helpers.GetEnv("QUEUE_ENABLED", "false") == "true" {
    workerManager = workers.NewWorkerManager(queueCfg)

    // Auto-register all handlers
    providers.RegisterQueueHandlers()
    registry := workers.GetGlobalRegistry()
    registry.RegisterAll(workerManager)

    facades.SetWorkerManager(workerManager)
    go workerManager.Start()
}
```

## Priority Queues

| Queue    | Weight | Use Case                          |
|----------|--------|-----------------------------------|
| critical | 60%    | Password resets, OTP              |
| default  | 30%    | Regular emails, notifications     |
| low      | 10%    | Analytics, cleanup                |

```go
// Critical priority
queueClient.EnqueueWithPriority(ctx, task, "critical")

// Default priority
queueClient.EnqueueWithPriority(ctx, task, "default")

// Low priority
queueClient.EnqueueWithPriority(ctx, task, "low")
```

## Scheduling

### Delay Execution

```go
delay := 5 * time.Minute
info, err := queueClient.EnqueueIn(ctx, task, delay)
```

### Schedule at Specific Time

```go
processAt := time.Now().Add(24 * time.Hour)
info, err := queueClient.EnqueueAt(ctx, task, processAt)
```

## Retry Configuration

```go
// Custom retry count
info, err := queueClient.Enqueue(ctx, task, asynq.MaxRetry(5))

// Disable retry
info, err := queueClient.Enqueue(ctx, task, asynq.MaxRetry(0))

// With timeout
info, err := queueClient.Enqueue(ctx, task, asynq.Timeout(30 * time.Second))
```

Default: 25 retries with exponential backoff

## Monitoring

### Queue Health Endpoint

```bash
curl http://localhost:8080/health/queue
```

Response:

```json
{
  "status": "healthy",
  "tasks_processed": 1234,
  "tasks_failed": 12,
  "tasks_retried": 45,
  "average_latency_ms": 125.5,
  "active_workers": 10,
  "failure_rate": 0.97,
  "last_updated": "2026-02-14T15:30:00Z",
  "queue_depth": {
    "critical": 5,
    "default": 12,
    "low": 3
  }
}
```

Status: `healthy` (< 50% failure) | `degraded` (> 50% failure)

### Get Metrics Programmatically

```go
workerManager := facades.GetWorkerManager()
metrics := workerManager.GetMetrics()

if metrics != nil {
    fmt.Printf("Tasks Processed: %d\n", metrics.TasksProcessed)
    fmt.Printf("Average Latency: %.2fms\n", metrics.AverageLatencyMs)
}
```

## Base Helpers

Reduce boilerplate with base helpers:

```go
// Create task
payload := MyPayload{UserID: 123}
task, err := workers.NewTaskFromPayload("my:task", payload)

// Process task
var payload MyPayload
err := workers.UnmarshalTaskPayload(task, &payload)
```

## Best Practices

### Task Design

✅ DO:
- Keep tasks small and focused
- Make tasks idempotent (safe to retry)
- Use clear task names
- Include all data in payload

❌ DON'T:
- Pass large objects (use IDs instead)
- Depend on global state
- Use for synchronous operations
- Include sensitive data

### Error Handling

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

### Task Deduplication

```go
// Prevent duplicates within 24 hours
info, err := queueClient.Enqueue(ctx, task,
    asynq.TaskID("send-email-" + userEmail),
    asynq.Unique(24 * time.Hour),
)
```

## Architecture

```
Application Layer (Services)
    ↓ uses
Queue Client (Producer)
    ↓
Redis (Message Broker)
    ↓
Worker Manager (Consumer)
    ↓ executes
Task Handlers
```

### Components

- **Client**: Enqueue tasks with priority/scheduling
- **WorkerManager**: Lifecycle management, metrics tracking
- **TaskRegistry**: Auto-discovery, centralized management
- **BaseHelpers**: DRY for task creation/processing
- **Interfaces**: Modular, testable dependencies

## Troubleshooting

### Redis Connection

```bash
# Check Redis
redis-cli ping  # Should return: PONG

# Check connection
redis-cli -h localhost -p 6379 -n 1 info
```

### Tasks Not Processing

1. Check `QUEUE_ENABLED=true` in `.env`
2. Verify workers started in `bootstrap/main.go`
3. Check handler registered in `providers/queue_provider.go`
4. View logs for errors

### Performance

- High memory: Reduce `QUEUE_CONCURRENCY`
- Slow processing: Increase `QUEUE_CONCURRENCY`, optimize handlers
- Check external service response times

## Testing

```bash
# Unit tests
go test ./tests/unit/... -v -run "Queue|Task"

# Integration tests (requires Redis)
go test ./tests/integration/... -v -run "QueueFlow"

# App package only
go test ./app/... -v
```

## References

- [Asynq Documentation](https://github.com/hibiken/asynq)
- [Asynq Best Practices](https://github.com/hibiken/asynq/wiki/Best-Practices)
