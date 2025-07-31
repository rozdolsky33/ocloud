package logger

import (
	"os"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestSetLogLevel tests the SetLogLevel function
func TestSetLogLevel(t *testing.T) {
	// Save original os.Args and logger settings
	originalArgs := os.Args
	originalLogLevel := logger.LogLevel
	originalColoredOutput := logger.ColoredOutput

	// Restore original values after the test
	defer func() {
		os.Args = originalArgs
		logger.LogLevel = originalLogLevel
		logger.ColoredOutput = originalColoredOutput
	}()

	// Create a test root command
	testRoot := &cobra.Command{
		Use:   "ocloud",
		Short: "Test command",
	}

	// Add the flags that SetLogLevel expects
	flags.AddGlobalFlags(testRoot)

	// Test case 1: Default settings (no flags)
	os.Args = []string{"ocloud"}
	err := SetLogLevel(testRoot)
	assert.NoError(t, err, "SetLogLevel should not return an error with default settings")
	assert.Equal(t, flags.FlagValueInfo, logger.LogLevel, "Default log level should be 'info'")
	assert.False(t, logger.ColoredOutput, "Default colored output should be false")

	// Test case 2: Debug flag
	os.Args = []string{"ocloud", "--debug"}
	err = SetLogLevel(testRoot)
	assert.NoError(t, err, "SetLogLevel should not return an error with debug flag")
	assert.Equal(t, flags.FlagNameDebug, logger.LogLevel, "Log level should be 'debug' with debug flag")

	// Test case 3: Short debug flag
	os.Args = []string{"ocloud", "-d"}
	err = SetLogLevel(testRoot)
	assert.NoError(t, err, "SetLogLevel should not return an error with short debug flag")
	assert.Equal(t, flags.FlagNameDebug, logger.LogLevel, "Log level should be 'debug' with short debug flag")

	// Test case 4: Log level flag
	// Reset the debug flag to ensure it doesn't override the log level
	testRoot.PersistentFlags().Set(flags.FlagNameDebug, "false")
	os.Args = []string{"ocloud", "--log-level=warn"}
	err = SetLogLevel(testRoot)
	assert.NoError(t, err, "SetLogLevel should not return an error with log level flag")
	assert.Equal(t, "warn", logger.LogLevel, "Log level should be 'warn' with log level flag")

	// Test case 5: Color flag
	os.Args = []string{"ocloud", "--color"}
	err = SetLogLevel(testRoot)
	assert.NoError(t, err, "SetLogLevel should not return an error with color flag")
	assert.True(t, logger.ColoredOutput, "Colored output should be true with color flag")

	// Test case 6: Multiple flags
	os.Args = []string{"ocloud", "--debug", "--color"}
	err = SetLogLevel(testRoot)
	assert.NoError(t, err, "SetLogLevel should not return an error with multiple flags")
	assert.Equal(t, flags.FlagNameDebug, logger.LogLevel, "Log level should be 'debug' with debug flag")
	assert.True(t, logger.ColoredOutput, "Colored output should be true with color flag")

	// Note: We can't easily test the version flag case because it calls os.Exit(0)
	// In a real test environment, we might use a custom exit function that can be mocked
}

// TestSetLogLevelWithVersionFlag tests the behavior when version flags are present
// This test is skipped because the function calls os.Exit(0) which would terminate the test
func TestSetLogLevelWithVersionFlag(t *testing.T) {
	t.Skip("Skipping test because SetLogLevel calls os.Exit(0) when version flags are present")

	// In a real test environment, we might use a custom exit function that can be mocked
	// For example:
	//
	// // Save original os.Args and exit function
	// originalArgs := os.Args
	// originalExit := exitFunc
	// exitCalled := false
	//
	// // Mock the exit function
	// exitFunc = func(code int) {
	//     exitCalled = true
	//     assert.Equal(t, 0, code, "Exit code should be 0")
	// }
	//
	// // Restore original values after the test
	// defer func() {
	//     os.Args = originalArgs
	//     exitFunc = originalExit
	// }()
	//
	// // Test with version flag
	// os.Args = []string{"ocloud", "--version"}
	// SetLogLevel(testRoot)
	// assert.True(t, exitCalled, "Exit should be called with version flag")
}
