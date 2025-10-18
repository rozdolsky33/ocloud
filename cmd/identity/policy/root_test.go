package policy

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestPolicyCommand tests the basic structure of the policy command
func TestPolicyCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new policy command
	cmd := NewPolicyCmd(appCtx)

	// Test that the policy command is properly configured
	assert.Equal(t, "policy", cmd.Use)
	assert.Equal(t, []string{"pol"}, cmd.Aliases)
	assert.Equal(t, "Explore OCI Policies", cmd.Short)
	assert.Equal(t, "Explore Oracle Cloud Infrastructure Policies: list, get, and search", cmd.Long)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the subcommands are added
	subCmds := cmd.Commands()
	assert.Equal(t, 3, len(subCmds), "policy command should have 2 subcommands")

	// Check that the list subcommand is present
	listCmd := policySubCommand(subCmds, "get")
	assert.NotNil(t, listCmd, "policy command should have get subcommand")

	// Check that the find subcommand is present
	findCmd := policySubCommand(subCmds, "search")
	assert.NotNil(t, findCmd, "policy command should have find subcommand")
}

// policySubCommand is a helper function to find a subcommand by name
func policySubCommand(cmds []*cobra.Command, name string) *cobra.Command {
	for _, cmd := range cmds {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
