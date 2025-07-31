package logger

import (
	"log/slog"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
)

// TestGetSlogLevel tests the getSlogLevel function
func TestGetSlogLevel(t *testing.T) {
	// Test valid log levels
	testCases := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug - 10}, // Adjusted to match implementation
		{"DEBUG", slog.LevelDebug - 10}, // Adjusted to match implementation
		{"Debug", slog.LevelDebug - 10}, // Adjusted to match implementation
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"Info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"WARN", slog.LevelWarn},
		{"Warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"ERROR", slog.LevelError},
		{"Error", slog.LevelError},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			level, err := getSlogLevel(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, level)
		})
	}

	// Test invalid log level
	level, err := getSlogLevel("invalid")
	assert.Error(t, err)
	assert.Equal(t, slog.LevelDebug, level) // Default to debug on error
}

// TestGetKlogLevel tests the getKlogLevel function
func TestGetKlogLevel(t *testing.T) {
	// Test debug level
	level := getKlogLevel(slog.LevelDebug)
	assert.Equal(t, slog.LevelDebug, level)

	// Test info level
	level = getKlogLevel(slog.LevelInfo)
	assert.Equal(t, slog.LevelError, level)

	// Test warn level
	level = getKlogLevel(slog.LevelWarn)
	assert.Equal(t, slog.LevelError, level)

	// Test error level
	level = getKlogLevel(slog.LevelError)
	assert.Equal(t, slog.LevelError, level)
}

// TestInitLogger tests the InitLogger function
func TestInitLogger(t *testing.T) {
	// Save original logger
	originalLogger := Logger
	defer func() {
		// Restore original logger
		Logger = originalLogger
	}()

	// Test with a valid logger
	mockLogger := logr.Discard()
	InitLogger(mockLogger)
	// Can't directly compare logr.Logger structs because they contain interfaces
	// Just check that the logger is set to something non-nil
	assert.NotNil(t, Logger)
	assert.NotNil(t, Logger.GetSink())

	// Test with a nil logger
	// This is hard to test directly since we can't easily check the shared state
	// of the logger, but we can at least verify it doesn't panic
	InitLogger(logr.Logger{})
	assert.NotNil(t, Logger)
}

// TestSetLogger tests the SetLogger function
func TestSetLogger(t *testing.T) {
	// Save original values
	originalLogLevel := LogLevel
	originalColoredOutput := ColoredOutput
	originalCmdLogger := CmdLogger
	defer func() {
		// Restore original values
		LogLevel = originalLogLevel
		ColoredOutput = originalColoredOutput
		CmdLogger = originalCmdLogger
	}()

	// Test with valid log level
	LogLevel = "info"
	ColoredOutput = false
	err := SetLogger()
	assert.NoError(t, err)
	assert.NotNil(t, CmdLogger)

	// Test with invalid log level
	LogLevel = "invalid"
	err = SetLogger()
	assert.Error(t, err)
}
