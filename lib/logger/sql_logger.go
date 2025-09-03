package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

const CtxRepoName string = "X-Repo-Name"

type SQLLogger struct {
	logger.Interface
	Env       string
	DebugMode bool
}

func (l *SQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rowsAffected := fc()
	duration := time.Since(begin)
	repoName := ""
	if ctxRepo := ctx.Value(CtxRepoName); ctxRepo != nil {
		if repo, ok := ctxRepo.(string); ok {
			repoName = repo
		}
	}
	logCategory, logLevel, isShouldLog := l.determineLogLevel(duration, err)

	if isShouldLog {
		fields := []zap.Field{
			zap.String("sql", sql),
			zap.Int64("rows_affected", rowsAffected),
			zap.String("env", l.Env),
			zap.String("repo", repoName),
			zap.String("duration", duration.String()),
			zap.String("query_category", logCategory),
			zap.Strings("tags", []string{"repo", repoName}),
		}
		if err != nil {
			fields = append(fields, zap.Error(err))
		}
		l.logWithLevel(ctx, logLevel, "SQL Execution", fields...)
	}

	l.Interface.Trace(ctx, begin, fc, err)
}

// determineLogLevel is helper function to determine log category and level
func (l *SQLLogger) determineLogLevel(duration time.Duration, err error) (category, level string, shouldLog bool) {
	switch {
	case err != nil:
		return "error", "error", true
	case duration >= 200*time.Millisecond:
		return "slow_query", "warn", true
	case l.Env == "development" || l.DebugMode:
		return "info", "info", true
	default:
		return "", "", false
	}
}

func (l *SQLLogger) logWithLevel(ctx context.Context, level, msg string, fields ...zap.Field) {
	switch level {
	case "error":
		LogError(ctx, msg, fields...)
	case "warn":
		LogWarn(ctx, msg, fields...)
	case "info":
		LogInfo(ctx, msg, fields...)
	}
}
