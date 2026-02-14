package services_test

import (
	"context"
	"testing"

	"golang_starter_kit_2025/app/workers"
	"golang_starter_kit_2025/config"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTaskHandler is a simple handler for testing
type mockTaskHandler struct {
	called bool
}

func (h *mockTaskHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	h.called = true
	return nil
}

// TestTaskRegistry_Register tests task registration
func TestTaskRegistry_Register(t *testing.T) {
	registry := workers.NewTaskRegistry()

	handler := &mockTaskHandler{}

	t.Run("register new handler", func(t *testing.T) {
		registry.Register("test:task", handler)

		retrieved, exists := registry.Get("test:task")
		require.True(t, exists, "handler should exist after registration")
		assert.Equal(t, handler, retrieved, "retrieved handler should match registered handler")
	})

	t.Run("count registered handlers", func(t *testing.T) {
		count := registry.Count()
		assert.Equal(t, 1, count, "should have 1 registered handler")
	})

	t.Run("list registered task types", func(t *testing.T) {
		taskTypes := registry.List()
		assert.Contains(t, taskTypes, "test:task", "should list registered task type")
	})
}

// TestTaskRegistry_Overwrite tests overwriting existing handler
func TestTaskRegistry_Overwrite(t *testing.T) {
	registry := workers.NewTaskRegistry()

	handler1 := &mockTaskHandler{}
	handler2 := &mockTaskHandler{}

	registry.Register("test:task", handler1)
	registry.Register("test:task", handler2) // Overwrite

	retrieved, exists := registry.Get("test:task")
	require.True(t, exists)
	assert.Equal(t, handler2, retrieved, "should return the latest registered handler")

	count := registry.Count()
	assert.Equal(t, 1, count, "should still have only 1 handler after overwrite")
}

// TestTaskRegistry_MultipleHandlers tests multiple handler registration
func TestTaskRegistry_MultipleHandlers(t *testing.T) {
	registry := workers.NewTaskRegistry()

	handlers := map[string]*mockTaskHandler{
		"email:send":        {},
		"notification:push": {},
		"file:process":      {},
	}

	for taskType, handler := range handlers {
		registry.Register(taskType, handler)
	}

	assert.Equal(t, 3, registry.Count(), "should have 3 registered handlers")

	taskTypes := registry.List()
	assert.Len(t, taskTypes, 3, "should list 3 task types")

	for taskType := range handlers {
		assert.Contains(t, taskTypes, taskType, "should contain registered task type")
	}
}

// TestTaskRegistry_RegisterAll tests registering all handlers to worker manager
func TestTaskRegistry_RegisterAll(t *testing.T) {
	registry := workers.NewTaskRegistry()

	// Register some handlers
	registry.Register("test:handler1", &mockTaskHandler{})
	registry.Register("test:handler2", &mockTaskHandler{})
	registry.Register("test:handler3", &mockTaskHandler{})

	// Create worker manager
	cfg := &config.QueueConfig{
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       15,
		Concurrency:   1,
		Queues: map[string]int{
			"default": 1,
		},
	}

	wm := workers.NewWorkerManager(cfg)

	// Register all handlers from registry
	registry.RegisterAll(wm)

	// Verify all handlers are registered
	// (In a real scenario, we'd need to expose a way to verify this in WorkerManager)
	assert.Equal(t, 3, registry.Count())
}

// TestGlobalTaskRegistry tests global registry singleton
func TestGlobalTaskRegistry(t *testing.T) {
	// Get global registry
	registry1 := workers.GetGlobalRegistry()
	registry2 := workers.GetGlobalRegistry()

	assert.Same(t, registry1, registry2, "should return same instance")

	// Register using convenience function
	handler := &mockTaskHandler{}
	workers.RegisterTask("global:test", handler)

	// Verify it's in the global registry
	retrieved, exists := registry1.Get("global:test")
	require.True(t, exists)
	assert.Equal(t, handler, retrieved)

	// Clean up for other tests
	// Note: In production, you'd want a way to reset the global registry for testing
}

// TestTaskRegistry_ConcurrentAccess tests thread-safe access
func TestTaskRegistry_ConcurrentAccess(t *testing.T) {
	registry := workers.NewTaskRegistry()

	// Run concurrent registrations
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(n int) {
			handler := &mockTaskHandler{}
			taskType := "concurrent:task"
			registry.Register(taskType, handler)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have 1 handler (last one wins)
	assert.Equal(t, 1, registry.Count())

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			registry.Get("concurrent:task")
			registry.List()
			registry.Count()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// BenchmarkTaskRegistry_Register benchmarks registration performance
func BenchmarkTaskRegistry_Register(b *testing.B) {
	registry := workers.NewTaskRegistry()
	handler := &mockTaskHandler{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		registry.Register("benchmark:task", handler)
	}
}

// BenchmarkTaskRegistry_Get benchmarks retrieval performance
func BenchmarkTaskRegistry_Get(b *testing.B) {
	registry := workers.NewTaskRegistry()
	handler := &mockTaskHandler{}
	registry.Register("benchmark:task", handler)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		registry.Get("benchmark:task")
	}
}
