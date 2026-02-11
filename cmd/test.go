package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

// TestCommand runs integration tests
// Usage: go run main.go test
var TestCommand = &cli.Command{
	Name:    "test",
	Aliases: []string{"t"},
	Usage:   "Run integration and unit tests (similar to 'php artisan test')",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Value:   "all",
			Usage:   "Test type: unit, integration, e2e, or all",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Value:   false,
			Usage:   "Enable verbose output",
		},
		&cli.BoolFlag{
			Name:    "coverage",
			Aliases: []string{"c"},
			Value:   false,
			Usage:   "Generate coverage report",
		},
		&cli.StringFlag{
			Name:  "filter",
			Value: "",
			Usage: "Run only tests matching pattern (e.g., 'Auth' or 'User')",
		},
		&cli.BoolFlag{
			Name:  "setup",
			Value: false,
			Usage: "Setup test database before running tests",
		},
	},
	Action: runTests,
}

func runTests(c *cli.Context) error {
	testType := c.String("type")
	verbose := c.Bool("verbose")
	coverage := c.Bool("coverage")
	filter := c.String("filter")
	setup := c.Bool("setup")

	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║   Golang Starter Kit - Test Runner        ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println()

	// Setup test database if requested
	if setup {
		fmt.Println("📦 Setting up test database...")
		if err := setupTestDatabase(); err != nil {
			return fmt.Errorf("failed to setup test database: %w", err)
		}
		fmt.Println("✅ Test database ready")
		fmt.Println()
	}

	// Determine test paths based on type
	testPaths := getTestPaths(testType)
	if len(testPaths) == 0 {
		return fmt.Errorf("invalid test type: %s (use: unit, integration, e2e, or all)", testType)
	}

	fmt.Printf("🧪 Running %s tests...\n", testType)
	fmt.Println()

	// Build go test command
	args := []string{"test"}
	args = append(args, testPaths...)

	if verbose {
		args = append(args, "-v")
	}

	if coverage {
		args = append(args, "-cover")
		args = append(args, "-coverprofile=coverage.out")
	}

	if filter != "" {
		args = append(args, "-run", filter)
	}

	// Add timeout
	args = append(args, "-timeout", "30s")

	// Run tests
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "APP_ENV=testing")

	fmt.Printf("Running: go %s\n", strings.Join(args, " "))
	fmt.Println(strings.Repeat("─", 50))

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Println()
			fmt.Println(strings.Repeat("─", 50))
			fmt.Printf("❌ Tests failed with exit code: %d\n", exitErr.ExitCode())
			return exitErr
		}
		return err
	}

	fmt.Println()
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ All tests passed!")

	// Show coverage report if requested
	if coverage {
		fmt.Println()
		fmt.Println("📊 Generating coverage report...")
		showCoverage()
	}

	return nil
}

// getTestPaths returns test paths based on test type
func getTestPaths(testType string) []string {
	switch testType {
	case "unit":
		return []string{"./app/services/..."}
	case "integration":
		return []string{"./tests/integration/..."}
	case "e2e":
		return []string{"./tests/e2e/..."}
	case "all":
		return []string{"./app/services/...", "./tests/integration/...", "./tests/e2e/..."}
	default:
		return []string{}
	}
}

// setupTestDatabase runs the test database setup script
func setupTestDatabase() error {
	scriptPath := filepath.Join("tests", "setup_test_db.sh")

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("setup script not found: %s", scriptPath)
	}

	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// showCoverage displays coverage report
func showCoverage() {
	cmd := exec.Command("go", "tool", "cover", "-func=coverage.out")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: failed to generate coverage report: %v\n", err)
	}

	fmt.Println()
	fmt.Println("💡 View HTML coverage report:")
	fmt.Println("   go tool cover -html=coverage.out")
}
