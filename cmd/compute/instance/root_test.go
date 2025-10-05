package instance

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
)

// TestInstanceCommand tests the basic structure of the instance command
func TestInstanceCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new instance command
	cmd := NewInstanceCmd(appCtx)

	// Test that the instance command is properly configured
	assert.Equal(t, "instance", cmd.Use)
	assert.Equal(t, "Manage OCI Compute instances — list, paginate, and search.", cmd.Short)
	assert.Equal(t, "List OCI Compute instances in a compartment. Supports paging through large result sets and filtering by name pattern.", cmd.Long)
	assert.Equal(t, "  ocloud compute instance get\n  ocloud compute instance list\n  ocloud compute instance search <value>", cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	assert.Nil(t, cmd.RunE, "RunE should be nil since the root command now has subcommands")

	// Test that the subcommands are added
	listCmd := findSubCommand(cmd, "get")
	assert.NotNil(t, listCmd, "get subcommand should be added")
	assert.Equal(t, "Paginated Instance Results", listCmd.Short)
	assert.NotNil(t, listCmd.RunE, "list subcommand should have a RunE function")

	// Test that the list subcommand has the appropriate flags
	limitFlag := listCmd.Flags().Lookup(flags.FlagNameLimit)
	assert.NotNil(t, limitFlag, "limit flag should be added to list subcommand")
	assert.Equal(t, flags.FlagShortLimit, limitFlag.Shorthand)
	assert.Equal(t, flags.FlagDescLimit, limitFlag.Usage)

	pageFlag := listCmd.Flags().Lookup(flags.FlagNamePage)
	assert.NotNil(t, pageFlag, "page flag should be added to list subcommand")
	assert.Equal(t, flags.FlagShortPage, pageFlag.Shorthand)
	assert.Equal(t, flags.FlagDescPage, pageFlag.Usage)

	// JSON flag is now a global flag, so it should not be in the local flags
	jsonFlag := listCmd.Flags().Lookup(flags.FlagNameJSON)
	assert.Nil(t, jsonFlag, "json flag should not be added as a local flag to list subcommand")

	// But we should still be able to get its value using flags.GetBoolFlag
	useJSON := flags.GetBoolFlag(listCmd, flags.FlagNameJSON, false)
	assert.False(t, useJSON, "default value of json flag should be false")

	// Test that the search subcommand is added
	searchCmd := findSubCommand(cmd, "search")
	assert.NotNil(t, searchCmd, "search subcommand should be added")
	assert.Equal(t, "Search instances by name pattern", searchCmd.Short)
	assert.NotNil(t, searchCmd.RunE, "search subcommand should have a RunE function")

	// Test that the search subcommand has the appropriate flags
	imageDetailsFlag := searchCmd.Flags().Lookup(flags.FlagNameAll)
	assert.NotNil(t, imageDetailsFlag, "image-details flag should be added to search subcommand")
	assert.Equal(t, flags.FlagShortAll, imageDetailsFlag.Shorthand)
	assert.Equal(t, flags.FlagDescAll, imageDetailsFlag.Usage)

	// JSON flag is now a global flag, so it should not be in the local flags
	jsonFlagSearch := searchCmd.Flags().Lookup(flags.FlagNameJSON)
	assert.Nil(t, jsonFlagSearch, "json flag should not be added as a local flag to search subcommand")

	// But we should still be able to get its value using flags.GetBoolFlag
	useJSONSearch := flags.GetBoolFlag(searchCmd, flags.FlagNameJSON, false)
	assert.False(t, useJSONSearch, "default value of json flag should be false")
}

// findSubCommand is a helper function to find a subcommand by name
func findSubCommand(cmd *cobra.Command, name string) *cobra.Command {
	for _, subCmd := range cmd.Commands() {
		if subCmd.Name() == name {
			return subCmd
		}
	}
	return nil
}

// TestInitApp tests the app.InitApp function
func TestInitApp(t *testing.T) {
	// This is just a placeholder test since we can't easily test InitApp without mocking the OCI SDK
	// The actual InitApp function is tested in the shared/app package
	t.Skip("Skipping test for InitApp since it requires mocking the OCI SDK")
}
