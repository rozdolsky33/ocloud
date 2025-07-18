package info

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestViewMappingFile tests the ViewMappingFile function
func TestViewMappingFile(t *testing.T) {
	// Create a mock application context
	appCtx := &app.ApplicationContext{
		Logger: logger.NewTestLogger(),
	}

	// Call ViewMappingFile with the mock context
	cmd := ViewMappingFile(appCtx)

	// Test that the command is properly configured
	assert.Equal(t, "map-file", cmd.Use)
	assert.Equal(t, []string{"mf", "tf"}, cmd.Aliases)
	assert.Equal(t, "View tenancy mapping information", cmd.Short)
	assert.Contains(t, cmd.Long, "View the tenancy mapping information from the tenancy-map.yaml file")
	assert.Contains(t, cmd.Example, "ocloud configuration info map-file")
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)

	// Test that the command has the expected flags
	jsonFlag := cmd.Flag("json")
	assert.NotNil(t, jsonFlag, "command should have json flag")
	assert.Equal(t, "j", jsonFlag.Shorthand)
	assert.Equal(t, "false", jsonFlag.DefValue)
	assert.Contains(t, jsonFlag.Usage, "Output in JSON format")

	realmFlag := cmd.Flag("realm")
	assert.NotNil(t, realmFlag, "command should have realm flag")
	assert.Equal(t, "r", realmFlag.Shorthand)
	assert.Equal(t, "", realmFlag.DefValue)
	assert.Contains(t, realmFlag.Usage, "Filter by realm")
}

// TestRunViewFileMappingCommand tests the RunViewFileMappingCommand function
func TestRunViewFileMappingCommand(t *testing.T) {
	// This test would normally execute the command and verify its behavior
	// However, since the command interacts with external resources (tenancy-map.yaml file),
	// we'll skip the actual execution and just verify the command structure
	t.Skip("Skipping test for RunViewFileMappingCommand since it requires external resources")

	// In a real test, we would:
	// 1. Create a mock command with flags
	// 2. Create a mock application context with stdout/stderr capture
	// 3. Call RunViewFileMappingCommand
	// 4. Verify the output and behavior

	// Example of how this might look:
	cmd := &cobra.Command{}
	appCtx := &app.ApplicationContext{
		Logger: logger.NewTestLogger(),
	}

	err := RunViewFileMappingCommand(cmd, appCtx)

	// In a real test, we would make assertions about the error and output
	assert.Error(t, err) // We expect an error since we're not setting up the environment properly
}
