package signoz

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func NewSignozTracer(sto SignozTracerOption) (*sdktrace.TracerProvider, error) {
	sto.setDefaultValue()

	var secureOption otlptracehttp.Option
	if sto.InsecureMode {
		secureOption = otlptracehttp.WithInsecure()
	} else {
		secureOption = otlptracehttp.WithTLSClientConfig(nil)
	}

	headers := map[string]string{
		"signoz-access-token": sto.SignozToken,
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(sto.CollectorURL),
			otlptracehttp.WithHeaders(headers),
			secureOption,
		),
	)
	if err != nil {
		return nil, err
	}

	// Set up the resource attributes
	resourceAttributes := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(sto.ServiceName),
		semconv.DeploymentEnvironmentKey.String(sto.Environment),
		semconv.ServiceNamespaceKey.String(sto.ServiceNamespace),
		semconv.TelemetrySDKLanguageGo,
		semconv.TelemetrySDKVersionKey.String(sto.ServiceVersion),
	)

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	traceProvider := sdktrace.NewTracerProvider(
		// sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sto.TraceSampleRate))),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
		sdktrace.WithResource(resourceAttributes),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return traceProvider, nil
}

func StartSpan(ctx context.Context, op string, attrs ...attribute.KeyValue) (context.Context, *Span) {
	ctx, span := otel.Tracer("http.server").Start(ctx, op)
	span.AddEvent(op, trace.WithAttributes(attrs...))

	return ctx, &Span{
		Stack:      "signoz",
		SignozSpan: &span,
	}
}
