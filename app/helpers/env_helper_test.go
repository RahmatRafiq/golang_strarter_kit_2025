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

			Expect(helpers.GetEnv("APP_NAME", "default")).To(Equal("Supply Chain Retail"))
			Expect(helpers.GetEnv("JWT_EXPIRE_MINUTES", "default")).To(Equal("60"))
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
			Expect(helpers.GetEnvInt("JWT_EXPIRE_MINUTES", 10)).To(Equal(60))
		})
	})
})
