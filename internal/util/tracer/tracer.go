package tracer

import (
	"context"
	"fmt"
	whttp "github.com/SyaibanAhmadRamadhan/http-wrapper"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"runtime"
)

func Error(err error) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Errorf("%s:%d: %w", frame.Function, frame.Line, err)
}

func RecordErrorOtel(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

func GetTraceParent(ctx context.Context) *string {
	traceParent := whttp.GetTraceParent(ctx)
	if traceParent != "" {
		return &traceParent
	}

	return nil
}
