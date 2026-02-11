package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer("gin-server")
		ctx, span := tracer.Start(c.Request.Context(), c.FullPath())
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.route", c.FullPath()),
		)

		c.Next()

		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)

		if c.Writer.Status() >= 400 {
			span.SetStatus(codes.Error, "HTTP error")
		}
	}
}
