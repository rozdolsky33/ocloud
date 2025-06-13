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

	// GLOBAL_VERBOSITY controls the verbosity level for V(n) calls
	// Higher values show more verbose logs
	GLOBAL_VERBOSITY int
)

// SetLogger initializes the loggers based on the current LogLevel and ColoredOutput settings
func SetLogger() error {
	l, err := getSlogLevel(LogLevel)
	if err != nil {
		return err
	}

	// Set the global verbosity level based on the log level
	switch strings.ToLower(LogLevel) {
	case "debug":
		// Turn on all verbose levels 0..5
		GLOBAL_VERBOSITY = 5
	default:
		GLOBAL_VERBOSITY = 0
	}

	slogger := slog.New(logger.NewHandler(os.Stderr, logger.Options{Level: l, Colored: ColoredOutput}))
	kslogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: getKlogLevel(l)}))
	baseLogger := logr.FromSlogHandler(slogger.Handler())
	klogger := logr.FromSlogHandler(kslogger.Handler())

	klog.SetLogger(klogger)
	ctrl.SetLogger(baseLogger)
	CmdLogger = baseLogger

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

// VerboseInfo logs a message at the specified verbosity level.
// If the verbosity level is less than or equal to GLOBAL_VERBOSITY,
// it logs the message using the logger's V(level).Info() method.
// Otherwise, it does nothing.
func VerboseInfo(logger logr.Logger, level int, msg string, keysAndValues ...interface{}) {
	if level <= GLOBAL_VERBOSITY {
		logger.V(level).Info(msg, keysAndValues...)
	}
}

// getSlogLevel converts a string log level to a slog.Level
func getSlogLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		// Set to a lower level than LevelDebug to ensure all debug logs are shown
		return slog.LevelDebug - 10, nil
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
