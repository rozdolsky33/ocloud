package configuration

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
)

// TestConfigCommand tests the basic structure of the configuration command
func TestConfigCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new configuration command
	cmd := NewConfigCmd(appCtx)

	// Test that the configuration command is properly configured
	assert.Equal(t, "config", cmd.Use)
	assert.Equal(t, []string{"conf"}, cmd.Aliases)
	assert.Equal(t, "Manage ocloud CLI configurations and authentication", cmd.Short)
	assert.Equal(t, "Manage ocloud CLI configurations and authentication with Oracle Cloud Infrastructure (OCI).\n\nThis command group provides functionality for:\n- Authenticating with OCI and refreshing session tokens\n- Viewing configuration information such as tenancy mappings", cmd.Long)
	assert.Contains(t, cmd.Example, "ocloud config session")
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the subcommands are added
	subCmds := cmd.Commands()
	assert.Equal(t, 2, len(subCmds), "configuration command should have 2 subcommands")

	// Check that the info subcommand is present
	infoCmd := findSubCommand(subCmds, "info")
	assert.NotNil(t, infoCmd, "configuration command should have info subcommand")

	// Check that the session subcommand is present
	sessionCmd := findSubCommand(subCmds, "session")
	assert.NotNil(t, sessionCmd, "configuration command should have session subcommand")
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
