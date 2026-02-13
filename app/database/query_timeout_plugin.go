package database

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// QueryTimeoutPlugin adds timeout context to all database queries
// MED-5: Prevent long-running queries from blocking resources
type QueryTimeoutPlugin struct {
	Timeout time.Duration
}

// Name returns the plugin name
func (p *QueryTimeoutPlugin) Name() string {
	return "query_timeout_plugin"
}

// Initialize initializes the plugin
func (p *QueryTimeoutPlugin) Initialize(db *gorm.DB) error {
	// Register callback before query execution
	db.Callback().Query().Before("gorm:query").Register("query_timeout:before", p.beforeQuery)
	db.Callback().Create().Before("gorm:create").Register("query_timeout:before", p.beforeQuery)
	db.Callback().Update().Before("gorm:update").Register("query_timeout:before", p.beforeQuery)
	db.Callback().Delete().Before("gorm:delete").Register("query_timeout:before", p.beforeQuery)
	db.Callback().Row().Before("gorm:row").Register("query_timeout:before", p.beforeQuery)
	db.Callback().Raw().Before("gorm:raw").Register("query_timeout:before", p.beforeQuery)

	log.Info().
		Dur("timeout", p.Timeout).
		Msg("Query timeout plugin initialized")

	return nil
}

// beforeQuery adds timeout context before query execution
func (p *QueryTimeoutPlugin) beforeQuery(db *gorm.DB) {
	// Skip if query already has a deadline
	if _, hasDeadline := db.Statement.Context.Deadline(); hasDeadline {
		return
	}

	// Add timeout context
	ctx, cancel := context.WithTimeout(db.Statement.Context, p.Timeout)
	db.Statement.Context = ctx

	// Store cancel function for cleanup (optional)
	// In production, consider using statement.AddVar to track cancel functions
	_ = cancel
}

// NewQueryTimeoutPlugin creates a new query timeout plugin
// Default timeout: 30 seconds
func NewQueryTimeoutPlugin(timeout time.Duration) *QueryTimeoutPlugin {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &QueryTimeoutPlugin{
		Timeout: timeout,
	}
}
