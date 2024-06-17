package supervisor

import (
	"context"
	"log/slog"

	"github.com/wastedcode/supervisor/errors"
	"go.opentelemetry.io/otel/sdk/trace"
)

type otelImpl struct {
	spanExporter  trace.SpanExporter
	traceProvider *trace.TracerProvider
	logger        *slog.Logger
}

func newOtelImpl() otelImpl {
	return otelImpl{}
}

func (o *otelImpl) GetSpanExporter() trace.SpanExporter {
	return o.spanExporter
}

func (o *otelImpl) GetTraceProvider() *trace.TracerProvider {
	return o.traceProvider
}

func (o *otelImpl) GetLogger() *slog.Logger {
	return o.logger
}

func (o *otelImpl) Shutdown(ctx context.Context) error {
	if err := o.traceProvider.ForceFlush(ctx); err != nil {
		return errors.WithInternalDetailsAndStack(err, "failed to flush trace provider")
	}

	if err := o.traceProvider.Shutdown(ctx); err != nil {
		return errors.WithInternalDetailsAndStack(err, "failed to shutdown trace provider")
	}

	if err := o.spanExporter.Shutdown(ctx); err != nil {
		return errors.WithInternalDetailsAndStack(err, "failed to shutdown span exporter")
	}

	return nil
}
