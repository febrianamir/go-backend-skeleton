package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const CtxRequestID string = "X-Request-ID"
const CtxProcessID string = "X-Process-ID"

var instance *zap.Logger
var env = ""
var logPath = ""
var l = &lumberjack.Logger{}

func Init(setup LoggerSetup) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "@timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	instance, _ = config.Build(zap.AddCallerSkip(1))

	setup.valueDefault()
	env = setup.Env
	logPath = setup.Path

	if env == LOGGER_ENV_SETUP_DEVELOPMENT_VALUE {
		timeString := time.Now().Format("2006-01-02")
		l = &lumberjack.Logger{
			Filename:   fmt.Sprintf("%sgo-%s.log", logPath, timeString),
			MaxSize:    50, // Megabytes
			MaxBackups: 3,
			MaxAge:     14,   // Days
			Compress:   true, // Disabled by default
		}
		outputWriter := io.MultiWriter(os.Stdout, l)
		log.SetOutput(outputWriter)
	}
}

func CommonLog(ctx context.Context, level, message string, fields ...zap.Field) {
	if reqID, ok := ctx.Value(CtxRequestID).(string); ok {
		fields = append(fields, zap.String("request_id", reqID))
	}
	if reqID, ok := ctx.Value(CtxProcessID).(string); ok {
		fields = append(fields, zap.String("process_id", reqID))
	}
	switch level {
	case "info":
		instance.Info(message, fields...)
	case "error":
		instance.Error(message, fields...)
	case "warn":
		instance.Warn(message, fields...)
	case "panic":
		instance.Panic(message, fields...)
	case "fatal":
		instance.Fatal(message, fields...)
	case "debug":
		instance.Debug(message, fields...)
	}

	if env == LOGGER_ENV_SETUP_DEVELOPMENT_VALUE {
		writeMessageLog(level, message, fields)
	}
}

func writeMessageLog(level string, message string, fields []zap.Field) {
	var messageLog = messageLog{
		Timestamps: time.Now(),
		Level:      level,
		Message:    message,
		Fields:     fields,
	}
	logJson, _ := json.Marshal(messageLog)
	l.Write(fmt.Appendf([]byte{}, "%s\n", logJson))
}

func TrafficLogInfo(ctx context.Context, message string, fields ...zap.Field) {
	if reqID, ok := ctx.Value(CtxRequestID).(string); ok {
		fields = append(fields, zap.String("request_id", reqID))
	}
	fields = append(fields, zap.String("tag", "traffic-log"))
	instance.Info(message, fields...)

	if env == LOGGER_ENV_SETUP_DEVELOPMENT_VALUE {
		var messageLog = messageLog{
			Timestamps: time.Now(),
			Level:      "info",
			Message:    message,
			Fields:     fields,
		}
		logJson, _ := json.Marshal(messageLog)

		timeString := time.Now().Format("2006-01-02")
		l2 := &lumberjack.Logger{
			Filename:   fmt.Sprintf("%straffic-%s.log", logPath, timeString),
			MaxSize:    500, // Megabytes
			MaxBackups: 3,
			MaxAge:     14,   // Days
			Compress:   true, // Disabled by default
		}
		l2.Write(fmt.Appendf([]byte{}, "%s\n", logJson))
	}
}

func LogInfo(ctx context.Context, message string, fields ...zap.Field) {
	CommonLog(ctx, "info", message, fields...)
}

func LogError(ctx context.Context, message string, fields ...zap.Field) {
	CommonLog(ctx, "error", message, fields...)
}

func LogPanic(ctx context.Context, message string, fields ...zap.Field) {
	CommonLog(ctx, "panic", message, fields...)
}

func LogFatal(ctx context.Context, message string, fields ...zap.Field) {
	CommonLog(ctx, "fatal", message, fields...)
}

func LogWarn(ctx context.Context, message string, fields ...zap.Field) {
	CommonLog(ctx, "warn", message, fields...)
}

func LogDebug(ctx context.Context, message string, fields ...zap.Field) {
	CommonLog(ctx, "debug", message, fields...)
}
