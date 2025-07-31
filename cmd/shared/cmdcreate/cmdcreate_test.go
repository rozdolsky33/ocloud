package cmdcreate

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestCreateRootCmd tests the CreateRootCmd function with and without application context
func TestCreateRootCmd(t *testing.T) {
	// Test with nil application context
	rootCmd := CreateRootCmd(nil)
	assert.NotNil(t, rootCmd, "root command should not be nil")
	assert.Equal(t, "ocloud", rootCmd.Use)
	assert.Equal(t, "Interact with Oracle Cloud Infrastructure", rootCmd.Short)
	assert.True(t, rootCmd.SilenceUsage)

	// Verify that version command is added
	versionCmd := findSubcommand(rootCmd, "version")
	assert.NotNil(t, versionCmd, "version command should be added as a subcommand")

	// Verify that config command is added
	configCmd := findSubcommand(rootCmd, "config")
	assert.NotNil(t, configCmd, "config command should be added as a subcommand")

	// Verify that compute command is not added when appCtx is nil
	computeCmd := findSubcommand(rootCmd, "compute")
	assert.Nil(t, computeCmd, "compute command should not be added when appCtx is nil")

	// Test with mock application context
	mockAppCtx := &app.ApplicationContext{}
	rootCmdWithCtx := CreateRootCmd(mockAppCtx)
	assert.NotNil(t, rootCmdWithCtx, "root command should not be nil")

	// Verify that compute command is added when appCtx is not nil
	computeCmdWithCtx := findSubcommand(rootCmdWithCtx, "compute")
	assert.NotNil(t, computeCmdWithCtx, "compute command should be added when appCtx is not nil")

	// Verify that identity command is added when appCtx is not nil
	identityCmd := findSubcommand(rootCmdWithCtx, "identity")
	assert.NotNil(t, identityCmd, "identity command should be added when appCtx is not nil")

	// Verify that database command is added when appCtx is not nil
	databaseCmd := findSubcommand(rootCmdWithCtx, "database")
	assert.NotNil(t, databaseCmd, "database command should be added when appCtx is not nil")

	// Verify that network command is added when appCtx is not nil
	networkCmd := findSubcommand(rootCmdWithCtx, "network")
	assert.NotNil(t, networkCmd, "network command should be added when appCtx is not nil")
}

// TestCreateRootCmdWithoutContext tests the CreateRootCmdWithoutContext function
func TestCreateRootCmdWithoutContext(t *testing.T) {
	rootCmd := CreateRootCmdWithoutContext()
	assert.NotNil(t, rootCmd, "root command should not be nil")
	assert.Equal(t, "ocloud", rootCmd.Use)
	assert.Equal(t, "Interact with Oracle Cloud Infrastructure", rootCmd.Short)
	assert.True(t, rootCmd.SilenceUsage)

	// Verify that version command is added
	versionCmd := findSubcommand(rootCmd, "version")
	assert.NotNil(t, versionCmd, "version command should be added as a subcommand")

	// Verify that config command is added
	configCmd := findSubcommand(rootCmd, "config")
	assert.NotNil(t, configCmd, "config command should be added as a subcommand")

	// Verify that placeholder commands are added
	computeCmd := findSubcommand(rootCmd, "compute")
	assert.NotNil(t, computeCmd, "compute command should be added as a placeholder")

	// Verify it's a placeholder by checking that it returns an error when run
	err := computeCmd.RunE(computeCmd, []string{})
	assert.Error(t, err, "placeholder command should return an error when run")
	assert.Contains(t, err.Error(), "requires application initialization", "error message should indicate initialization is required")

	// Test that the RunE function is set on the root command
	assert.NotNil(t, rootCmd.RunE, "RunE function should be set on the root command")
}

// TestAddPlaceholderCommands tests the addPlaceholderCommands function
func TestAddPlaceholderCommands(t *testing.T) {
	rootCmd := &cobra.Command{
		Use:   "ocloud",
		Short: "Interact with Oracle Cloud Infrastructure",
	}

	addPlaceholderCommands(rootCmd)

	// Verify that placeholder commands are added
	commandTypes := []string{"compute", "identity", "database", "network"}
	for _, cmdType := range commandTypes {
		cmd := findSubcommand(rootCmd, cmdType)
		assert.NotNil(t, cmd, cmdType+" command should be added as a placeholder")

		// Verify it's a placeholder by checking that it returns an error when run
		err := cmd.RunE(cmd, []string{})
		assert.Error(t, err, "placeholder command should return an error when run")
		assert.Contains(t, err.Error(), "requires application initialization", "error message should indicate initialization is required")
	}
}

// findSubcommand is a helper function to find a subcommand by name
func findSubcommand(cmd *cobra.Command, name string) *cobra.Command {
	for _, subCmd := range cmd.Commands() {
		if subCmd.Name() == name {
			return subCmd
		}
	}
	return nil
}
