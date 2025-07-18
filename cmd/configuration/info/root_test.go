package info

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestInfoCommand tests the basic structure of the info command
func TestInfoCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new info command
	cmd := NewInfoCmd(appCtx)

	// Test that the info command is properly configured
	assert.Equal(t, "info", cmd.Use)
	assert.Equal(t, []string{"i"}, cmd.Aliases)
	assert.Equal(t, "View information about ocloud environment configuration", cmd.Short)
	assert.Equal(t, "View information about ocloud environment configuration, such as tenancy mappings and other configuration details.", cmd.Long)
	assert.Contains(t, cmd.Example, "ocloud config info map-file")
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the subcommands are added
	subCmds := cmd.Commands()
	assert.Equal(t, 1, len(subCmds), "info command should have 1 subcommand")

	// Check that the map-file subcommand is present
	mapFileCmd := findSubCommand(subCmds, "map-file")
	assert.NotNil(t, mapFileCmd, "info command should have map-file subcommand")
}

// findSubCommand is a helper function to find a subcommand by name
func findSubCommand(cmds []*cobra.Command, name string) *cobra.Command {
	for _, cmd := range cmds {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
