package helpers

import (
	"log/slog"
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
)

// setupTest prepares the test environment and returns a cleanup function
func setupTest(t *testing.T) func() {
	// Save original values
	originalLogLevel := LogLevel
	originalColoredOutput := ColoredOutput
	originalLogger := Logger
	originalCmdLogger := CmdLogger

	// Return a cleanup function
	return func() {
		// Restore original values
		LogLevel = originalLogLevel
		ColoredOutput = originalColoredOutput
		Logger = originalLogger
		CmdLogger = originalCmdLogger
	}
}

// TestGetSlogLevel tests the getSlogLevel function
func TestGetSlogLevel(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedLevel slog.Level
		expectError   bool
	}{
		{
			name:          "debug level",
			input:         "debug",
			expectedLevel: slog.LevelDebug,
			expectError:   false,
		},
		{
			name:          "info level",
			input:         "info",
			expectedLevel: slog.LevelInfo,
			expectError:   false,
		},
		{
			name:          "warn level",
			input:         "warn",
			expectedLevel: slog.LevelWarn,
			expectError:   false,
		},
		{
			name:          "error level",
			input:         "error",
			expectedLevel: slog.LevelError,
			expectError:   false,
		},
		{
			name:          "uppercase DEBUG",
			input:         "DEBUG",
			expectedLevel: slog.LevelDebug,
			expectError:   false,
		},
		{
			name:          "mixed case InFo",
			input:         "InFo",
			expectedLevel: slog.LevelInfo,
			expectError:   false,
		},
		{
			name:          "invalid level",
			input:         "invalid",
			expectedLevel: slog.LevelDebug, // Default return value on error
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := getSlogLevel(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedLevel, level)
		})
	}
}

// TestGetKlogLevel tests the getKlogLevel function
func TestGetKlogLevel(t *testing.T) {
	tests := []struct {
		name          string
		input         slog.Level
		expectedLevel slog.Level
	}{
		{
			name:          "debug level",
			input:         slog.LevelDebug,
			expectedLevel: slog.LevelDebug,
		},
		{
			name:          "info level",
			input:         slog.LevelInfo,
			expectedLevel: slog.LevelError,
		},
		{
			name:          "warn level",
			input:         slog.LevelWarn,
			expectedLevel: slog.LevelError,
		},
		{
			name:          "error level",
			input:         slog.LevelError,
			expectedLevel: slog.LevelError,
		},
		{
			name:          "custom level below info",
			input:         slog.LevelInfo - 1,
			expectedLevel: slog.LevelInfo - 1,
		},
		{
			name:          "custom level above info",
			input:         slog.LevelInfo + 1,
			expectedLevel: slog.LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := getKlogLevel(tt.input)
			assert.Equal(t, tt.expectedLevel, level)
		})
	}
}

// TestInitLogger tests the InitLogger function
func TestInitLogger(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Test with nil logger
	InitLogger(logr.Logger{})
	assert.NotNil(t, Logger)
	assert.NotNil(t, Logger.GetSink())

	// Test with valid logger
	// Create a mock logger
	slogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
	mockLogger := logr.FromSlogHandler(slogger.Handler())

	InitLogger(mockLogger)
	assert.Equal(t, mockLogger, Logger)
}

// TestSetLogger tests the SetLogger function
func TestSetLogger(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

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

	// Test with debug log level
	LogLevel = "debug"
	err = SetLogger()
	assert.NoError(t, err)

	// Test with colored output
	LogLevel = "info"
	ColoredOutput = true
	err = SetLogger()
	assert.NoError(t, err)
}
