package workers

import (
	"sync"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

// TaskRegistry manages task handlers registration
// This pattern allows for auto-discovery and centralized task management
type TaskRegistry struct {
	handlers map[string]asynq.Handler
	mu       sync.RWMutex
}

// NewTaskRegistry creates a new task registry
func NewTaskRegistry() *TaskRegistry {
	return &TaskRegistry{
		handlers: make(map[string]asynq.Handler),
	}
}

// Register adds a task handler to the registry
func (r *TaskRegistry) Register(taskType string, handler asynq.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[taskType]; exists {
		log.Warn().
			Str("task_type", taskType).
			Msg("Task handler already registered, overwriting")
	}

	r.handlers[taskType] = handler

	log.Info().
		Str("task_type", taskType).
		Msg("Task handler registered in registry")
}

// RegisterAll registers all handlers to a worker manager
func (r *TaskRegistry) RegisterAll(wm *WorkerManager) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for taskType, handler := range r.handlers {
		wm.RegisterHandler(taskType, handler)
	}

	log.Info().
		Int("handlers_count", len(r.handlers)).
		Msg("All task handlers registered to worker manager")
}

// Get retrieves a handler by task type
func (r *TaskRegistry) Get(taskType string) (asynq.Handler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, exists := r.handlers[taskType]
	return handler, exists
}

// List returns all registered task types
func (r *TaskRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	taskTypes := make([]string, 0, len(r.handlers))
	for taskType := range r.handlers {
		taskTypes = append(taskTypes, taskType)
	}

	return taskTypes
}

// Count returns the number of registered handlers
func (r *TaskRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.handlers)
}

// GlobalTaskRegistry is the global task registry instance
var (
	globalRegistry     *TaskRegistry
	globalRegistryOnce sync.Once
)

// GetGlobalRegistry returns the global task registry
func GetGlobalRegistry() *TaskRegistry {
	globalRegistryOnce.Do(func() {
		globalRegistry = NewTaskRegistry()
	})
	return globalRegistry
}

// RegisterTask is a convenience function to register a task to the global registry
func RegisterTask(taskType string, handler asynq.Handler) {
	GetGlobalRegistry().Register(taskType, handler)
}
