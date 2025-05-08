package casts_test

import (
	"time"

	"golang_starter_kit_2025/app/casts"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewJwtClaims", func() {
	It("should return jwt claims", func() {
		userID := uint(1)
		expiredAt := time.Now().Unix()

		claims := casts.NewJwtClaims(userID, expiredAt)

		Expect(claims["user_id"]).To(Equal(userID))
		Expect(claims["expired_at"]).To(Equal(expiredAt))
	})
})

var _ = Describe("ParseJwtClaims", func() {
	Context("when parse JWT claims", func() {
		It("should parse valid JWT claims", func() {
			userID := uint(123)
			expiredAt := time.Now().Add(time.Hour).Unix()

			claims := jwt.MapClaims{
				"user_id":    float64(userID),
				"expired_at": float64(expiredAt),
			}

			parsedClaims := casts.ParseJwtClaims(claims)

			Expect(parsedClaims.UserID).To(Equal(userID))
			Expect(parsedClaims.ExpiredAt).To(Equal(expiredAt))
		})

		It("should handle invalid JWT claims", func() {
			claims := jwt.MapClaims{
				"user_id":    "invalid_user_id",
				"expired_at": "invalid_expired_at",
			}

			defer func() {
				if r := recover(); r != nil {
					Expect(r).To(HaveOccurred())
				}
			}()

			casts.ParseJwtClaims(claims)
		})
	})
})
