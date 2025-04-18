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

	// load .env.test file for all specs in this package
	err := godotenv.Load("../../.env.test")
	if err != nil {
		Fail("Error loading .env.test file")
	}

	gin.SetMode(gin.TestMode)

	RunSpecs(t, "Helpers Test Suite")
}
