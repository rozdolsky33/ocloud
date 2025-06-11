package config

import (
	"log/slog"
	"os"

	"github.com/go-logr/logr"
)

var (
	// Logger is the package-level logger
	Logger logr.Logger
)

// InitLogger initializes the package-level logger
func InitLogger(logger logr.Logger) {
	Logger = logger
	if Logger.GetSink() == nil {
		// If no logger is provided, create a default one
		slogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
		Logger = logr.FromSlogHandler(slogger.Handler())
	}
}

func init() {
	// Initialize with a default logger
	slogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	Logger = logr.FromSlogHandler(slogger.Handler())
}