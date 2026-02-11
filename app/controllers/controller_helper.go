package controllers

import (
	"strconv"

	"golang_starter_kit_2025/app/helpers"

	"github.com/gin-gonic/gin"
)

// ParseIDParam parses the "id" URL parameter to uint and returns it.
// If parsing fails, it sends an error response and returns (0, false).
// Usage: if id, ok := ParseIDParam(ctx); ok { ... }
func ParseIDParam(ctx *gin.Context) (uint, bool) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || idUint == 0 {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Invalid ID format"},
			Message:   "Invalid ID",
			Reference: "ERROR-INVALID-ID",
		}, 400)
		return 0, false
	}
	return uint(idUint), true
}
