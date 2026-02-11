package services

import (
	"context"
	"fmt"

	"golang_starter_kit_2025/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type TracingService struct {
	tracer   trace.Tracer
	provider *sdktrace.TracerProvider
}

func NewTracingService() (*TracingService, error) {
	cfg := config.GetTracingConfig()

	if !cfg.Enabled {
		return &TracingService{}, nil
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerURL)))
	if err != nil {
		return nil, fmt.Errorf("failed to create jaeger exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		)),
	)

	otel.SetTracerProvider(tp)

	return &TracingService{
		tracer:   tp.Tracer(cfg.ServiceName),
		provider: tp,
	}, nil
}

func (s *TracingService) Shutdown(ctx context.Context) error {
	if s.provider != nil {
		return s.provider.Shutdown(ctx)
	}
	return nil
}

func (s *TracingService) GetTracer() trace.Tracer {
	return s.tracer
}
