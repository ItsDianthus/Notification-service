package slogger

import (
	"log/slog"
	"os"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

func NewOperationLogger(logger *slog.Logger, operation string) *slog.Logger {
	return logger.With(slog.String("operation", operation))
}

func NewLoggerByEnvironment(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case EnvLocal:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	case EnvProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	return slog.New(handler)
}

func ErrAttr(err error) slog.Attr {
	if err == nil {
		return slog.Attr{Key: "error", Value: slog.StringValue("<nil>")}
	}
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
