package main

import (
	jaegerExporter "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func NewTracer(url, name string) (*tracesdk.TracerProvider, error) {
	exp, err := jaegerExporter.New(jaegerExporter.WithCollectorEndpoint(jaegerExporter.WithEndpoint(url)))

	if err != nil {
		return nil, err
	}

	tracerObj := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
		)),
	)

	return tracerObj, nil
}
