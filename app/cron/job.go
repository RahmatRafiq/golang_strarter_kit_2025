package cron

import (
	"time"
)

// CronJobInterface defines the interface for cron jobs
type CronJobInterface interface {
	Execute() error
	GetName() string
	GetSchedule() string
}

// BaseCronJob provides a basic implementation
type BaseCronJob struct {
	Name     string
	Schedule string
}

func (c *BaseCronJob) GetName() string {
	return c.Name
}

func (c *BaseCronJob) GetSchedule() string {
	return c.Schedule
}

func (c *BaseCronJob) Execute() error {
	// Default implementation - should be overridden
	return nil
}

// CronSchedule represents different schedule types
type CronSchedule struct {
	Expression string
	Interval   time.Duration
	Type       ScheduleType
}

type ScheduleType string

const (
	ScheduleTypeExpression ScheduleType = "expression"
	ScheduleTypeInterval   ScheduleType = "interval"
	ScheduleTypeOnce       ScheduleType = "once"
)

// Predefined schedules
var (
	EveryMinute    = CronSchedule{Expression: "* * * * *", Type: ScheduleTypeExpression}
	Every5Minutes  = CronSchedule{Expression: "*/5 * * * *", Type: ScheduleTypeExpression}
	Every10Minutes = CronSchedule{Expression: "*/10 * * * *", Type: ScheduleTypeExpression}
	Every15Minutes = CronSchedule{Expression: "*/15 * * * *", Type: ScheduleTypeExpression}
	Every30Minutes = CronSchedule{Expression: "*/30 * * * *", Type: ScheduleTypeExpression}
	EveryHour      = CronSchedule{Expression: "0 * * * *", Type: ScheduleTypeExpression}
	EveryDay       = CronSchedule{Expression: "0 0 * * *", Type: ScheduleTypeExpression}
	EveryWeek      = CronSchedule{Expression: "0 0 * * 0", Type: ScheduleTypeExpression}
	EveryMonth     = CronSchedule{Expression: "0 0 1 * *", Type: ScheduleTypeExpression}
)

// Helper functions for interval-based schedules
func EverySeconds(seconds int) CronSchedule {
	return CronSchedule{
		Interval: time.Duration(seconds) * time.Second,
		Type:     ScheduleTypeInterval,
	}
}

func EveryMinutes(minutes int) CronSchedule {
	return CronSchedule{
		Interval: time.Duration(minutes) * time.Minute,
		Type:     ScheduleTypeInterval,
	}
}

func EveryHours(hours int) CronSchedule {
	return CronSchedule{
		Interval: time.Duration(hours) * time.Hour,
		Type:     ScheduleTypeInterval,
	}
}

// CronJobRegistry stores registered cron jobs
var CronJobRegistry = make(map[string]CronJobInterface)

// RegisterCronJob registers a cron job
func RegisterCronJob(job CronJobInterface) {
	CronJobRegistry[job.GetName()] = job
}

// GetRegisteredJobs returns all registered cron jobs
func GetRegisteredJobs() map[string]CronJobInterface {
	return CronJobRegistry
}
