package heatwave

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestHeatWaveDatabaseCommand tests the basic structure of the heatwave command
func TestHeatWaveDatabaseCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new heatwave command
	cmd := NewHeatWaveDatabaseCmd(appCtx)

	// Test that the heatwave command is properly configured
	assert.Equal(t, "heatwave", cmd.Use)
	assert.Equal(t, []string{"hw"}, cmd.Aliases)
	assert.Equal(t, "Explore OCI HeatWave Databases.", cmd.Short)
	assert.Equal(t, "Explore Oracle Cloud Infrastructure databases: list, get, and search", cmd.Long)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the subcommands are added
	subCmds := cmd.Commands()
	assert.Equal(t, 3, len(subCmds), "heatwave command should have 3 subcommands")

	// Check that the get subcommand is present
	getCmd := hwSubCommand(subCmds, "get")
	assert.NotNil(t, getCmd, "heatwave command should have get subcommand")

	// Check that the list subcommand is present
	listCmd := hwSubCommand(subCmds, "list")
	assert.NotNil(t, listCmd, "heatwave command should have list subcommand")

	// Check that the search subcommand is present
	searchCmd := hwSubCommand(subCmds, "search")
	assert.NotNil(t, searchCmd, "heatwave command should have search subcommand")
}

// hwSubCommand is a helper function to search a subcommand by name
func hwSubCommand(cmds []*cobra.Command, name string) *cobra.Command {
	for _, cmd := range cmds {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
