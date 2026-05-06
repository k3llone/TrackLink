package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New returns a centralized app logger with INFO level.
func New(format string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	format = strings.ToLower(strings.TrimSpace(format))

	switch format {
	case "json":
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	case "text":
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
}
