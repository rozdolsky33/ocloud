package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rozdolsky33/ocloud/internal/config"
)

// TestInstanceCommand tests the basic structure of the instance command
func TestInstanceCommand(t *testing.T) {
	// Test that the instance command is properly configured
	assert.Equal(t, "instance", InstanceCmd.Use)
	assert.Equal(t, "Find and list OCI instances", InstanceCmd.Short)
	assert.True(t, InstanceCmd.SilenceUsage)
	assert.True(t, InstanceCmd.SilenceErrors)
	assert.NotNil(t, InstanceCmd.PreRunE)
	assert.NotNil(t, InstanceCmd.RunE)

	// Test that the flags are added
	listFlag := InstanceCmd.Flags().Lookup(config.FlagNameList)
	assert.NotNil(t, listFlag, "list flag should be added")
	assert.Equal(t, config.FlagShortList, listFlag.Shorthand)
	assert.Equal(t, config.FlagDescList, listFlag.Usage)

	findFlag := InstanceCmd.Flags().Lookup(config.FlagNameFind)
	assert.NotNil(t, findFlag, "find flag should be added")
	assert.Equal(t, config.FlagShortFind, findFlag.Shorthand)
	assert.Equal(t, config.FlagDescFind, findFlag.Usage)

	imageDetailsFlag := InstanceCmd.Flags().Lookup(config.FlagNameImageDetails)
	assert.NotNil(t, imageDetailsFlag, "image-details flag should be added")
	assert.Equal(t, config.FlagShortImageDetails, imageDetailsFlag.Shorthand)
	assert.Equal(t, config.FlagDescImageDetails, imageDetailsFlag.Usage)
}

// TestInitApp tests the app.InitApp function
func TestInitApp(t *testing.T) {
	// This is just a placeholder test since we can't easily test InitApp without mocking the OCI SDK
	// The actual InitApp function is tested in the internal/app package
	t.Skip("Skipping test for InitApp since it requires mocking the OCI SDK")
}
