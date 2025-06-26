package queue

import (
	"fmt"
	"golang_starter_kit_2025/config"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisQueue implements QueueInterface using Redis
type RedisQueue struct {
	defaultQueue string
	prefix       string
}

// NewRedisQueue creates a new Redis queue instance
func NewRedisQueue() *RedisQueue {
	return &RedisQueue{
		defaultQueue: "default",
		prefix:       "queue:",
	}
}

// Push adds a job to the queue
func (rq *RedisQueue) Push(job JobInterface, delay ...time.Duration) error {
	if config.Redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	serialized, err := SerializeJob(job)
	if err != nil {
		return fmt.Errorf("failed to serialize job: %w", err)
	}

	queueName := rq.getQueueName(rq.defaultQueue)

	// If delay is specified, schedule the job
	if len(delay) > 0 && delay[0] > 0 {
		return rq.scheduleJob(queueName, serialized, delay[0])
	}

	// Push to queue immediately
	return config.Redis.LPush(queueName, string(serialized))
}

// Pop retrieves and removes a job from the queue
func (rq *RedisQueue) Pop(queueName string) (*QueueJob, error) {
	if config.Redis == nil {
		return nil, fmt.Errorf("redis not initialized")
	}

	if queueName == "" {
		queueName = rq.defaultQueue
	}

	fullQueueName := rq.getQueueName(queueName)

	// Use blocking pop with timeout
	result, err := config.Redis.BRPop(5*time.Second, fullQueueName)
	if err != nil {
		return nil, err
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("invalid result from queue")
	}

	job, err := DeserializeJob([]byte(result[1]))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize job: %w", err)
	}

	// Generate ID if not exists
	if job.ID == "" {
		job.ID = uuid.New().String()
	}

	return job, nil
}

// Size returns the number of jobs in the queue
func (rq *RedisQueue) Size(queueName string) (int64, error) {
	if config.Redis == nil {
		return 0, fmt.Errorf("redis not initialized")
	}

	if queueName == "" {
		queueName = rq.defaultQueue
	}

	return config.Redis.LLen(rq.getQueueName(queueName))
}

// Clear removes all jobs from the queue
func (rq *RedisQueue) Clear(queueName string) error {
	if config.Redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	if queueName == "" {
		queueName = rq.defaultQueue
	}

	return config.Redis.Del(rq.getQueueName(queueName))
}

// Failed handles failed jobs
func (rq *RedisQueue) Failed(job *QueueJob, err error) error {
	if config.Redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	job.Attempts++

	// If max retries exceeded, move to failed queue
	if job.Attempts >= job.MaxTries {
		failedKey := rq.prefix + "failed"
		serialized, serErr := SerializeJob(&BaseJob{
			Name: job.Name,
			Data: job.Data,
		})
		if serErr != nil {
			return serErr
		}
		return config.Redis.LPush(failedKey, string(serialized))
	}

	// Retry the job
	return rq.Retry(job)
}

// Retry adds a failed job back to the queue
func (rq *RedisQueue) Retry(job *QueueJob) error {
	if config.Redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	serialized, err := SerializeJob(&BaseJob{
		Name: job.Name,
		Data: job.Data,
	})
	if err != nil {
		return err
	}

	queueName := rq.getQueueName(rq.defaultQueue)
	return config.Redis.LPush(queueName, string(serialized))
}

// scheduleJob schedules a job to be executed later
func (rq *RedisQueue) scheduleJob(queueName string, jobData []byte, delay time.Duration) error {
	delayedKey := rq.prefix + "delayed"
	score := float64(time.Now().Add(delay).Unix())

	return config.Redis.Client.ZAdd(config.Redis.Ctx, delayedKey, redis.Z{
		Score:  score,
		Member: string(jobData),
	}).Err()
}

// getQueueName returns the full queue name with prefix
func (rq *RedisQueue) getQueueName(queueName string) string {
	return rq.prefix + queueName
}

// ProcessDelayedJobs moves delayed jobs to the main queue when their time comes
func (rq *RedisQueue) ProcessDelayedJobs() error {
	if config.Redis == nil {
		return fmt.Errorf("redis not initialized")
	}

	delayedKey := rq.prefix + "delayed"
	now := float64(time.Now().Unix())

	// Get jobs that should be processed now
	results, err := config.Redis.Client.ZRangeByScore(config.Redis.Ctx, delayedKey, &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%f", now),
	}).Result()

	if err != nil {
		return err
	}

	// Move each job to the main queue
	for _, jobData := range results {
		// Add to main queue
		queueName := rq.getQueueName(rq.defaultQueue)
		if err := config.Redis.LPush(queueName, jobData); err != nil {
			continue
		}

		// Remove from delayed queue
		config.Redis.Client.ZRem(config.Redis.Ctx, delayedKey, jobData)
	}

	return nil
}
