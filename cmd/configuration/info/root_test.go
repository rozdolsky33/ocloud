package info

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewInfoCmd tests the NewInfoCmd function
func TestNewInfoCmd(t *testing.T) {
	// Create a new info command
	cmd := NewInfoCmd()

	// Verify that the command is not nil
	assert.NotNil(t, cmd, "Command should not be nil")

	// Verify that the command has the expected properties
	assert.Equal(t, "info", cmd.Use, "Command should have the correct use")
	assert.Contains(t, cmd.Aliases, "i", "Command should have the correct alias")
	assert.NotEmpty(t, cmd.Short, "Command should have a short description")
	assert.NotEmpty(t, cmd.Long, "Command should have a long description")
	assert.NotEmpty(t, cmd.Example, "Command should have examples")
	assert.True(t, cmd.SilenceUsage, "Command should silence usage")
	assert.True(t, cmd.SilenceErrors, "Command should silence errors")

	// Verify that the command has subcommands
	assert.NotEmpty(t, cmd.Commands(), "Command should have subcommands")

	// Verify that the view-mapping-file command is a subcommand
	var hasViewMappingFile bool
	for _, subCmd := range cmd.Commands() {
		if subCmd.Use == "map-file" {
			hasViewMappingFile = true
			break
		}
	}
	assert.True(t, hasViewMappingFile, "Command should have the view-mapping-file subcommand")
}
