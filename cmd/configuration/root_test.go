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
	assert.Equal(t, "Manage ocloud CLI configurations", cmd.Short)
	assert.Equal(t, "Manage ocloud CLI configurations with OCI such as authentication, view configuration information, and more.", cmd.Long)
	assert.Contains(t, cmd.Example, "ocloud config info map-file")
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the subcommands are added
	subCmds := cmd.Commands()
	assert.Equal(t, 1, len(subCmds), "configuration command should have 1 subcommand")

	// Check that the info subcommand is present
	infoCmd := findSubCommand(subCmds, "info")
	assert.NotNil(t, infoCmd, "configuration command should have info subcommand")
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
