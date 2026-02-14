package controllers_test

import (
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestControllerssSuite(t *testing.T) {
	RegisterFailHandler(Fail)

	// Try to load .env.test, but don't fail if it doesn't exist (CI/CD compatibility)
	_ = godotenv.Load("../../.env.test")

	RunSpecs(t, "Controllers Test Suite")
}
