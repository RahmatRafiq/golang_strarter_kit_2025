package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%s%s⚠️  SECURITY WARNING ⚠️%s\n", ColorBold, ColorRed, ColorReset)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	fmt.Printf("%s%sAuthentication is currently DISABLED!%s\n", ColorBold, ColorYellow, ColorReset)
	fmt.Printf("%sEnvironment: %s%s%s\n", ColorWhite, ColorCyan, appEnv, ColorReset)
	fmt.Printf("%sSKIP_AUTH=%s%strue%s\n\n", ColorWhite, ColorBold, ColorRed, ColorReset)

	fmt.Printf("%sThis means:%s\n", ColorBold, ColorReset)
	fmt.Printf("%s  • All API endpoints are accessible WITHOUT authentication%s\n", ColorYellow, ColorReset)
	fmt.Printf("%s  • No JWT token required%s\n", ColorYellow, ColorReset)
	fmt.Printf("%s  • All requests will use default user_id=1%s\n", ColorYellow, ColorReset)
	fmt.Printf("%s  • This is ONLY for development/testing purposes%s\n\n", ColorYellow, ColorReset)

	if appEnv == "production" || appEnv == "prod" {
		fmt.Printf("%s%s🚨 CRITICAL: You are running in PRODUCTION mode with SKIP_AUTH=true!%s\n", ColorBold, ColorRed, ColorReset)
		fmt.Printf("%s%sThis is a SEVERE SECURITY RISK!%s\n\n", ColorBold, ColorRed, ColorReset)
		fmt.Printf("%sPlease set SKIP_AUTH=false in your .env file immediately!%s\n\n", ColorRed, ColorReset)

		fmt.Printf("%sDo you want to continue anyway? (type 'yes' to continue, anything else to exit): %s", ColorRed, ColorReset)

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input != "yes" {
			fmt.Printf("\n%s%s❌ Application startup cancelled for security reasons.%s\n", ColorBold, ColorRed, ColorReset)
			fmt.Printf("%sPlease update your .env file:%s\n", ColorYellow, ColorReset)
			fmt.Printf("%s  SKIP_AUTH=false%s\n\n", ColorGreen, ColorReset)
			return false
		}

		fmt.Printf("\n%s%s⚠️  Proceeding with authentication DISABLED in production...%s\n", ColorBold, ColorRed, ColorReset)
		fmt.Printf("%s%sYou have been warned! This is extremely dangerous.%s\n\n", ColorBold, ColorYellow, ColorReset)
		return true
	}

	// Development mode - auto-accept with warning
	fmt.Printf("%s⚙️  Development Mode Detected%s\n", ColorGreen, ColorReset)
	fmt.Printf("%sThis setting is acceptable for local development.%s\n", ColorWhite, ColorReset)
	fmt.Printf("%s%s✅ Auto-continuing in development mode...%s\n\n", ColorBold, ColorGreen, ColorReset)
	fmt.Printf("%sRemember: This is for development only!%s\n\n", ColorYellow, ColorReset)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	return true
}

func PrintServerStartupInfo() {
	skipAuth := os.Getenv("SKIP_AUTH")
	appPort := GetEnv("APP_PORT", "8080")
	appHost := GetEnv("APP_HOST", "localhost")
	appScheme := GetEnv("APP_SCHEME", "http")

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%s%s🚀 Server Started Successfully%s\n", ColorBold, ColorGreen, ColorReset)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%sServer URL:%s %s%s://%s:%s%s\n", ColorBold, ColorReset, ColorCyan, appScheme, appHost, appPort, ColorReset)
	fmt.Printf("%sSwagger UI:%s %s%s://%s:%s/swagger/index.html%s\n", ColorBold, ColorReset, ColorCyan, appScheme, appHost, appPort, ColorReset)

	if skipAuth == "true" {
		fmt.Printf("%sAuth Mode:%s  %s%s⚠️  DISABLED (SKIP_AUTH=true)%s\n", ColorBold, ColorReset, ColorBold, ColorYellow, ColorReset)
	} else {
		fmt.Printf("%sAuth Mode:%s  %s%s✅ ENABLED (JWT Required)%s\n", ColorBold, ColorReset, ColorBold, ColorGreen, ColorReset)
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
}
