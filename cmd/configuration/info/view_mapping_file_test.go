package info

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestViewMappingFile tests the ViewMappingFile function
func TestViewMappingFile(t *testing.T) {
	// Create a new view mapping file command
	cmd := ViewMappingFile()

	// Verify that the command is not nil
	assert.NotNil(t, cmd, "Command should not be nil")

	// Verify that the command has the expected properties
	assert.Equal(t, "map-file", cmd.Use, "Command should have the correct use")
	assert.Contains(t, cmd.Aliases, "mf", "Command should have the correct alias")
	assert.Contains(t, cmd.Aliases, "tf", "Command should have the correct alias")
	assert.NotEmpty(t, cmd.Short, "Command should have a short description")
	assert.NotEmpty(t, cmd.Long, "Command should have a long description")
	assert.NotEmpty(t, cmd.Example, "Command should have examples")
	assert.True(t, cmd.SilenceUsage, "Command should silence usage")
	assert.True(t, cmd.SilenceErrors, "Command should silence errors")

	// Verify that the command has the RunE function
	assert.NotNil(t, cmd.RunE, "Command should have a RunE function")

	// Verify that the command has the expected flags
	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag, "Command should have the JSON flag")

	realmFlag := cmd.Flags().Lookup("realm")
	assert.NotNil(t, realmFlag, "Command should have the realm flag")
}

// TestRunViewFileMappingCommand tests the RunViewFileMappingCommand function
// This is a smoke test since we can't easily test the actual execution
func TestRunViewFileMappingCommand(t *testing.T) {
	// Create a new command for testing
	cmd := &cobra.Command{
		Use: "test",
	}

	// Add the flags that the function expects
	cmd.Flags().Bool("json", false, "")
	cmd.Flags().String("realm", "", "")

	// This is a smoke test since we can't easily test the actual execution
	// In a real test environment, we would mock the info.ViewConfiguration function
	// For now, we'll just verify that the function doesn't panic
	assert.NotPanics(t, func() {
		// We expect an error since the tenancy mapping file is not available in the test environment
		_ = RunViewFileMappingCommand(cmd)
	}, "RunViewFileMappingCommand should not panic")
}
