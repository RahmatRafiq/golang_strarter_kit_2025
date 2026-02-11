package workers

import (
	"context"
	"fmt"
	"time"

	"golang_starter_kit_2025/config"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

// Client is a wrapper around asynq.Client for enqueueing tasks
type Client struct {
	client *asynq.Client
}

// NewClient creates a new job client
func NewClient(cfg *config.QueueConfig) *Client {
	client := asynq.NewClient(cfg.GetRedisClientOpt())
	return &Client{client: client}
}

// Enqueue adds a task to the queue with default options
func (c *Client) Enqueue(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	info, err := c.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		log.Error().
			Err(err).
			Str("task_type", task.Type()).
			Msg("Failed to enqueue task")
		return nil, fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("task_id", info.ID).
		Str("task_type", task.Type()).
		Str("queue", info.Queue).
		Msg("Task enqueued successfully")

	return info, nil
}

// EnqueueWithPriority enqueues a task to a specific priority queue
func (c *Client) EnqueueWithPriority(ctx context.Context, task *asynq.Task, priority string, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.Queue(priority))
	return c.Enqueue(ctx, task, opts...)
}

// EnqueueIn schedules a task to be processed after a delay
func (c *Client) EnqueueIn(ctx context.Context, task *asynq.Task, delay time.Duration, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.ProcessIn(delay))
	return c.Enqueue(ctx, task, opts...)
}

// EnqueueAt schedules a task to be processed at a specific time
func (c *Client) EnqueueAt(ctx context.Context, task *asynq.Task, processAt time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.ProcessAt(processAt))
	return c.Enqueue(ctx, task, opts...)
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.client.Close()
}
