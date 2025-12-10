package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New(level string) *slog.Logger {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

func WithItemID(logger *slog.Logger, itemID string) *slog.Logger {
	return logger.With("item_id", itemID)
}

func WithTraceID(logger *slog.Logger, traceID string) *slog.Logger {
	return logger.With("trace_id", traceID)
}

func WithTenantID(logger *slog.Logger, tenantID string) *slog.Logger {
	return logger.With("tenant_id", tenantID)
}
