package signoz

import "go.opentelemetry.io/otel/trace"

type SignozTracerOption struct {
	CollectorURL     string
	InsecureMode     bool
	ServiceName      string
	ServiceNamespace string
	Environment      string
	SignozToken      string
	ServiceVersion   string
	TraceSampleRate  float64
}

func (sto *SignozTracerOption) setDefaultValue() {
	if sto.ServiceVersion == "" {
		sto.ServiceVersion = "0.1.0"
	}

	if sto.Environment == "" {
		sto.Environment = "development"
	}

	if sto.TraceSampleRate == 0 {
		sto.TraceSampleRate = 0.2
	}

	if sto.Environment == "development" {
		sto.InsecureMode = true
	}
}

type Span struct {
	Stack      string
	SignozSpan *trace.Span
}

func (s *Span) Finish() {
	var span trace.Span
	if s.SignozSpan != nil {
		span = *s.SignozSpan
		span.End()
	}
}

func (s *Span) TraceID() string {
	var traceID string
	var span trace.Span

	if s.SignozSpan != nil {
		span = *s.SignozSpan
		traceID = span.SpanContext().TraceID().String()
	}

	return traceID
}
