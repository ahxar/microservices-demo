package observability

import (
	"log/slog"
	"os"
)

// SetupDefaultLogger configures the process logger to emit structured JSON logs.
func SetupDefaultLogger(serviceName string) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).With(
		slog.String("service", serviceName),
	)
	slog.SetDefault(logger)
	return logger
}
