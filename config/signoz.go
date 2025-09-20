package config

import (
	"app/lib/logger"
	"app/lib/signoz"
	"context"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

func (c *Config) NewSignoz() (*sdktrace.TracerProvider, error) {
	tp, err := signoz.NewSignozTracer(signoz.SignozTracerOption{
		CollectorURL:     c.SIGNOZ_URL,
		ServiceName:      c.SIGNOZ_SERVICE_NAME,
		ServiceNamespace: c.SIGNOZ_SERVICE_NAMESPACE,
		Environment:      c.ENV,
		TraceSampleRate:  c.SIGNOZ_TRACE_SAMPLE_RATE,
	})
	if err != nil {
		logger.LogError(context.Background(), "signoz not connected", []zap.Field{
			zap.Error(err),
		}...)
		return nil, err
	}

	logger.LogInfo(context.Background(), "success connect to signoz", []zap.Field{}...)
	return tp, nil
}
