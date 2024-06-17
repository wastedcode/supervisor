package supervisor

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/sdk/trace"
)

type Telemetry interface {
	StartChild(ctx context.Context, moduleName, childName string) (context.Context, Child)
	Shutdown(ctx context.Context) error
}

type Otel interface {
	GetSpanExporter() trace.SpanExporter
	GetTraceProvider() *trace.TracerProvider
	GetLogger() *slog.Logger
	Shutdown(ctx context.Context) error
}
