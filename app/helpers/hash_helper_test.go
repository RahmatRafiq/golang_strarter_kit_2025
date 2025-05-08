package helpers_test

import (
	"golang_starter_kit_2025/app/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HashPasswordArgon2", func() {
	Context("when password is hashed", func() {
		It("should return hashed password", func() {
			password := "password"
			hashedPassword, err := helpers.HashPasswordArgon2(password, helpers.DefaultParams)
			Expect(err).NotTo(HaveOccurred())
			Expect(helpers.ComparePasswordArgon2(password, hashedPassword)).To(BeTrue())
		})
	})

	Context("when password is hashed with different params", func() {
		It("should return hashed password", func() {
			password := "password"
			hashedPassword, err := helpers.HashPasswordArgon2(password, &helpers.Argon2Params{
				Memory:      64 * 1024, // 64 MB
				Iterations:  1,
				Parallelism: 1,
				SaltLength:  8,
				KeyLength:   32,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(helpers.ComparePasswordArgon2(password, hashedPassword)).To(BeTrue())
		})
	})

	Context("when password is hashed with invalid params", func() {
		It("should return error", func() {
			password := "password"
			hashedPassword, err := helpers.HashPasswordArgon2(password, &helpers.Argon2Params{
				Memory:      64 * 1024, // 64 MB
				Iterations:  1,
				Parallelism: 1,
				SaltLength:  8,
				KeyLength:   32,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(helpers.ComparePasswordArgon2(password, hashedPassword+"a")).To(BeFalse())
		})
	})
})

var _ = Describe("ComparePasswordArgon2", func() {
	Context("when password is compared", func() {
		It("should return true", func() {
			password := "password"
			hashedPassword, err := helpers.HashPasswordArgon2(password, helpers.DefaultParams)
			Expect(err).NotTo(HaveOccurred())
			Expect(helpers.ComparePasswordArgon2(password, hashedPassword)).To(BeTrue())
		})
	})

	Context("when password is compared with invalid hash", func() {
		It("should return false", func() {
			password := "password"
			hashedPassword, err := helpers.HashPasswordArgon2(password, helpers.DefaultParams)
			Expect(err).NotTo(HaveOccurred())
			Expect(helpers.ComparePasswordArgon2(password, hashedPassword+"a")).To(BeFalse())
		})
	})
})
