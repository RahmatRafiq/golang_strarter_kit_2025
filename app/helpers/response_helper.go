package helpers

import (
	"golang_starter_kit_2025/app/casts"

	"github.com/gin-gonic/gin"
)

// Data
type ResponseParams[T any] struct {
	Status    *string           `json:"status"`
	Total     *int64            `json:"total,omitempty"`
	Data      *[]T              `json:"data,omitempty"`
	Errors    map[string]string `json:"errors,omitempty"`
	Item      *T                `json:"item,omitempty"`
	Message   string            `json:"message,omitempty"`
	Token     *casts.Token      `json:"token,omitempty"`
	Reference string            `json:"reference,omitempty"`
}

// ResponseSuccess Response
func ResponseSuccess[T any](ctx *gin.Context, params *ResponseParams[T], code int) {
	ctx.JSON(code, ResponseParams[T]{
		Status:  func() *string { s := "success"; return &s }(),
		Total:   params.Total,
		Data:    params.Data,
		Item:    params.Item,
		Message: params.Message,
		Token:   params.Token,
	})
}

// ResponseError Response
func ResponseError(ctx *gin.Context, params *ResponseParams[any], code int) {
	ctx.JSON(code, ResponseParams[any]{
		Status:    func() *string { s := "error"; return &s }(),
		Data:      params.Data,
		Item:      params.Item,
		Errors:    params.Errors,
		Message:   params.Message,
		Reference: params.Reference,
	})
}
