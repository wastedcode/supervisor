package supervisor

import (
	"context"
	"log/slog"
	"sync/atomic"

	"github.com/wastedcode/supervisor/errors"
)

var (
	defaultTelemetry = &atomic.Value{}
)

type DefaultTelemetry struct {
	config   Config
	otelImpl Otel
}

func NewDefaultTelemetry(ctx context.Context, applicationName string) (*DefaultTelemetry, error) {
	config := Config{
		OtelEnvironment: OtelEnvLocal,
		ApplicationName: applicationName,
	}

	otelImpl, err := getOtelImpl(config)
	if err != nil {
		return nil, err
	}

	defaultTelemetry := &DefaultTelemetry{
		config:   config,
		otelImpl: otelImpl,
	}

	return defaultTelemetry, nil
}

func (d *DefaultTelemetry) Shutdown(ctx context.Context) error {
	return d.otelImpl.Shutdown(ctx)
}

func (d *DefaultTelemetry) StartChild(ctx context.Context, moduleName, childName string) (context.Context, Child) {
	ctx, span := d.otelImpl.GetTraceProvider().Tracer(moduleName).Start(ctx, childName)

	return ctx, Child{
		span:   span,
		logger: d.otelImpl.GetLogger().With(slog.String("module", moduleName), slog.String("child", childName), slog.String("trace_id", span.SpanContext().TraceID().String()), slog.String("span_id", span.SpanContext().SpanID().String())),
	}
}

func GetDefaultTelemetry() Telemetry {
	return defaultTelemetry.Load().(Telemetry)
}

func SetDefaultTelemetry(t Telemetry) {
	defaultTelemetry.Store(t)
}

func getOtelImpl(config Config) (Otel, error) {
	switch config.OtelEnvironment {
	case OtelEnvLocal:
		return initLocalOtel(config)
	default:
		return nil, errors.WithInternalDetailsAndStack(ErrUnsupportedOtelEnv, "env: %s", config.OtelEnvironment)
	}
}
