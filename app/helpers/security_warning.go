package helpers

import (
	"bufio"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

func CheckSecurityWarning() bool {
	skipAuth := os.Getenv("SKIP_AUTH")
	appEnv := os.Getenv("APP_ENV")

	if skipAuth != "true" {
		return true
	}

	log.Warn().
		Str("skip_auth", skipAuth).
		Str("environment", appEnv).
		Msg("SECURITY WARNING: Authentication is DISABLED")

	log.Warn().
		Msg("All API endpoints are accessible WITHOUT authentication")
	log.Warn().
		Msg("No JWT token required")
	log.Warn().
		Msg("All requests will use default user_id=1")
	log.Warn().
		Msg("This is ONLY for development/testing purposes")

	if appEnv == "production" || appEnv == "prod" {
		log.Error().
			Str("environment", appEnv).
			Str("skip_auth", skipAuth).
			Msg("CRITICAL: Running in PRODUCTION mode with SKIP_AUTH=true - SEVERE SECURITY RISK")

		log.Error().
			Msg("Please set SKIP_AUTH=false in your .env file immediately")

		log.Warn().
			Msg("Do you want to continue anyway? (type 'yes' to continue, anything else to exit)")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error reading user input")
			return false
		}
		input = strings.TrimSpace(strings.ToLower(input))

		if input != "yes" {
			log.Info().
				Msg("Application startup canceled for security reasons")
			log.Info().
				Msg("Please update your .env file: SKIP_AUTH=false")
			return false
		}

		log.Warn().
			Msg("Proceeding with authentication DISABLED in production - You have been warned!")
		return true
	}

	// Development mode - auto-accept with warning
	log.Info().
		Str("environment", appEnv).
		Msg("Development mode detected - SKIP_AUTH is acceptable for local development")
	log.Warn().
		Msg("Auto-continuing in development mode - Remember: This is for development only!")

	return true
}

func PrintServerStartupInfo() {
	skipAuth := os.Getenv("SKIP_AUTH")
	appPort := GetEnv("APP_PORT", "8080")
	appHost := GetEnv("APP_HOST", "localhost")
	appScheme := GetEnv("APP_SCHEME", "http")

	log.Info().
		Str("url", appScheme+"://"+appHost+":"+appPort).
		Str("swagger", appScheme+"://"+appHost+":"+appPort+"/swagger/index.html").
		Bool("auth_enabled", skipAuth != "true").
		Msg("Server started successfully")

	if skipAuth == "true" {
		log.Warn().
			Str("auth_mode", "DISABLED").
			Str("skip_auth", "true").
			Msg("Authentication is DISABLED")
	} else {
		log.Info().
			Str("auth_mode", "ENABLED").
			Msg("Authentication is ENABLED (JWT Required)")
	}
}
