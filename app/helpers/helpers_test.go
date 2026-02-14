package helpers_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHelpersSuite(t *testing.T) {
	RegisterFailHandler(Fail)

	// Try to load .env.test, but don't fail if it doesn't exist (CI/CD compatibility)
	_ = godotenv.Load("../../.env.test")

	gin.SetMode(gin.TestMode)

	RunSpecs(t, "Helpers Test Suite")
}
