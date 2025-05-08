package helpers_test

import (
	"net/http/httptest"
	"time"

	"golang_starter_kit_2025/app/casts"
	"golang_starter_kit_2025/app/helpers"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tidwall/gjson"
)

var _ = Describe("ResponseSuccess", func() {
	var (
		w   *httptest.ResponseRecorder
		ctx *gin.Context
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(w)
	})

	Context("when response is success", func() {
		It("should return success response", func() {
			response := helpers.ResponseParams[any]{
				Message: "This is a message",
			}

			helpers.ResponseSuccess(ctx, &response, 200)

			Expect(w.Code).To(Equal(200))
			Expect(w.Body.String()).To(Equal("{\"status\":\"success\",\"data\":[{\"key\":\"value\"}],\"item\":{\"key\":\"value\"},\"message\":\"This is a message\"}"))
		})
	})

	Context("when response is success with token", func() {
		It("should return success response with token", func() {
			response := helpers.ResponseParams[any]{
				Message: "This is a message",
				Token: &casts.Token{
					Token:     "access_token",
					ExpiredAt: time.Now().Add(time.Hour),
				},
			}

			helpers.ResponseSuccess(ctx, &response, 200)

			Expect(w.Code).To(Equal(200))

			Expect(w.Body.String()).To(MatchRegexp(`{"status":"success","data":\[{"key":"value"}\],"item":{"key":"value"},"message":"This is a message","token":{"token":"access_token","expired_at":".*"}}`))

			// get expired at from response
			expiredAt := gjson.Get(w.Body.String(), "token.expired_at").String()

			// convert expired at to time
			expiredAtTime, err := time.Parse(time.RFC3339, expiredAt)
			Expect(err).To(BeNil())

			// check if expired at is valid
			Expect(expiredAtTime).ToNot(BeNil())

			Expect(expiredAtTime).To(BeTemporally("~", time.Now().Add(time.Hour)))
		})
	})
})

var _ = Describe("ResponseError", func() {
	var (
		w   *httptest.ResponseRecorder
		ctx *gin.Context
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(w)
	})

	Context("when response is error", func() {
		It("should return error response", func() {
			response := helpers.ResponseParams[any]{
				Message:   "This is a message",
				Reference: "INV-123456",
			}

			helpers.ResponseError(ctx, &response, 400)

			Expect(w.Code).To(Equal(400))
			Expect(w.Body.String()).To(Equal("{\"status\":\"error\",\"data\":[{\"key\":\"value\"}],\"item\":{\"key\":\"value\"},\"message\":\"This is a message\",\"reference\":\"INV-123456\"}"))
		})
	})
})
