package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestRootCommand tests the basic structure of the root command
func TestRootCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new root command
	rootCmd := NewRootCmd(appCtx)

	// Test that the root command is properly configured
	assert.Equal(t, "ocloud", rootCmd.Use)
	assert.Equal(t, "Interact with Oracle Cloud Infrastructure", rootCmd.Short)
	assert.True(t, rootCmd.SilenceUsage)

	// Test that the compute command is added as a subcommand
	computeCmd := findSubcommand(rootCmd, "compute")
	assert.NotNil(t, computeCmd, "compute command should be added as a subcommand")

	// Test that the instance command is added as a subcommand of the compute command
	if computeCmd != nil {
		instanceCmd := findSubcommand(computeCmd, "instance")
		assert.NotNil(t, instanceCmd, "instance command should be added as a subcommand of the compute command")
	}
}

// findSubcommand is a helper function to find a subcommand by name
func findSubcommand(cmd *cobra.Command, name string) *cobra.Command {
	for _, subCmd := range cmd.Commands() {
		if subCmd.Name() == name {
			return subCmd
		}
	}
	return nil
}

// TestRootCommandWithoutContext tests the root command created without context
func TestRootCommandWithoutContext(t *testing.T) {
	// Create a root command without context
	rootCmd := createRootCmdWithoutContext()

	// Test that the root command is properly configured
	assert.Equal(t, "ocloud", rootCmd.Use)
	assert.Equal(t, "Interact with Oracle Cloud Infrastructure", rootCmd.Short)
	assert.True(t, rootCmd.SilenceUsage)

	// Test that the version command is added as a subcommand
	versionCmd := findSubcommand(rootCmd, "version")
	assert.NotNil(t, versionCmd, "version command should be added as a subcommand")

	// Test that other commands are not added
	computeCmd := findSubcommand(rootCmd, "compute")
	assert.Nil(t, computeCmd, "compute command should not be added as a subcommand")
}

// TestIsNoContextCommand tests the IsNoContextCommand method of the CommandRegistry
func TestIsNoContextCommand(t *testing.T) {
	// Save original os.Args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Create a test registry
	registry := &CommandRegistry{
		NoContextCommands: map[string]bool{
			"version": true,
			"test":    true,
		},
		NoContextFlags: map[string]bool{
			"--version": true,
			"-v":        true,
			"--test":    true,
		},
	}

	// Test with version command
	os.Args = []string{"ocloud", "version"}
	assert.True(t, registry.IsNoContextCommand(), "should return true for 'version' command")

	// Test with test command
	os.Args = []string{"ocloud", "test"}
	assert.True(t, registry.IsNoContextCommand(), "should return true for 'test' command")

	// Test with version flag (short)
	os.Args = []string{"ocloud", "-v"}
	assert.True(t, registry.IsNoContextCommand(), "should return true for '-v' flag")

	// Test with version flag (long)
	os.Args = []string{"ocloud", "--version"}
	assert.True(t, registry.IsNoContextCommand(), "should return true for '--version' flag")

	// Test with test flag
	os.Args = []string{"ocloud", "--test"}
	assert.True(t, registry.IsNoContextCommand(), "should return true for '--test' flag")

	// Test with other command
	os.Args = []string{"ocloud", "compute", "instance", "list"}
	assert.False(t, registry.IsNoContextCommand(), "should return false for other commands")

	// Test with no arguments
	os.Args = []string{"ocloud"}
	assert.False(t, registry.IsNoContextCommand(), "should return false when no arguments are provided")
}

// TestCommandRegistryMethods tests the methods of the CommandRegistry
func TestCommandRegistryMethods(t *testing.T) {
	// Create a test registry
	registry := &CommandRegistry{
		NoContextCommands: map[string]bool{},
		NoContextFlags:    map[string]bool{},
	}

	// Test RegisterNoContextCommand
	registry.RegisterNoContextCommand("test")
	assert.True(t, registry.NoContextCommands["test"], "test command should be registered")

	// Test RegisterNoContextFlag
	registry.RegisterNoContextFlag("--test")
	assert.True(t, registry.NoContextFlags["--test"], "test flag should be registered")
}
