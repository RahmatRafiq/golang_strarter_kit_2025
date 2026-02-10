package helpers

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the global logger with appropriate configuration
func InitLogger() {
	// Set log level from environment variable
	logLevel := GetEnv("LOG_LEVEL", "info")
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Configure time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Use console writer for development (pretty, colored output)
	appEnv := GetEnv("APP_ENV", "production")
	if appEnv == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	log.Info().
		Str("env", appEnv).
		Str("log_level", logLevel).
		Msg("Logger initialized")
}

// GetLogger returns a logger with context
func GetLogger() *zerolog.Logger {
	return &log.Logger
}
