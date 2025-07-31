package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/cmd/shared/cmdcreate"
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
	rootCmd := cmdcreate.CreateRootCmdWithoutContext()

	// Test that the root command is properly configured
	assert.Equal(t, "ocloud", rootCmd.Use)
	assert.Equal(t, "Interact with Oracle Cloud Infrastructure", rootCmd.Short)
	assert.True(t, rootCmd.SilenceUsage)

	// Test that the version command is added as a subcommand
	versionCmd := findSubcommand(rootCmd, "version")
	assert.NotNil(t, versionCmd, "version command should be added as a subcommand")

	// Test that placeholder commands are added
	computeCmd := findSubcommand(rootCmd, "compute")
	assert.NotNil(t, computeCmd, "compute command should be added as a placeholder subcommand")

	// Verify it's a placeholder by checking that it returns an error when run
	err := computeCmd.RunE(computeCmd, []string{})
	assert.Error(t, err, "placeholder command should return an error when run")
	assert.Contains(t, err.Error(), "requires application initialization", "error message should indicate initialization is required")
}

// TestIsNoContextCommand tests the isNoContextCommand function
func TestIsNoContextCommand(t *testing.T) {
	// Save original os.Args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test with version command
	os.Args = []string{"ocloud", "version"}
	assert.True(t, isNoContextCommand(), "should return true for 'version' command")

	// Test with config command
	os.Args = []string{"ocloud", "config"}
	assert.True(t, isNoContextCommand(), "should return true for 'config' command")

	// Test with a version flag (short)
	os.Args = []string{"ocloud", "-v"}
	assert.True(t, isNoContextCommand(), "should return true for '-v' flag")

	// Test with a version flag (long)
	os.Args = []string{"ocloud", "--version"}
	assert.True(t, isNoContextCommand(), "should return true for '--version' flag")

	// Test with another command
	os.Args = []string{"ocloud", "compute", "instance", "list"}
	assert.False(t, isNoContextCommand(), "should return false for other commands")

	// Test with no arguments
	os.Args = []string{"ocloud"}
	assert.True(t, isNoContextCommand(), "should return true when no arguments are provided (just the program name)")
}
