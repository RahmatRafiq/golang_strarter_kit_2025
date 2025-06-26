package controllers

import (
	"golang_starter_kit_2025/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// QueueController handles queue-related API endpoints
type QueueController struct {
	queueService *services.QueueService
}

// CronController handles cron-related API endpoints
type CronController struct {
	cronService *services.CronService
}

// NewQueueController creates a new queue controller
func NewQueueController() *QueueController {
	queueService := services.NewQueueService(3)
	return &QueueController{
		queueService: queueService,
	}
}

// NewCronController creates a new cron controller
func NewCronController() *CronController {
	cronService := services.NewCronService()
	return &CronController{
		cronService: cronService,
	}
}

// GetQueueStats godoc
// @Summary Get queue statistics
// @Description Get statistics about the queue system
// @Tags Queue
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/queue/stats [get]
func (qc *QueueController) GetQueueStats(c *gin.Context) {
	stats, err := qc.queueService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get queue stats",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Queue statistics retrieved successfully",
		"data":    stats,
	})
}

// DispatchEmailJob godoc
// @Summary Dispatch email job
// @Description Add an email job to the queue
// @Tags Queue
// @Accept json
// @Produce json
// @Param request body object{to=string,subject=string,body=string} true "Email job data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/queue/email [post]
func (qc *QueueController) DispatchEmailJob(c *gin.Context) {
	var request struct {
		To      string `json:"to" binding:"required,email"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	err := qc.queueService.DispatchEmailJob(request.To, request.Subject, request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to dispatch email job",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Email job dispatched successfully",
		"data": gin.H{
			"to":      request.To,
			"subject": request.Subject,
		},
	})
}

// DispatchFileCleanupJob godoc
// @Summary Dispatch file cleanup job
// @Description Add a file cleanup job to the queue
// @Tags Queue
// @Accept json
// @Produce json
// @Param request body object{directory=string,older_than_hours=number} true "File cleanup job data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/queue/file-cleanup [post]
func (qc *QueueController) DispatchFileCleanupJob(c *gin.Context) {
	var request struct {
		Directory      string  `json:"directory" binding:"required"`
		OlderThanHours float64 `json:"older_than_hours" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	err := qc.queueService.DispatchFileCleanupJob(request.Directory, request.OlderThanHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to dispatch file cleanup job",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "File cleanup job dispatched successfully",
		"data": gin.H{
			"directory":        request.Directory,
			"older_than_hours": request.OlderThanHours,
		},
	})
}

// SendTestJobs godoc
// @Summary Send test jobs to queue
// @Description Send multiple test jobs to the queue for testing purposes
// @Tags Queue
// @Accept json
// @Produce json
// @Param count query int false "Number of test jobs to send" default(1)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/queue/test [post]
func (qc *QueueController) SendTestJobs(c *gin.Context) {
	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		count = 1
	}

	successCount := 0
	errors := []string{}

	for i := 0; i < count; i++ {
		err := qc.queueService.DispatchEmailJob(
			"test@example.com",
			"Test Email from API",
			"This is a test email sent from the API endpoint.",
		)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			successCount++
		}
	}

	response := gin.H{
		"requested":     count,
		"successful":    successCount,
		"failed":        len(errors),
		"error_details": errors,
	}

	if len(errors) > 0 {
		c.JSON(http.StatusPartialContent, response)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "All test jobs dispatched successfully",
			"data":    response,
		})
	}
}

// GetCronStatus godoc
// @Summary Get cron job status
// @Description Get status of all registered cron jobs
// @Tags Cron
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/cron/status [get]
func (cc *CronController) GetCronStatus(c *gin.Context) {
	status := cc.cronService.GetStatus()
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Cron status retrieved successfully",
		"data":    status,
	})
}

// Health check endpoints

// QueueHealthCheck godoc
// @Summary Queue health check
// @Description Check if the queue system is healthy
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /api/health/queue [get]
func (qc *QueueController) HealthCheck(c *gin.Context) {
	stats, err := qc.queueService.GetStats()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Queue system is unhealthy",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Queue system is healthy",
		"data": gin.H{
			"status": "healthy",
			"stats":  stats,
		},
	})
}

// CronHealthCheck godoc
// @Summary Cron health check
// @Description Check if the cron system is healthy
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/health/cron [get]
func (cc *CronController) HealthCheck(c *gin.Context) {
	status := cc.cronService.GetStatus()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Cron system is healthy",
		"data": gin.H{
			"status": "healthy",
			"jobs":   len(status),
		},
	})
}
