package facades

import (
	"sync"

	"golang_starter_kit_2025/app/workers"
	"golang_starter_kit_2025/config"
)

var (
	QueueClient     *workers.Client
	queueClientOnce sync.Once
)

// InitQueue initializes the queue client
func InitQueue() *workers.Client {
	queueClientOnce.Do(func() {
		cfg := config.GetQueueConfig()
		QueueClient = workers.NewClient(cfg)
	})
	return QueueClient
}

// GetQueueClient returns the initialized queue client
func GetQueueClient() *workers.Client {
	if QueueClient == nil {
		return InitQueue()
	}
	return QueueClient
}

// CloseQueue closes the queue client connection
func CloseQueue() error {
	if QueueClient != nil {
		return QueueClient.Close()
	}
	return nil
}
