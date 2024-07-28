package observability

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"time"
	"user-service/internal/primitive"
)

type otelProvider struct {
	name     string
	exporter trace.SpanExporter
}

func newOtelProvider(name string) *otelProvider {
	return &otelProvider{
		name:     name,
		exporter: nil,
	}
}

func (t *otelProvider) start(exp trace.SpanExporter) (*trace.TracerProvider, primitive.CloseFn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(t.name),
		),
	)

	if err != nil {
		err = fmt.Errorf("failed to created resource: %w", err)
		return nil, nil, err
	}

	t.exporter = exp
	bsp := trace.NewBatchSpanProcessor(t.exporter)

	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)

	closeFn := func(ctx context.Context) (err error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err = t.exporter.Shutdown(ctx); err != nil {
			return err
		}

		if err = provider.Shutdown(ctx); err != nil {
			return err
		}

		return
	}
	return provider, closeFn, nil
}
