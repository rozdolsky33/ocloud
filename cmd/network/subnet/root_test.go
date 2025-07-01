package subnet

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestSubnetCommand tests the basic structure of the subnet command
func TestSubnetCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new subnet command
	cmd := NewSubnetCmd(appCtx)

	// Test that the subnet command is properly configured
	assert.Equal(t, "subnet", cmd.Use)
	assert.Equal(t, []string{"sub"}, cmd.Aliases)
	assert.Equal(t, "Manage OCI Subnets", cmd.Short)
	assert.Equal(t, "Manage Oracle Cloud Infrastructure Subnets - list all subnets or find subnet by pattern.", cmd.Long)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the subcommands are added
	subCmds := cmd.Commands()
	assert.Equal(t, 2, len(subCmds), "subnet command should have 2 subcommands")

	// Check that the list subcommand is present
	listCmd := findSubCommand(subCmds, "list")
	assert.NotNil(t, listCmd, "subnet command should have list subcommand")

	// Check that the find subcommand is present
	findCmd := findSubCommand(subCmds, "find")
	assert.NotNil(t, findCmd, "subnet command should have find subcommand")
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
