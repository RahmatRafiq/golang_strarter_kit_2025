package cron

import (
	"log"
	"sync"
	"time"
)

// CronManager manages and executes cron jobs
type CronManager struct {
	jobs     map[string]*ScheduledJob
	running  bool
	stopChan chan bool
	wg       sync.WaitGroup
	mutex    sync.RWMutex
}

// ScheduledJob wraps a cron job with its schedule information
type ScheduledJob struct {
	Job      CronJobInterface
	Schedule CronSchedule
	NextRun  time.Time
	LastRun  *time.Time
	IsActive bool
	stopChan chan bool
}

// NewCronManager creates a new cron manager
func NewCronManager() *CronManager {
	return &CronManager{
		jobs:     make(map[string]*ScheduledJob),
		stopChan: make(chan bool),
	}
}

// AddJob adds a job to the cron manager
func (cm *CronManager) AddJob(job CronJobInterface, schedule CronSchedule) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	scheduledJob := &ScheduledJob{
		Job:      job,
		Schedule: schedule,
		IsActive: true,
		stopChan: make(chan bool),
	}

	// Calculate next run time
	scheduledJob.NextRun = cm.calculateNextRun(schedule)

	cm.jobs[job.GetName()] = scheduledJob

	log.Printf("📅 Cron job '%s' scheduled for %v", job.GetName(), scheduledJob.NextRun)

	// If manager is already running, start this job
	if cm.running {
		cm.wg.Add(1)
		go cm.runJob(scheduledJob)
	}

	return nil
}

// RemoveJob removes a job from the cron manager
func (cm *CronManager) RemoveJob(jobName string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if job, exists := cm.jobs[jobName]; exists {
		job.IsActive = false
		close(job.stopChan)
		delete(cm.jobs, jobName)
		log.Printf("🗑️ Cron job '%s' removed", jobName)
	}
}

// Start begins executing all scheduled jobs
func (cm *CronManager) Start() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.running {
		return
	}

	cm.running = true
	log.Println("🚀 Starting cron manager")

	// Start all jobs
	for _, job := range cm.jobs {
		cm.wg.Add(1)
		go cm.runJob(job)
	}
}

// Stop stops all cron jobs
func (cm *CronManager) Stop() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if !cm.running {
		return
	}

	log.Println("🛑 Stopping cron manager")
	cm.running = false
	close(cm.stopChan)

	// Stop all jobs
	for _, job := range cm.jobs {
		job.IsActive = false
		close(job.stopChan)
	}

	cm.wg.Wait()
	log.Println("✅ Cron manager stopped")
}

// runJob runs a single scheduled job
func (cm *CronManager) runJob(scheduledJob *ScheduledJob) {
	defer cm.wg.Done()

	for {
		select {
		case <-scheduledJob.stopChan:
			return
		case <-cm.stopChan:
			return
		default:
			if !scheduledJob.IsActive {
				return
			}

			// Check if it's time to run
			if time.Now().After(scheduledJob.NextRun) {
				cm.executeJob(scheduledJob)

				// Calculate next run time
				scheduledJob.NextRun = cm.calculateNextRun(scheduledJob.Schedule)
			}

			// Sleep for a short duration to avoid busy waiting
			time.Sleep(1 * time.Second)
		}
	}
}

// executeJob executes a cron job
func (cm *CronManager) executeJob(scheduledJob *ScheduledJob) {
	jobName := scheduledJob.Job.GetName()
	log.Printf("🔄 Executing cron job: %s", jobName)

	startTime := time.Now()
	err := scheduledJob.Job.Execute()
	duration := time.Since(startTime)

	now := time.Now()
	scheduledJob.LastRun = &now

	if err != nil {
		log.Printf("❌ Cron job '%s' failed (took %v): %v", jobName, duration, err)
	} else {
		log.Printf("✅ Cron job '%s' completed successfully (took %v)", jobName, duration)
	}
}

// calculateNextRun calculates the next run time based on schedule
func (cm *CronManager) calculateNextRun(schedule CronSchedule) time.Time {
	now := time.Now()

	switch schedule.Type {
	case ScheduleTypeInterval:
		return now.Add(schedule.Interval)
	case ScheduleTypeExpression:
		// Simple cron expression parsing (basic implementation)
		return cm.parseBasicCronExpression(schedule.Expression, now)
	case ScheduleTypeOnce:
		return now // Run immediately once
	default:
		return now.Add(1 * time.Hour) // Default to 1 hour
	}
}

// parseBasicCronExpression provides basic cron expression parsing
func (cm *CronManager) parseBasicCronExpression(expression string, from time.Time) time.Time {
	// This is a simplified implementation
	// For production, consider using a proper cron library like "github.com/robfig/cron/v3"

	switch expression {
	case "* * * * *": // Every minute
		return from.Add(1 * time.Minute).Truncate(time.Minute)
	case "*/5 * * * *": // Every 5 minutes
		return from.Add(5 * time.Minute).Truncate(5 * time.Minute)
	case "*/10 * * * *": // Every 10 minutes
		return from.Add(10 * time.Minute).Truncate(10 * time.Minute)
	case "*/15 * * * *": // Every 15 minutes
		return from.Add(15 * time.Minute).Truncate(15 * time.Minute)
	case "*/30 * * * *": // Every 30 minutes
		return from.Add(30 * time.Minute).Truncate(30 * time.Minute)
	case "0 * * * *": // Every hour
		return from.Add(1 * time.Hour).Truncate(time.Hour)
	case "0 0 * * *": // Every day at midnight
		return from.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	case "0 0 * * 0": // Every week (Sunday)
		return from.AddDate(0, 0, 7).Truncate(24 * time.Hour)
	case "0 0 1 * *": // Every month (1st day)
		return from.AddDate(0, 1, 0).Truncate(24 * time.Hour)
	default:
		// Default to 1 hour for unknown expressions
		return from.Add(1 * time.Hour)
	}
}

// GetJobStatus returns the status of all jobs
func (cm *CronManager) GetJobStatus() map[string]JobStatus {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	status := make(map[string]JobStatus)
	for name, job := range cm.jobs {
		status[name] = JobStatus{
			Name:     name,
			Schedule: job.Schedule.Expression,
			NextRun:  job.NextRun,
			LastRun:  job.LastRun,
			IsActive: job.IsActive,
		}
	}
	return status
}

// JobStatus represents the status of a cron job
type JobStatus struct {
	Name     string     `json:"name"`
	Schedule string     `json:"schedule"`
	NextRun  time.Time  `json:"next_run"`
	LastRun  *time.Time `json:"last_run,omitempty"`
	IsActive bool       `json:"is_active"`
}
