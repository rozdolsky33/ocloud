package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
)

// TestInstanceCommand tests the basic structure of the instance command
func TestInstanceCommand(t *testing.T) {
	// Create a mock AppContext
	appCtx := &app.AppContext{}

	// Create a new instance command
	cmd := NewInstanceCmd(appCtx)

	// Test that the instance command is properly configured
	assert.Equal(t, "instance", cmd.Use)
	assert.Equal(t, "Find and list OCI instances", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	assert.NotNil(t, cmd.RunE)

	// Test that the flags are added
	listFlag := cmd.Flags().Lookup(flags.FlagNameList)
	assert.NotNil(t, listFlag, "list flag should be added")
	assert.Equal(t, flags.FlagShortList, listFlag.Shorthand)
	assert.Equal(t, flags.FlagDescList, listFlag.Usage)

	findFlag := cmd.Flags().Lookup(flags.FlagNameFind)
	assert.NotNil(t, findFlag, "find flag should be added")
	assert.Equal(t, flags.FlagShortFind, findFlag.Shorthand)
	assert.Equal(t, flags.FlagDescFind, findFlag.Usage)

	imageDetailsFlag := cmd.Flags().Lookup(flags.FlagNameImageDetails)
	assert.NotNil(t, imageDetailsFlag, "image-details flag should be added")
	assert.Equal(t, flags.FlagShortImageDetails, imageDetailsFlag.Shorthand)
	assert.Equal(t, flags.FlagDescImageDetails, imageDetailsFlag.Usage)
}

// TestInitApp tests the app.InitApp function
func TestInitApp(t *testing.T) {
	// This is just a placeholder test since we can't easily test InitApp without mocking the OCI SDK
	// The actual InitApp function is tested in the internal/app package
	t.Skip("Skipping test for InitApp since it requires mocking the OCI SDK")
}
