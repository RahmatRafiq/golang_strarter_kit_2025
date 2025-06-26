package jobs

import (
	"golang_starter_kit_2025/app/queue"
	"log"
	"time"
)

// FileCleanupJob handles cleaning up temporary files
type FileCleanupJob struct {
	queue.BaseJob
}

// NewFileCleanupJob creates a new file cleanup job
func NewFileCleanupJob() queue.JobInterface {
	return &FileCleanupJob{
		BaseJob: queue.BaseJob{
			Name: "file_cleanup_job",
		},
	}
}

// Handle processes the file cleanup job
func (f *FileCleanupJob) Handle() error {
	data := f.GetData()

	directory, ok := data["directory"].(string)
	if !ok {
		directory = "./tmp" // Default directory
	}

	olderThanHours, ok := data["older_than_hours"].(float64)
	if !ok {
		olderThanHours = 24 // Default to 24 hours
	}

	log.Printf("🧹 Starting file cleanup in directory: %s", directory)
	log.Printf("🧹 Cleaning files older than %.0f hours", olderThanHours)

	// Simulate file cleanup logic
	// In a real implementation, you would:
	// 1. Scan the directory
	// 2. Check file modification times
	// 3. Delete files older than the specified time

	cleanupTime := time.Duration(olderThanHours) * time.Hour
	log.Printf("🧹 Files older than %v would be cleaned up", cleanupTime)

	// Simulate some processing time
	time.Sleep(100 * time.Millisecond)

	log.Printf("✅ File cleanup completed for directory: %s", directory)
	return nil
}

// QueueFileCleanupJob is a helper function to queue a file cleanup job
func QueueFileCleanupJob(directory string, olderThanHours float64) error {
	job := NewFileCleanupJob()
	job.SetData(map[string]interface{}{
		"directory":        directory,
		"older_than_hours": olderThanHours,
	})

	// You would get the queue instance from your app context
	// For now, this is just a placeholder
	return nil
}
