package logger

import (
	"io"
	"log/slog"

	"github.com/go-logr/logr"
)

// NewTestLogger creates a logger suitable for testing that doesn't produce output.
// It uses a discard handler to ensure no logs are written to stdout/stderr.
func NewTestLogger() logr.Logger {
	// Create a slog.Logger with a handler that discards all output
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug})
	slogger := slog.New(handler)

	// Convert the slog.Logger to a logr.Logger
	return logr.FromSlogHandler(slogger.Handler())
}
