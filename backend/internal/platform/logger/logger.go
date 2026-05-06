package logger

import (
	"log/slog"
	"os"
)

// New returns a centralized app logger with INFO level.
func New() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler)
}
