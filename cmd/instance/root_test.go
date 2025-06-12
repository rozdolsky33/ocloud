package instance

import (
	"testing"

	"github.com/spf13/cobra"
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

// TestGetAppContext tests the getAppContext function
func TestGetAppContext(t *testing.T) {
	// Create a command with no context
	cmd := &cobra.Command{}

	// Test that getAppContext returns an error when the context is nil
	_, err := getAppContext(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "command context is nil")
}
