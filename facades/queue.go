package facades

import (
	"sync"

	"golang_starter_kit_2025/app/workers"
	"golang_starter_kit_2025/config"
)

var (
	QueueClient        *workers.Client
	queueClientOnce    sync.Once
	WorkerManager      *workers.WorkerManager
	workerManagerMutex sync.Mutex
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

// SetWorkerManager sets the global worker manager instance
// This is called from bootstrap during initialization
func SetWorkerManager(wm *workers.WorkerManager) {
	workerManagerMutex.Lock()
	defer workerManagerMutex.Unlock()
	WorkerManager = wm
}

// GetWorkerManager returns the global worker manager instance
func GetWorkerManager() *workers.WorkerManager {
	workerManagerMutex.Lock()
	defer workerManagerMutex.Unlock()
	return WorkerManager
}
