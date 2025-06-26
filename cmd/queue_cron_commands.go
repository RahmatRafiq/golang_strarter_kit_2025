package cmd

import (
	"golang_starter_kit_2025/app/services"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

// QueueWorkCommand defines the command to run queue workers
var QueueWorkCommand = &cli.Command{
	Name:  "queue:work",
	Usage: "Start processing queue jobs",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "workers",
			Value: 3,
			Usage: "Number of workers to start",
		},
		&cli.StringFlag{
			Name:  "queue",
			Value: "default",
			Usage: "Name of the queue to process",
		},
	},
	Action: func(c *cli.Context) error {
		workerCount := c.Int("workers")
		queueName := c.String("queue")

		log.Printf("🚀 Starting %d queue workers for queue: %s", workerCount, queueName)

		// Create and start queue service
		queueService := services.NewQueueService(workerCount)
		queueService.Start()

		// Handle graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Wait for shutdown signal
		<-sigChan
		log.Println("📡 Received shutdown signal")

		// Stop queue service
		queueService.Stop()

		log.Println("✅ Queue workers stopped")
		return nil
	},
}

// CronRunCommand defines the command to run cron jobs
var CronRunCommand = &cli.Command{
	Name:  "cron:run",
	Usage: "Start running scheduled cron jobs",
	Action: func(c *cli.Context) error {
		log.Println("🚀 Starting cron job scheduler")

		// Create and start cron service
		cronService := services.NewCronService()
		cronService.Start()

		// Handle graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Wait for shutdown signal
		<-sigChan
		log.Println("📡 Received shutdown signal")

		// Stop cron service
		cronService.Stop()

		log.Println("✅ Cron scheduler stopped")
		return nil
	},
}

// QueueStatsCommand shows queue statistics
var QueueStatsCommand = &cli.Command{
	Name:  "queue:stats",
	Usage: "Show queue statistics",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "queue",
			Value: "default",
			Usage: "Name of the queue to check",
		},
	},
	Action: func(c *cli.Context) error {
		queueName := c.String("queue")

		// Create queue service (without starting workers)
		queueService := services.NewQueueService(1)

		// Get stats
		stats, err := queueService.GetStats()
		if err != nil {
			log.Printf("❌ Error getting queue stats: %v", err)
			return err
		}

		log.Printf("📊 Queue Statistics for '%s':", queueName)
		log.Printf("   Size: %d jobs", stats.Size)
		log.Printf("   Workers: %d", stats.Workers)

		return nil
	},
}

// CronStatusCommand shows cron job status
var CronStatusCommand = &cli.Command{
	Name:  "cron:status",
	Usage: "Show status of all cron jobs",
	Action: func(c *cli.Context) error {
		// Create cron service (without starting)
		cronService := services.NewCronService()

		// Register jobs to get their info
		cronService.Start()
		defer cronService.Stop()

		// Get status
		statuses := cronService.GetStatus()

		log.Println("📊 Cron Job Status:")
		for name, status := range statuses {
			log.Printf("   %s:", name)
			log.Printf("     Schedule: %s", status.Schedule)
			log.Printf("     Next Run: %v", status.NextRun)
			if status.LastRun != nil {
				log.Printf("     Last Run: %v", *status.LastRun)
			} else {
				log.Printf("     Last Run: Never")
			}
			log.Printf("     Active: %v", status.IsActive)
			log.Println()
		}

		return nil
	},
}

// QueueTestCommand sends test jobs to the queue
var QueueTestCommand = &cli.Command{
	Name:  "queue:test",
	Usage: "Send test jobs to the queue",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "count",
			Value: 1,
			Usage: "Number of test jobs to send",
		},
	},
	Action: func(c *cli.Context) error {
		count := c.Int("count")

		// Create queue service
		queueService := services.NewQueueService(1)

		log.Printf("📤 Sending %d test jobs to queue", count)

		for i := 0; i < count; i++ {
			// Send test email job
			err := queueService.DispatchEmailJob(
				"test@example.com",
				"Test Email",
				"This is a test email from the queue system.",
			)
			if err != nil {
				log.Printf("❌ Error sending test job %d: %v", i+1, err)
				continue
			}
			log.Printf("✅ Test job %d sent successfully", i+1)
		}

		return nil
	},
}
