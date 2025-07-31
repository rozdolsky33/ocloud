package cmdutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsNoContextCommand tests the IsNoContextCommand function
func TestIsNoContextCommand(t *testing.T) {
	// Save original os.Args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test with version command
	os.Args = []string{"ocloud", "version"}
	assert.True(t, IsNoContextCommand(), "should return true for 'version' command")

	// Test with config command
	os.Args = []string{"ocloud", "config"}
	assert.True(t, IsNoContextCommand(), "should return true for 'config' command")

	// Test with version flag (short)
	os.Args = []string{"ocloud", "-v"}
	assert.True(t, IsNoContextCommand(), "should return true for '-v' flag")

	// Test with version flag (long)
	os.Args = []string{"ocloud", "--version"}
	assert.True(t, IsNoContextCommand(), "should return true for '--version' flag")

	// Test with other command
	os.Args = []string{"ocloud", "compute", "instance", "list"}
	assert.False(t, IsNoContextCommand(), "should return false for other commands")

	// Test with no arguments
	os.Args = []string{"ocloud"}
	assert.True(t, IsNoContextCommand(), "should return true when no arguments are provided (just the program name)")
}

// TestIsRootCommandWithoutSubcommands tests the IsRootCommandWithoutSubcommands function
func TestIsRootCommandWithoutSubcommands(t *testing.T) {
	// Save original os.Args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test with no subcommands
	os.Args = []string{"ocloud"}
	assert.True(t, IsRootCommandWithoutSubcommands(), "should return true when no subcommands are provided")

	// Test with subcommand
	os.Args = []string{"ocloud", "compute"}
	assert.False(t, IsRootCommandWithoutSubcommands(), "should return false when a subcommand is provided")

	// Test with flag
	os.Args = []string{"ocloud", "--version"}
	assert.False(t, IsRootCommandWithoutSubcommands(), "should return false when a flag is provided")

	// Test with multiple arguments
	os.Args = []string{"ocloud", "compute", "instance", "list"}
	assert.False(t, IsRootCommandWithoutSubcommands(), "should return false when multiple arguments are provided")
}
