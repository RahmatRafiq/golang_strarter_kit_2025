package queue

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Worker represents a queue worker
type Worker struct {
	ID        int
	Queue     QueueInterface
	QueueName string
	stopChan  chan bool
	wg        *sync.WaitGroup
}

// WorkerPool manages multiple workers
type WorkerPool struct {
	workers     []*Worker
	queue       QueueInterface
	queueName   string
	workerCount int
	stopChan    chan bool
	wg          sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workerCount int, queue QueueInterface, queueName string) *WorkerPool {
	if queueName == "" {
		queueName = "default"
	}

	return &WorkerPool{
		workers:     make([]*Worker, 0, workerCount),
		queue:       queue,
		queueName:   queueName,
		workerCount: workerCount,
		stopChan:    make(chan bool),
	}
}

// Start begins processing jobs with multiple workers
func (wp *WorkerPool) Start() {
	log.Printf("🚀 Starting worker pool with %d workers for queue: %s", wp.workerCount, wp.queueName)

	for i := 0; i < wp.workerCount; i++ {
		worker := &Worker{
			ID:        i + 1,
			Queue:     wp.queue,
			QueueName: wp.queueName,
			stopChan:  make(chan bool),
			wg:        &wp.wg,
		}
		wp.workers = append(wp.workers, worker)

		wp.wg.Add(1)
		go worker.start()
	}

	// Handle graceful shutdown
	go wp.handleShutdown()
}

// Stop gracefully stops all workers
func (wp *WorkerPool) Stop() {
	log.Println("🛑 Stopping worker pool...")

	close(wp.stopChan)
	for _, worker := range wp.workers {
		close(worker.stopChan)
	}

	wp.wg.Wait()
	log.Println("✅ Worker pool stopped")
}

// Wait blocks until all workers are stopped
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// handleShutdown handles graceful shutdown signals
func (wp *WorkerPool) handleShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case <-c:
		log.Println("📡 Received shutdown signal")
		wp.Stop()
	case <-wp.stopChan:
		return
	}
}

// start begins the worker's job processing loop
func (w *Worker) start() {
	defer w.wg.Done()

	log.Printf("👷 Worker %d started", w.ID)

	for {
		select {
		case <-w.stopChan:
			log.Printf("👷 Worker %d stopped", w.ID)
			return
		default:
			w.processJob()
		}
	}
}

// processJob processes a single job from the queue
func (w *Worker) processJob() {
	job, err := w.Queue.Pop(w.QueueName)
	if err != nil {
		// No job available or error occurred, wait a bit
		time.Sleep(1 * time.Second)
		return
	}

	if job == nil {
		return
	}

	log.Printf("🔄 Worker %d processing job: %s (ID: %s)", w.ID, job.Name, job.ID)

	startTime := time.Now()

	// Create job instance from registry
	jobInstance, err := CreateJob(job.Name, job.Data)
	if err != nil {
		log.Printf("❌ Worker %d failed to create job %s: %v", w.ID, job.Name, err)
		w.Queue.Failed(job, err)
		return
	}

	// Execute the job
	err = jobInstance.Handle()
	duration := time.Since(startTime)

	if err != nil {
		log.Printf("❌ Worker %d job %s failed (took %v): %v", w.ID, job.Name, duration, err)
		w.Queue.Failed(job, err)
	} else {
		log.Printf("✅ Worker %d job %s completed successfully (took %v)", w.ID, job.Name, duration)
	}
}

// QueueStats provides statistics about the queue
type QueueStats struct {
	QueueName string `json:"queue_name"`
	Size      int64  `json:"size"`
	Workers   int    `json:"workers"`
}

// GetStats returns current queue statistics
func (wp *WorkerPool) GetStats() (*QueueStats, error) {
	size, err := wp.queue.Size(wp.queueName)
	if err != nil {
		return nil, err
	}

	return &QueueStats{
		QueueName: wp.queueName,
		Size:      size,
		Workers:   wp.workerCount,
	}, nil
}

// DelayedJobProcessor processes delayed jobs
type DelayedJobProcessor struct {
	queue    *RedisQueue
	stopChan chan bool
	interval time.Duration
}

// NewDelayedJobProcessor creates a new delayed job processor
func NewDelayedJobProcessor(queue *RedisQueue) *DelayedJobProcessor {
	return &DelayedJobProcessor{
		queue:    queue,
		stopChan: make(chan bool),
		interval: 30 * time.Second, // Check every 30 seconds
	}
}

// Start begins processing delayed jobs
func (djp *DelayedJobProcessor) Start() {
	log.Println("🕐 Starting delayed job processor")

	go func() {
		ticker := time.NewTicker(djp.interval)
		defer ticker.Stop()

		for {
			select {
			case <-djp.stopChan:
				log.Println("🕐 Delayed job processor stopped")
				return
			case <-ticker.C:
				if err := djp.queue.ProcessDelayedJobs(); err != nil {
					log.Printf("❌ Error processing delayed jobs: %v", err)
				}
			}
		}
	}()
}

// Stop stops the delayed job processor
func (djp *DelayedJobProcessor) Stop() {
	close(djp.stopChan)
}
