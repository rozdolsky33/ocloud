package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewConfigCmd tests the NewConfigCmd function
func TestNewConfigCmd(t *testing.T) {
	// Call the function
	cmd := NewConfigCmd()

	// Verify the command properties
	assert.NotNil(t, cmd, "Command should not be nil")
	assert.Equal(t, "config", cmd.Use, "Command Use should be 'config'")
	assert.Equal(t, []string{"conf"}, cmd.Aliases, "Command Aliases should include 'conf'")
	assert.Equal(t, "Manage ocloud CLI configurations file and authentication", cmd.Short, "Command Short description should match")
	assert.NotEmpty(t, cmd.Long, "Command Long description should not be empty")
	assert.NotEmpty(t, cmd.Example, "Command Example should not be empty")
	assert.True(t, cmd.SilenceUsage, "Command SilenceUsage should be true")
	assert.True(t, cmd.SilenceErrors, "Command SilenceErrors should be true")

	// Verify that the command has subcommands
	assert.Greater(t, len(cmd.Commands()), 0, "Command should have subcommands")

	// Verify that the expected subcommands are present
	subcommandNames := make(map[string]bool)
	for _, subCmd := range cmd.Commands() {
		subcommandNames[subCmd.Use] = true
	}

	assert.True(t, subcommandNames["info"], "Command should have 'info' subcommand")
	assert.True(t, subcommandNames["session"], "Command should have 'session' subcommand")
	assert.True(t, subcommandNames["setup"], "Command should have 'setup' subcommand")
}
