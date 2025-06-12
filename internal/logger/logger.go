package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/logger"
	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	// Logger is the package-level logger
	// It should be initialized using InitLogger before use
	Logger logr.Logger

	// LogLevel sets the verbosity level for logging
	LogLevel string

	// LogLevelMsg provides help text for the log-level flag
	LogLevelMsg = "Set the log verbosity. Supported values are: debug, info, warn, and error."

	// CmdLogger is the logger used by command-line operations
	CmdLogger logr.Logger

	// ColoredOutput determines whether log output should be colored
	ColoredOutput bool

	// ColoredOutputMsg provides help text for the color flag
	ColoredOutputMsg = "Enable colored log messages."
)

// SetLogger initializes the loggers based on the current LogLevel and ColoredOutput settings
func SetLogger() error {
	l, err := getSlogLevel(LogLevel)
	if err != nil {
		return err
	}

	slogger := slog.New(logger.NewHandler(os.Stderr, logger.Options{Level: l, Colored: ColoredOutput}))
	kslogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: getKlogLevel(l)}))
	logger := logr.FromSlogHandler(slogger.Handler())
	klogger := logr.FromSlogHandler(kslogger.Handler())

	klog.SetLogger(klogger)
	ctrl.SetLogger(logger)
	CmdLogger = logger
	return nil
}

// InitLogger initializes the package-level logger
// If no logger is provided, it creates a default one
func InitLogger(logger logr.Logger) {
	Logger = logger
	if Logger.GetSink() == nil {
		// If no logger is provided, or it has a nil sink, create a default one
		slogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
		Logger = logr.FromSlogHandler(slogger.Handler())
	}
}

// getSlogLevel converts a string log level to a slog.Level
func getSlogLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelDebug, fmt.Errorf("%s is not a valid log level", s)
	}
}

// For end users, klog messages are mostly useless. We set it to the error level unless debug logging is enabled.
func getKlogLevel(l slog.Level) slog.Level {
	if l < slog.LevelInfo {
		return l
	}
	return slog.LevelError
}