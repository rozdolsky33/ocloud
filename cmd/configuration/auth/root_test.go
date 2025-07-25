package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewSessionCmd tests the NewSessionCmd function
func TestNewSessionCmd(t *testing.T) {
	// Create a new session command
	cmd := NewSessionCmd()

	// Verify that the command is not nil
	assert.NotNil(t, cmd, "Command should not be nil")

	// Verify that the command has the expected properties
	assert.Equal(t, "session", cmd.Use, "Command should have the correct use")
	assert.Contains(t, cmd.Aliases, "s", "Command should have the correct alias")
	assert.Equal(t, sessionShort, cmd.Short, "Command should have the correct short description")
	assert.Equal(t, sessionLong, cmd.Long, "Command should have the correct long description")
	assert.Equal(t, sessionExamples, cmd.Example, "Command should have the correct examples")
	assert.True(t, cmd.SilenceUsage, "Command should silence usage")
	assert.True(t, cmd.SilenceErrors, "Command should silence errors")

	// Verify that the RunE function is set
	assert.NotNil(t, cmd.RunE, "Command should have a RunE function")

	// Verify that the command has subcommands
	assert.NotEmpty(t, cmd.Commands(), "Command should have subcommands")

	// Verify that the authenticate command is a subcommand
	var hasAuthenticateCmd bool
	for _, subCmd := range cmd.Commands() {
		if subCmd.Use == "authenticate" {
			hasAuthenticateCmd = true
			break
		}
	}
	assert.True(t, hasAuthenticateCmd, "Command should have the authenticate subcommand")
}
