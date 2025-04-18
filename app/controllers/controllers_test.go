package controllers_test

import (
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestControllerssSuite(t *testing.T) {
	RegisterFailHandler(Fail)

	// load .env.test
	err := godotenv.Load("../../.env.test")
	if err != nil {
		Fail("Error loading .env.test file")
	}

	RunSpecs(t, "Controllers Test Suite")
}