package casts_test

import (
	"time"

	"golang_starter_kit_2025/app/casts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Token", func() {
	var (
		token     string
		expiredAt time.Time
	)

	BeforeEach(func() {
		token = "sample_token"
		expiredAt = time.Now().Add(24 * time.Hour)
	})

	Describe("Creating a new Token", func() {
		It("should create a token with the correct fields", func() {
			t := casts.Token{
				Token:     token,
				ExpiredAt: expiredAt,
			}

			Expect(t.Token).To(Equal(token))
			Expect(t.ExpiredAt).To(BeTemporally("~", expiredAt, time.Second))
		})
	})

	Describe("Token expiration", func() {
		It("should correctly identify if the token is expired", func() {
			t := casts.Token{
				Token:     token,
				ExpiredAt: expiredAt,
			}

			Expect(t.ExpiredAt.After(time.Now())).To(BeTrue())
		})

		It("should correctly identify if the token is not expired", func() {
			t := casts.Token{
				Token:     token,
				ExpiredAt: time.Now().Add(-1 * time.Hour),
			}

			Expect(t.ExpiredAt.Before(time.Now())).To(BeTrue())
		})
	})
})
