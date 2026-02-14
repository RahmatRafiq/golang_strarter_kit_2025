package services_test

import (
	"context"
	"testing"
	"time"

	"golang_starter_kit_2025/app/workers"
	"golang_starter_kit_2025/config"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Load .env.test for all tests in this package
	_ = godotenv.Load("../../.env.test")
}

// TestQueueClient_Enqueue tests basic enqueue functionality
func TestQueueClient_Enqueue(t *testing.T) {
	// Use centralized Redis config for testing
	redisConfig := config.GetTestRedisConfig()

	cfg := &config.QueueConfig{
		RedisAddr:     redisConfig.Addr,
		RedisPassword: redisConfig.Password,
		RedisDB:       redisConfig.DB,
		Concurrency:   1,
		Queues: map[string]int{
			"default": 1,
		},
	}

	client := workers.NewClient(cfg)
	defer client.Close()

	ctx := context.Background()

	t.Run("enqueue task successfully", func(t *testing.T) {
		task := asynq.NewTask("test:task", []byte(`{"test": "data"}`))

		info, err := client.Enqueue(ctx, task)

		require.NoError(t, err, "should enqueue task without error")
		assert.NotEmpty(t, info.ID, "task should have an ID")
		assert.Equal(t, "test:task", info.Type, "task type should match")
		assert.Equal(t, "default", info.Queue, "task should be in default queue")
	})

	t.Run("enqueue with priority", func(t *testing.T) {
		task := asynq.NewTask("test:priority", []byte(`{"priority": "high"}`))

		info, err := client.EnqueueWithPriority(ctx, task, "critical")

		require.NoError(t, err)
		assert.Equal(t, "critical", info.Queue, "task should be in critical queue")
	})

	t.Run("enqueue with delay", func(t *testing.T) {
		task := asynq.NewTask("test:delayed", []byte(`{"delayed": true}`))
		delay := 5 * time.Second

		info, err := client.EnqueueIn(ctx, task, delay)

		require.NoError(t, err)
		assert.NotEmpty(t, info.ID)
		// Check that process time is in the future
		assert.True(t, info.NextProcessAt.After(time.Now()))
	})

	t.Run("enqueue at specific time", func(t *testing.T) {
		task := asynq.NewTask("test:scheduled", []byte(`{"scheduled": true}`))
		processAt := time.Now().Add(10 * time.Second)

		info, err := client.EnqueueAt(ctx, task, processAt)

		require.NoError(t, err)
		assert.NotEmpty(t, info.ID)
		// Check that process time matches (with some tolerance)
		diff := info.NextProcessAt.Sub(processAt).Abs()
		assert.Less(t, diff, time.Second, "process time should match scheduled time")
	})

	t.Run("enqueue with max retry", func(t *testing.T) {
		task := asynq.NewTask("test:retry", []byte(`{"retry": true}`))

		info, err := client.Enqueue(ctx, task, asynq.MaxRetry(3))

		require.NoError(t, err)
		assert.Equal(t, 3, info.MaxRetry, "max retry should be set")
	})
}

// TestQueueClient_Close tests client cleanup
func TestQueueClient_Close(t *testing.T) {
	redisConfig := config.GetTestRedisConfig()

	cfg := &config.QueueConfig{
		RedisAddr:     redisConfig.Addr,
		RedisPassword: redisConfig.Password,
		RedisDB:       redisConfig.DB,
		Concurrency:   1,
		Queues: map[string]int{
			"default": 1,
		},
	}

	client := workers.NewClient(cfg)

	err := client.Close()
	assert.NoError(t, err, "should close client without error")
}

// BenchmarkQueueClient_Enqueue benchmarks enqueue performance
func BenchmarkQueueClient_Enqueue(b *testing.B) {
	redisConfig := config.GetTestRedisConfig()

	cfg := &config.QueueConfig{
		RedisAddr:     redisConfig.Addr,
		RedisPassword: redisConfig.Password,
		RedisDB:       redisConfig.DB,
		Concurrency:   10,
		Queues: map[string]int{
			"default": 1,
		},
	}

	client := workers.NewClient(cfg)
	defer client.Close()

	ctx := context.Background()
	task := asynq.NewTask("benchmark:task", []byte(`{"benchmark": true}`))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = client.Enqueue(ctx, task)
	}
}
