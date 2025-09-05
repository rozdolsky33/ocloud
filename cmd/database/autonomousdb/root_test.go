package autonomousdb

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestAutonomousDatabaseCommand tests the basic structure of the autonomousdb command
func TestAutonomousDatabaseCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new autonomousdb command
	cmd := NewAutonomousDatabaseCmd(appCtx)

	// Test that the autonomousdb command is properly configured
	assert.Equal(t, "autonomous", cmd.Use)
	assert.Equal(t, "Manage OCI Compartments", cmd.Short)
	assert.Equal(t, "Manage Oracle Cloud Infrastructure Databases - list all databases or find database by pattern.", cmd.Long)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the subcommands are added
	subCmds := cmd.Commands()
	assert.Equal(t, 2, len(subCmds), "autonomousdb command should have 2 subcommands")

	// Check that the list subcommand is present
	getCmd := findSubCommand(subCmds, "get")
	assert.NotNil(t, getCmd, "autonomousdb command should have list subcommand")

	// Check that the find subcommand is present
	findCmd := findSubCommand(subCmds, "find")
	assert.NotNil(t, findCmd, "autonomousdb command should have find subcommand")
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
