package observability

import (
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"user-service/internal/infra"
	"user-service/internal/primitive"
)

func InitOtel(endpointExp, nameTracer string) primitive.CloseFn {
	spanExp := infra.NewOTLP(endpointExp)
	traceProvide, closeFnTracer, err := newOtelProvider(nameTracer).start(spanExp)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed initializing the tracer provider")
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(traceProvide)

	return closeFnTracer
}
