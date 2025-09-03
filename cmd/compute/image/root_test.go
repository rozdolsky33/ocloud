package image

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
)

// TestImageCommand tests the basic structure of the image command
func TestImageCommand(t *testing.T) {
	// Create a mock ApplicationContext
	appCtx := &app.ApplicationContext{}

	// Create a new image command
	cmd := NewImageCmd(appCtx)

	// Test that the image command is properly configured
	assert.Equal(t, "image", cmd.Use)
	assert.Equal(t, "Manage OCI Image", cmd.Short)
	assert.Equal(t, "Manage Oracle Cloud Infrastructure Compute Image - list all image or find image by name pattern.", cmd.Long)
	assert.Equal(t, "  ocloud compute image list\n  ocloud compute image find <image-name>", cmd.Example)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	assert.Nil(t, cmd.RunE, "RunE should be nil since the root command now has subcommands")

	// Test that the list subcommand is added and configured (TUI, no pagination flags)
	listCmd := findSubCommand(cmd, "list")
	assert.NotNil(t, listCmd, "list subcommand should be added")
	assert.Equal(t, "List all images", listCmd.Short)
	assert.NotNil(t, listCmd.RunE, "list subcommand should have a RunE function")

	// JSON flag is now a global flag, so it should not be in the local flags for a list
	jsonFlag := listCmd.Flags().Lookup(flags.FlagNameJSON)
	assert.Nil(t, jsonFlag, "json flag should not be added as a local flag to list subcommand")

	// But we should still be able to get its value using flags.GetBoolFlag
	useJSON := flags.GetBoolFlag(listCmd, flags.FlagNameJSON, false)
	assert.False(t, useJSON, "default value of json flag should be false")

	// Test that the get subcommand is added and has pagination flags
	getCmd := findSubCommand(cmd, "get")
	assert.NotNil(t, getCmd, "get subcommand should be added")
	assert.Equal(t, "Paginated Image Results", getCmd.Short)
	assert.NotNil(t, getCmd.RunE, "get subcommand should have a RunE function")

	limitFlag := getCmd.Flags().Lookup(flags.FlagNameLimit)
	assert.NotNil(t, limitFlag, "limit flag should be added to get subcommand")
	assert.Equal(t, flags.FlagShortLimit, limitFlag.Shorthand)
	assert.Equal(t, flags.FlagDescLimit, limitFlag.Usage)

	pageFlag := getCmd.Flags().Lookup(flags.FlagNamePage)
	assert.NotNil(t, pageFlag, "page flag should be added to get subcommand")
	assert.Equal(t, flags.FlagShortPage, pageFlag.Shorthand)
	assert.Equal(t, flags.FlagDescPage, pageFlag.Usage)

	// JSON flag is global, so it should not be in the local flags for get either
	jsonFlagGet := getCmd.Flags().Lookup(flags.FlagNameJSON)
	assert.Nil(t, jsonFlagGet, "json flag should not be added as a local flag to get subcommand")

	// Test that the find subcommand is added
	findCmd := findSubCommand(cmd, "find")
	assert.NotNil(t, findCmd, "find subcommand should be added")
	assert.Equal(t, "Find image by name pattern", findCmd.Short)
	assert.NotNil(t, findCmd.RunE, "find subcommand should have a RunE function")

	// JSON flag is now a global flag, so it should not be in the local flags
	jsonFlagFind := findCmd.Flags().Lookup(flags.FlagNameJSON)
	assert.Nil(t, jsonFlagFind, "json flag should not be added as a local flag to find subcommand")

	// But we should still be able to get its value using flags.GetBoolFlag
	useJSONFind := flags.GetBoolFlag(findCmd, flags.FlagNameJSON, false)
	assert.False(t, useJSONFind, "default value of json flag should be false")
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
