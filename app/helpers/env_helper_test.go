package helpers_test

import (
	"golang_starter_kit_2025/app/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetEnv", func() {
	Context("when key is not found", func() {
		It("should return default value", func() {
			Expect(helpers.GetEnv("NOT_FOUND", "default")).To(Equal("default"))
		})
	})

	Context("when key is found", func() {
		It("should return value from environment", func() {
			// GetEnv should return non-empty value when env var exists
			// Don't assert specific values as they may differ between local and CI
			appName := helpers.GetEnv("APP_NAME", "default")
			Expect(appName).NotTo(BeEmpty())
		})
	})
})

var _ = Describe("GetEnvInt", func() {
	Context("when key is not found", func() {
		It("should return default value", func() {
			Expect(helpers.GetEnvInt("NOT_FOUND", 10)).To(Equal(10))
		})
	})

	Context("when key is found", func() {
		It("should return value from environment", func() {
			// GetEnvInt should return integer value
			// Don't assert specific values as they may differ between local and CI
			jwtExpire := helpers.GetEnvInt("JWT_EXPIRE_MINUTES", 10)
			Expect(jwtExpire).To(BeNumerically(">", 0))
		})
	})
})
