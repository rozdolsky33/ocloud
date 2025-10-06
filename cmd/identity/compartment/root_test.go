package compartment

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestCompartmentCommand tests the basic structure of the compartment command
func TestCompartmentCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new compartment command
	cmd := NewCompartmentCmd(appCtx)

	// Test that the compartment command is properly configured
	assert.Equal(t, "compartment", cmd.Use)
	assert.Equal(t, []string{"compart", "comp", "cmp", "c"}, cmd.Aliases)
	assert.Equal(t, "Manage OCI Compartments", cmd.Short)
	assert.Equal(t, "Manage Oracle Cloud Infrastructure Compartments: list, get and search", cmd.Long)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the subcommands are added
	subCmds := cmd.Commands()
	assert.Equal(t, 3, len(subCmds), "compartment command should have 3 subcommands")

	// Check that the list subcommand is present
	listCmd := compartmentSubCommand(subCmds, "list")
	assert.NotNil(t, listCmd, "compartment command should have list subcommand")

	// Check that the get subcommand is present
	getCmd := compartmentSubCommand(subCmds, "get")
	assert.NotNil(t, getCmd, "compartment command should have get subcommand")

	// Check that the find subcommand is present
	findCmd := compartmentSubCommand(subCmds, "search")
	assert.NotNil(t, findCmd, "compartment command should have search subcommand")
}

// compartmentSubCommand is a helper function to find a subcommand by name
func compartmentSubCommand(cmds []*cobra.Command, name string) *cobra.Command {
	for _, cmd := range cmds {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
