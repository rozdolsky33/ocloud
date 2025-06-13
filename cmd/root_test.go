package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestRootCommand tests the basic structure of the root command
func TestRootCommand(t *testing.T) {
	// Create a mock AppContext
	appCtx := &app.AppContext{}

	// Create a new root command
	rootCmd := NewRootCmd(appCtx)

	// Test that the root command is properly configured
	assert.Equal(t, "ocloud", rootCmd.Use)
	assert.Equal(t, "Interact with Oracle Cloud Infrastructure", rootCmd.Short)
	assert.True(t, rootCmd.SilenceUsage)

	// Test that the instance command is added as a subcommand
	instanceCmd := findSubcommand(rootCmd, "instance")
	assert.NotNil(t, instanceCmd, "instance command should be added as a subcommand")
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
