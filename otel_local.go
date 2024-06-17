package supervisor

import (
	"os"

	"github.com/wastedcode/supervisor/errors"
	otelbase "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

func initLocalOtel(config Config) (Otel, error) {
	otel := newOtelImpl()

	spanExporter, err := localSpanExporter()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.ApplicationName),
		),
	)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(spanExporter),
	)

	otel.spanExporter = spanExporter
	otel.traceProvider = traceProvider
	otel.logger = getSimpleJsonLogger()

	// In GCP gcppropagator.CloudTraceOneWayPropagator{}
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otelbase.SetTextMapPropagator(propagator)

	return &otel, nil
}

func localSpanExporter() (trace.SpanExporter, error) {
	file, err := os.CreateTemp("", "otel-local-*.log")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return stdouttrace.New(stdouttrace.WithWriter(file))
}
