package services

import (
	"golang_starter_kit_2025/app/cron"
	"golang_starter_kit_2025/app/jobs"
	"golang_starter_kit_2025/app/queue"
	"golang_starter_kit_2025/config"
	"log"
)

// QueueService manages the queue system
type QueueService struct {
	Queue            queue.QueueInterface
	WorkerPool       *queue.WorkerPool
	DelayedProcessor *queue.DelayedJobProcessor
}

// CronService manages the cron system
type CronService struct {
	Manager *cron.CronManager
}

// NewQueueService creates a new queue service
func NewQueueService(workerCount int) *QueueService {
	// Initialize Redis queue
	redisQueue := queue.NewRedisQueue()

	// Create worker pool
	workerPool := queue.NewWorkerPool(workerCount, redisQueue, "default")

	// Create delayed job processor
	delayedProcessor := queue.NewDelayedJobProcessor(redisQueue)

	return &QueueService{
		Queue:            redisQueue,
		WorkerPool:       workerPool,
		DelayedProcessor: delayedProcessor,
	}
}

// NewCronService creates a new cron service
func NewCronService() *CronService {
	return &CronService{
		Manager: cron.NewCronManager(),
	}
}

// StartQueue starts the queue processing
func (qs *QueueService) Start() {
	// Register all job types
	qs.registerJobs()

	// Start worker pool
	qs.WorkerPool.Start()

	// Start delayed job processor
	qs.DelayedProcessor.Start()

	log.Println("🚀 Queue service started")
}

// StopQueue stops the queue processing
func (qs *QueueService) Stop() {
	qs.WorkerPool.Stop()
	qs.DelayedProcessor.Stop()
	log.Println("🛑 Queue service stopped")
}

// StartCron starts the cron job manager
func (cs *CronService) Start() {
	// Register all cron jobs
	cs.registerCronJobs()

	// Start cron manager
	cs.Manager.Start()

	log.Println("🚀 Cron service started")
}

// StopCron stops the cron job manager
func (cs *CronService) Stop() {
	cs.Manager.Stop()
	log.Println("🛑 Cron service stopped")
}

// registerJobs registers all available queue job types
func (qs *QueueService) registerJobs() {
	// Register email job
	queue.RegisterJob("email_job", jobs.NewEmailJob)

	// Register file cleanup job
	queue.RegisterJob("file_cleanup_job", jobs.NewFileCleanupJob)

	log.Println("📝 Queue jobs registered")
}

// registerCronJobs registers all available cron jobs
func (cs *CronService) registerCronJobs() {
	// Register database cleanup cron job
	dbCleanupJob := jobs.NewDatabaseCleanupCronJob()
	cs.Manager.AddJob(dbCleanupJob, cron.EveryDay)

	// Register backup cron job
	backupJob := jobs.NewBackupCronJob()
	cs.Manager.AddJob(backupJob, cron.EveryDay)

	// Register report generation cron job
	reportJob := jobs.NewReportGenerationCronJob()
	cs.Manager.AddJob(reportJob, cron.EveryWeek)

	log.Println("📝 Cron jobs registered")
}

// QueueStats returns queue statistics
func (qs *QueueService) GetStats() (*queue.QueueStats, error) {
	return qs.WorkerPool.GetStats()
}

// GetCronStatus returns cron job statuses
func (cs *CronService) GetStatus() map[string]cron.JobStatus {
	return cs.Manager.GetJobStatus()
}

// DispatchEmailJob dispatches an email job to the queue
func (qs *QueueService) DispatchEmailJob(to, subject, body string) error {
	job := jobs.NewEmailJob()
	job.SetData(map[string]interface{}{
		"to":      to,
		"subject": subject,
		"body":    body,
	})

	return qs.Queue.Push(job)
}

// DispatchFileCleanupJob dispatches a file cleanup job to the queue
func (qs *QueueService) DispatchFileCleanupJob(directory string, olderThanHours float64) error {
	job := jobs.NewFileCleanupJob()
	job.SetData(map[string]interface{}{
		"directory":        directory,
		"older_than_hours": olderThanHours,
	})

	return qs.Queue.Push(job)
}

// BootstrapServices initializes and starts all services
func BootstrapServices() (*QueueService, *CronService) {
	// Initialize Redis first
	config.InitRedis()

	// Create services
	queueService := NewQueueService(3) // 3 workers by default
	cronService := NewCronService()

	// Start services
	queueService.Start()
	cronService.Start()

	return queueService, cronService
}

// ShutdownServices gracefully shuts down all services
func ShutdownServices(queueService *QueueService, cronService *CronService) {
	if queueService != nil {
		queueService.Stop()
	}

	if cronService != nil {
		cronService.Stop()
	}

	// Close Redis connection
	config.CloseRedis()

	log.Println("✅ All services shut down gracefully")
}
